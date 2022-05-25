package envelope

import (
	"resk/services"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// 插入
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	ret, err := dao.runner.Insert(po)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

// 更新红包余额和数量
func (dao *RedEnvelopeGoodsDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := `
		UPDATE
			red_envelope_goods
		SET
			remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)),
			remain_quantity=remain_quantity-1
		WHERE
			envelope_no=?
			AND remain_quantity>=1
			AND remain_amount>=CAST(? AS DECIMAL(30,6))
	`
	ret, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

// 更新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := `
		UPDATE
			red_envelope_goods
		SET
			order_status=?
		WHERE
			envelope_no=?
	`
	ret, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

// 查询: 根据红包编号
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	po := &RedEnvelopeGoods{
		EnvelopeNo: envelopeNo,
	}
	ok, err := dao.runner.GetOne(po)
	if err != nil || !ok {
		return nil
	}
	return po
}

// 过期: 把过期的所有红包都查询出来,分页
func (dao *RedEnvelopeGoodsDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := `
		SELECT
			*
		FROM
			red_envelope_goods
		WHERE
			expired_at>?
		LIMIT
			?,?
	`
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
