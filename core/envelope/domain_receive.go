package envelope

import (
	"context"
	"database/sql"
	"errors"
	"resk/core/account"
	"resk/infra/algo"
	"resk/infra/base"
	"resk/services"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

// 人民币单位转换: 元 -> 分
var multiple = decimal.NewFromFloat(100)

func (d *goodsDomain) Receive(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO,
) (item *services.RedEnvelopeItemDTO, err error) {
	// 1. 创建收红包的订单明细
	d.preCreateItem(dto)
	// 2. 查询出当前红包的剩余数量和剩余金额信息
	goods := d.Get(dto.EnvelopeNo)
	if goods == nil {
		return nil, errors.New("红包商品不存在" + dto.EnvelopeNo)
	}
	// 3. 校验剩余红包和剩余金额
	// - 如果没有剩余,直接返回无可用红包金额
	if goods.Quantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("没有足够的红包金额和数量了")
	}
	// 4. 使用红包算法计算红包金额
	nextAmount := d.nextAmount(goods)
	logrus.Infof("账户 NO %s, nextAmount %s", dto.AccountNo, nextAmount.String())
	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 5. 使用乐观锁更新语句,尝试更新剩余数量和剩余金额
		dao := RedEnvelopeGoodsDao{runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		// - 如果更新成功,也就是返回 1,表示抢到了红包
		// - 如果更新失败,也就是返回 0,表示无可用红包数量和金额,抢红包失败
		if rows <= 0 || err != nil {
			return errors.New("没有足够的红包和金额了")
		}
		// 6. 保存订单明细数据
		d.itemDomain.Quantity = 1
		d.itemDomain.PayStatus = int(services.Paying)
		d.itemDomain.AccountNo = dto.AccountNo
		d.itemDomain.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.itemDomain.Amount = nextAmount
		txCtx := base.WithValueContext(ctx, runner)
		_, err = d.itemDomain.Save(txCtx)
		if err != nil {
			return err
		}
		// 7. 将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		status, err := d.transfer(txCtx, dto)
		if status == services.TransferStatusSuccess {
			return nil
		}
		return err
	})
	return d.itemDomain.ToDTO(), err
}

func (d *goodsDomain) transfer(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO,
) (status services.TransferStatus, err error) {
	systemAccount := base.GetSystemAccount()

	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}

	transferDto := services.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     dto.EnvelopeNo,
		Amount:      d.itemDomain.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagAccountIn,
		Decs:        "红包收入",
	}
	accountDomain := account.NewAccountDomain()
	return accountDomain.TransferWithContext(ctx, transferDto)
}

// 预创建收红包的订单明细
func (d *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	d.itemDomain.AccountNo = dto.AccountNo
	d.itemDomain.EnvelopeNo = dto.EnvelopeNo
	d.itemDomain.RecvUsername = sql.NullString{
		String: dto.RecvUsername,
		Valid:  true,
	}
	d.itemDomain.RecvUserId = dto.RecvUserId
	d.itemDomain.createItemNo()
}

// 计算红包金额``
func (d *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		amount = goods.RemainAmount
	} else if goods.EnvelopeType == services.GeneralEnvelopeType {
		amount = goods.AmountOne
	} else if goods.EnvelopeType == services.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromInt(next).Div(multiple)
	} else {
		logrus.Error("不支持的红包类型")
	}
	return amount
}
