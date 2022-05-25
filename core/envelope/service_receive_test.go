package envelope

import (
	"strconv"
	"testing"

	"resk/services"
	_ "resk/testx"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedEnvelopeService_Receive(t *testing.T) {
	// 1. 准备红包资金账户
	envelopeService := services.GetRedEnvelopeService()
	accountService := services.GetAccountService()
	size := 10
	Convey("收红包", t, func() {
		accounts := make([]*services.AccountDTO, 0)
		for i := 0; i < size; i++ {
			account := services.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i+1),
				AccountName:  "测试账户" + strconv.Itoa(i+1),
				AccountType:  int8(services.EnvelopeAccountType),
				CurrencyCode: "CNY",
				Amount:       "2000",
			}
			accountDto, err := accountService.CreateAccount(account)
			So(accountDto, ShouldNotBeNil)
			So(err, ShouldBeNil)
			accounts = append(accounts, accountDto)
		}
		So(len(accounts), ShouldEqual, size)

		Convey("收普通红包", func() {
			// 先发一个普通红包
			accountDto := accounts[0]
			goods := services.RedEnvelopeSendingDTO{
				UserId:       accountDto.UserId,
				Username:     accountDto.Username,
				EnvelopeType: services.GeneralEnvelopeType,
				Amount:       decimal.NewFromFloat(8.88),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			activity, err := envelopeService.SendOut(goods)
			So(activity, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(activity.Link, ShouldNotBeEmpty)
			So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 所有的账户抢红包
			remainAmount := activity.Amount
			for _, account := range accounts {
				receiveDto := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   activity.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				itemDto, err := envelopeService.Receive(receiveDto)
				So(itemDto, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(itemDto.Amount, ShouldEqual, activity.AmountOne)
				// 剩余金额
				remainAmount = remainAmount.Sub(activity.AmountOne)
				So(itemDto.RemainAmount, ShouldEqual, remainAmount)
			}
		})
		Convey("收碰运气红包", func() {
			// 先发一个普通红包
			accountDto := accounts[0]
			goods := services.RedEnvelopeSendingDTO{
				UserId:       accountDto.UserId,
				Username:     accountDto.Username,
				EnvelopeType: services.LuckyEnvelopeType,
				Amount:       decimal.NewFromFloat(8.88),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			activity, err := envelopeService.SendOut(goods)
			So(activity, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(activity.Link, ShouldNotBeEmpty)
			So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 所有的账户抢红包
			remainAmount := activity.Amount
			total := decimal.NewFromFloat(0)
			for _, account := range accounts {
				receiveDto := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   activity.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				itemDto, err := envelopeService.Receive(receiveDto)
				if itemDto != nil {
					total = total.Add(itemDto.Amount)
				}
				So(itemDto, ShouldNotBeNil)
				So(err, ShouldBeNil)
				// 剩余金额
				remainAmount = remainAmount.Sub(itemDto.Amount)
				So(itemDto.RemainAmount, ShouldEqual, remainAmount)
			}
			So(total.String(), ShouldEqual, goods.Amount.String())
		})
	})
}
