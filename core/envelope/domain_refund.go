package envelope

import (
	"context"
	"errors"
	"resk/core/account"
	"resk/infra/base"
	"resk/services"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

const pageSize = 100

type ExpiredEnvelopeDomain struct {
	expiredGoods []RedEnvelopeGoods
	offset       int
}

// 查询出过期红包
func (e *ExpiredEnvelopeDomain) Next() (ok bool) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{
			runner,
		}
		e.expiredGoods = dao.FindExpired(e.offset, pageSize)
		if len(e.expiredGoods) > 0 {
			e.offset += len(e.expiredGoods)
			ok = true
		}
		return nil
	})
	return ok
}

func (e *ExpiredEnvelopeDomain) Expired() (err error) {
	for e.Next() {
		for _, goods := range e.expiredGoods {
			logrus.Debugf("过期红包退款开始, %v", goods)
			err := e.ExpiredOne(goods)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("过期红包退款结束, %v", goods)
		}
	}
	return err
}

func (e *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	// 创建退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	refund.RemainAmount = goods.RemainAmount.Mul(decimal.NewFromFloat(-1))
	refund.RemainQuantity = -goods.RemainQuantity
	refund.Status = services.OrderExpired
	refund.PayStatus = services.Refunding
	refund.OriginEnvelopeNo = goods.EnvelopeNo
	domain := goodsDomain{
		RedEnvelopeGoods: refund,
	}
	domain.createEnvelopeNo()
	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		id, err := domain.Save(txCtx)
		if err != nil || id == 0 {
			return errors.New("创建退款订单失败")
		}
		// 修改原订单的状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 调用资金账户接口进行转账: 把退款订单里的钱,返还给发红包的账户
	accountDomain := account.NewAccountDomain()
	systemAccount := base.GetSystemAccount()
	account := services.GetAccountService().GetAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的红包资金账户:" + goods.UserId)
	}

	// 调用资金账户接口进行转账
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}

	target := services.TradeParticipator{
		AccountNo: account.AccountNo,
		UserId:    account.UserId,
		Username:  account.Username,
	}

	dto := services.AccountTransferDTO{
		TradeNo:     refund.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      goods.RemainAmount,
		ChangeType:  services.EnvelopeExpiredRefund,
		ChangeFlag:  services.FlagAccountIn,
		Decs:        "过期红包退款" + goods.EnvelopeNo,
	}
	status, err := accountDomain.Transfer(dto)
	if status != services.TransferStatusSuccess {
		return err
	}

	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 修改原订单的状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredSuccess)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		// 修改退款订单状态
		rows, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredSuccess)
		if err != nil || rows == 0 {
			return errors.New("更新退款订单状态失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return err
}

// 发起退款流程
