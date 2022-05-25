package envelope

import (
	"context"
	"resk/infra/base"
	"resk/services"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type goodsDomain struct {
	RedEnvelopeGoods
	itemDomain itemDomain
}

// 生成一个红包编号
func (domain *goodsDomain) createEnvelopeNo() {
	domain.EnvelopeNo = ksuid.New().Next().String()
}

// 创建一个红包商品对象
func (domain *goodsDomain) Create(goods services.RedEnvelopeGoodsDTO) {
	domain.RedEnvelopeGoods.FromDTO(&goods)
	domain.RemainQuantity = goods.Quantity
	domain.Username.Valid = true
	domain.Blessing.Valid = true
	if domain.EnvelopeType == services.GeneralEnvelopeType {
		domain.Amount = goods.AmountOne.Mul(
			decimal.NewFromFloat(float64(goods.Quantity)))
	}
	if domain.EnvelopeType == services.LuckyEnvelopeType {
		domain.AmountOne = decimal.NewFromFloat(0)
	}
	domain.RemainAmount = domain.Amount
	// 过期时间
	domain.ExpiredAt = time.Now().Add(24 * time.Hour)
	domain.Status = services.OrderCreate
	domain.createEnvelopeNo()
}

// 保存到红包商品表
func (domain *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		id, err = dao.Insert(&domain.RedEnvelopeGoods)
		return err
	})
	return id, err
}

func (domain *goodsDomain) CreateAndSave(ctx context.Context, goods services.RedEnvelopeGoodsDTO) (id int64, err error) {
	domain.Create(goods)
	return domain.Save(ctx)
}

// 查询商品信息
func (domain *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return goods
}
