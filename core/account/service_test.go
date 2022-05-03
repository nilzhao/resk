package account

import (
	"resk/services"
	// _ "resk/testx"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccountService_CreateAccount(t *testing.T) {
	amount := "100"
	dto := services.AccountCreatedDTO{
		UserId:      ksuid.New().Next().String(),
		Username:    "测试用户",
		Amount:      amount,
		AccountName: "测试账户",
	}
	service := accountService{}
	Convey("账户创建", t, func() {
		accountDto, err := service.CreateAccount(dto)
		So(err, ShouldBeNil)
		So(accountDto, ShouldNotBeNil)
		So(accountDto.UserId, ShouldEqual, dto.UserId)
		So(accountDto.Username, ShouldEqual, dto.Username)
		So(accountDto.Balance.String(), ShouldEqual, amount)
	})
}

func TestAccountService_Transfer(t *testing.T) {
	Convey("转账", t, func() {
		// 准备两个账户
		// 账户 1
		accountDto1 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户1",
			AccountName:  "测试账户1",
			AccountType:  1,
			CurrencyCode: "CNY",
			Amount:       "100",
		}
		// 账户 2
		accountDto2 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户2",
			AccountName:  "测试账户2",
			AccountType:  1,
			CurrencyCode: "CNY",
			Amount:       "100",
		}
		service := new(accountService)
		account1, err := service.CreateAccount(accountDto1)
		So(err, ShouldBeNil)
		So(account1, ShouldNotBeNil)
		account2, err := service.CreateAccount(accountDto2)
		So(err, ShouldBeNil)
		So(account2, ShouldNotBeNil)
		Convey("从账户 1 转入账户 2:余额足够,转账成功", func() {
			body := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: account2.AccountNo,
				UserId:    account2.UserId,
				Username:  account2.Username,
			}
			amount, _ := decimal.NewFromString("1")
			transferDto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.EnvelopeOutgoing,
				ChangeFlag:  services.FlagAccountOut,
				Decs:        "转出",
			}
			status, err := service.Transfer(transferDto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferStatusSuccess)
			retAccount1 := service.GetAccount(account1.AccountNo)
			So(retAccount1, ShouldNotBeNil)
			So(retAccount1.Balance.String(), ShouldEqual, account1.Balance.Sub(amount).String())
		})

		//从账户1转入账户2一定金额，但余额不足，转账会失败
		Convey("余额不足，从账户1转入账户2一定金额", func() {
			//转账过程需要2个交易的参与者：交易主体和交易对象
			body := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: account2.AccountNo,
				UserId:    account2.UserId,
				Username:  account2.Username,
			}
			amount := account1.Balance.Add(decimal.NewFromFloat(200))
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagAccountOut,
				Decs:        "转出",
			}
			status, err := service.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferStatusInsufficient)

			retAccount1 := service.GetAccount(account1.AccountNo)
			So(retAccount1, ShouldNotBeNil)
			So(retAccount1.Balance.String(), ShouldEqual, account1.Balance.String())
		})
		//给账户1储值
		Convey("给账户1储值", func() {
			//转账过程需要2个交易的参与者：交易主体和交易对象
			body := services.TradeParticipator{
				AccountNo: account1.AccountNo,
				UserId:    account1.UserId,
				Username:  account1.Username,
			}
			target := body
			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagAccountIn,
				Decs:        "储值",
			}
			status, err := service.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferStatusSuccess)

			retAccount1 := service.GetAccount(account1.AccountNo)
			So(retAccount1, ShouldNotBeNil)
			So(retAccount1.Balance.String(), ShouldEqual, account1.Balance.Add(amount).String())

		})
	})
}
