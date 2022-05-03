package account

import (
	"resk/services"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户",
		Balance:  decimal.NewFromFloat(0),
		Status:   1,
	}
	domain := new(accountDomain)
	Convey("账户创建", t, func() {
		retDto, err := domain.Create(dto)
		So(err, ShouldBeNil)
		So(retDto, ShouldNotBeNil)
		So(retDto.UserId, ShouldEqual, dto.UserId)
		So(retDto.Username, ShouldEqual, dto.Username)
		So(retDto.Status, ShouldEqual, dto.Status)
		So(retDto.Balance.String(), ShouldEqual, dto.Balance.String())
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	// 两个账户,交易主体要有余额
	dto1 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户1",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	dto2 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户2",
		Balance:  decimal.NewFromFloat(0),
		Status:   1,
	}
	domain1 := accountDomain{}
	domain2 := accountDomain{}
	Convey("转账", t, func() {
		// 准备两个账号
		// 账户 1
		retDto1, err := domain1.Create(*dto1)
		So(err, ShouldBeNil)
		So(retDto1, ShouldNotBeNil)
		dto1 = retDto1
		// 账户 2
		retDto2, err := domain2.Create(*dto2)
		So(err, ShouldBeNil)
		So(retDto2, ShouldNotBeNil)
		dto2 = retDto2

		Convey("余额充足，应转出成功", func() {
			amount := decimal.NewFromFloat(1)
			transferDto := services.AccountTransferDTO{
				TradeBody: services.TradeParticipator{
					AccountNo: dto1.AccountNo,
					UserId:    dto1.UserId,
					Username:  dto1.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: dto2.AccountNo,
					UserId:    dto2.UserId,
					Username:  dto2.Username,
				},
				Amount:     amount,
				ChangeFlag: services.FlagAccountOut,
				ChangeType: services.EnvelopeOutgoing,
				Decs:       "转账给他人",
			}
			// 转账状态应该是 `成功`
			status, err := domain1.Transfer(transferDto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferStatusSuccess)
			// 验证转账后账户的余额
			account1 := domain1.GetAccount(dto1.AccountNo)
			So(account1, ShouldNotBeNil)
			So(account1.Balance, ShouldEqual, dto1.Balance.Sub(amount))
		})

		Convey("余额不足，应转出失败", func() {
			amount := dto1.Balance.Add(decimal.NewFromFloat(1))
			transferDto := services.AccountTransferDTO{
				TradeBody: services.TradeParticipator{
					AccountNo: dto1.AccountNo,
					UserId:    dto1.UserId,
					Username:  dto1.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: dto2.AccountNo,
					UserId:    dto2.UserId,
					Username:  dto2.Username,
				},
				Amount:     amount,
				ChangeFlag: services.FlagAccountOut,
				ChangeType: services.EnvelopeOutgoing,
				Decs:       "转账给他人",
			}
			// 转账状态应该是 `成功`
			status, err := domain1.Transfer(transferDto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferStatusInsufficient)
			// 验证转账后账户的余额
			account1 := domain1.GetAccount(dto1.AccountNo)
			So(account1, ShouldNotBeNil)
			So(account1.Balance, ShouldEqual, dto1.Balance)
		})

		Convey("充值，应成功", func() {
			amount := decimal.NewFromFloat(1)
			transferDto := services.AccountTransferDTO{
				TradeBody: services.TradeParticipator{
					AccountNo: dto1.AccountNo,
					UserId:    dto1.UserId,
					Username:  dto1.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: dto2.AccountNo,
					UserId:    dto2.UserId,
					Username:  dto2.Username,
				},
				Amount:     amount,
				ChangeFlag: services.FlagAccountIn,
				ChangeType: services.AccountStoreValue,
				Decs:       "充值",
			}
			// 转账状态应该是 `成功`
			status, err := domain1.Transfer(transferDto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferStatusSuccess)
			// 验证转账后账户的余额
			account1 := domain1.GetAccount(dto1.AccountNo)
			So(account1, ShouldNotBeNil)
			So(account1.Balance, ShouldEqual, dto1.Balance.Add(amount))
		})

	})
}
