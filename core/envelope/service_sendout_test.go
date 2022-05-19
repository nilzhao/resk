package envelope

import (
	"resk/services"
	"testing"

	_ "resk/testx"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedEnvelopeService_SendOut(t *testing.T) {
	accountService := services.GetAccountService()

	accountDto := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户1",
		AccountName:  "测试账户1",
		AccountType:  int8(services.EnvelopeAccountType),
		CurrencyCode: "CNY",
		Amount:       "1000",
	}
	envelopeService := services.GetRedEnvelopeService()
	Convey("发红包", t, func() {
		// 准备资金账户
		acDto, err := accountService.CreateAccount(accountDto)
		So(acDto, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("发普通红包", func() {
			goods := services.RedEnvelopeSendingDTO{
				UserId:       accountDto.UserId,
				Username:     accountDto.Username,
				EnvelopeType: services.GeneralEnvelopeType,
				Amount:       decimal.NewFromFloat(8.88),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			at, err := envelopeService.SendOut(goods)
			So(at, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			q := decimal.NewFromFloat(float64(dto.Quantity))
			So(dto.Amount, ShouldEqual, goods.Amount.Mul(q))
		})
		Convey("发碰运气红包", func() {
			goods := services.RedEnvelopeSendingDTO{
				UserId:       accountDto.UserId,
				Username:     accountDto.Username,
				EnvelopeType: services.LuckyEnvelopeType,
				Amount:       decimal.NewFromFloat(88.8),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			at, err := envelopeService.SendOut(goods)
			So(at, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			So(dto.Amount, ShouldEqual, goods.Amount)
		})
	})
}
