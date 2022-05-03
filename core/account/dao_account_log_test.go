package account

import (
	"resk/infra/base"
	"resk/services"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
)

func TestAccountLogDao(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountLogDao{
			runner: runner,
		}
		accountLog := &AccountLog{
			LogNo:      ksuid.New().Next().String(),
			TradeNo:    ksuid.New().Next().String(),
			AccountNo:  ksuid.New().Next().String(),
			UserId:     ksuid.New().Next().String(),
			Username:   "测试用户",
			Status:     1,
			Amount:     decimal.NewFromFloat(1),
			Balance:    decimal.NewFromFloat(100),
			ChangeFlag: services.FlagAccountCreated,
			ChangeType: services.AccountCreated,
		}
		Convey("测试账户流水", t, func() {
			Convey("插入一条流水", func() {
				id, err := dao.Insert(accountLog)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)
			})

			Convey("通过Log编号获取账户日志", func() {
				ret := dao.GetOne(accountLog.LogNo)
				So(ret, ShouldNotBeNil)
				So(ret.AccountNo, ShouldEqual, accountLog.AccountNo)
				So(ret.LogNo, ShouldEqual, accountLog.LogNo)
			})

			Convey("通过交易编号获取账户日志", func() {
				ret := dao.GetByTradeNo(accountLog.TradeNo)
				So(ret, ShouldNotBeNil)
				So(ret.TradeNo, ShouldEqual, accountLog.TradeNo)
				So(ret.AccountNo, ShouldEqual, accountLog.AccountNo)
			})
		})

		return nil

	})

	if err != nil {
		logrus.Error(err)
	}
}
