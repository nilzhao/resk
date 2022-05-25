package envelope

import (
	"context"
	"path"
	"resk/core/account"
	"resk/infra/base"
	"resk/services"

	"github.com/tietang/dbx"
)

func (d *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包商品
	d.Create(goods)
	// 创建活动
	activity = new(services.RedEnvelopeActivity)
	// 红包链接
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeActivityDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)
	accountDomain := account.NewAccountDomain()
	err = base.Tx(func(runner *dbx.TxRunner) (err error) {
		ctx := base.WithValueContext(context.Background(), runner)
		// 保存红包商品
		id, err := d.Save(ctx)
		if id <= 0 || err != nil {
			return err
		}
		// 红包金额支付
		// 1. 需要红包中间商的红包资金账户,定义在配置文件中,事先初始化到资金账户表中
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}
		transfer := services.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     goods.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagAccountOut,
			Decs:        "红包金额支付",
		}
		// 2. 从红包发送人的资金账户中扣减红包金额
		status, err := accountDomain.TransferWithContext(ctx, transfer)
		if status != services.TransferStatusSuccess {
			return err
		}
		// 3. 将扣减的红包总金额转入红包中间商的红包资金账户
		transfer = services.AccountTransferDTO{
			TradeBody:   target,
			TradeTarget: body,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagAccountIn,
			Decs:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContext(ctx, transfer)
		if status == services.TransferStatusSuccess {
			return nil
		}
		// 扣减金额没问题,返回活动
		return err
	})

	if err != nil {
		return nil, err
	}
	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()
	return activity, err
}
