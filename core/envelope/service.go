package envelope

import (
	"errors"
	"resk/infra/base"
	"resk/services"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.IRedEnvelopeService = new(redEnvelopeService)
	})
}

type redEnvelopeService struct {
}

// Get implements services.RedEnvelopeService
func (r *redEnvelopeService) Get(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {

	panic("unimplemented")
}

// Receive implements services.RedEnvelopeService
func (r *redEnvelopeService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 参数校验
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	// 获取当前收红包用户的账户信息
	account := services.GetAccountService().GetAccountByUserId(dto.RecvUserId)
	if account == nil {
		return nil, errors.New("红包账户不存在: user_id=" + dto.RecvUserId)
	}
	// 尝试接收红包
	domain := goodsDomain{}
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}

// Refund implements services.RedEnvelopeService
func (r *redEnvelopeService) Refund(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	panic("unimplemented")
}

// 发红包
func (r *redEnvelopeService) SendOut(dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 验证
	err = base.ValidateStruct(dto)
	if err != nil {
		return nil, err
	}

	// 获取红包发送人的资金账户信息
	account := services.GetAccountService().GetAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户账户不存在" + dto.UserId)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo
	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	// 执行发送红包的逻辑
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if activity == nil || err != nil {
		logrus.Error(err)
	}

	return activity, nil
}
