package account

import (
	"database/sql"
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"resk/infra/base"
	_ "resk/testx"
	"testing"
)

func TestAccountDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		Convey("通过编号查询账户数据", t, func() {
			account := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试账户",
				UserId:      ksuid.New().Next().String(),
				Username: sql.NullString{
					String: "测试用户",
					Valid:  true,
				},
			}
			id, err := dao.Insert(account)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			ret := dao.GetOne(account.AccountNo)
			So(ret, ShouldNotBeNil)
			So(ret.Balance.String(), ShouldEqual, account.Balance.String())
			So(ret.CreatedAt, ShouldNotBeNil)
			So(ret.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_GetByUserId(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		Convey("通过用户 ID 查询账户数据", t, func() {

			dao := &AccountDao{
				runner: runner,
			}

			account := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试账户",
				UserId:      ksuid.New().Next().String(),
				Username: sql.NullString{
					String: "测试用户",
					Valid:  true,
				},
				AccountType: 2,
			}
			id, err := dao.Insert(account)
			So(err, ShouldBeNil)
			So(id, ShouldNotBeNil)

			ret := dao.GetByUserId(account.UserId, account.AccountType)
			So(err, ShouldBeNil)
			So(ret, ShouldNotBeNil)
			So(ret.Id, ShouldNotBeNil)
			So(ret.Balance.String(), ShouldEqual, account.Balance.String())
			So(ret.CreatedAt, ShouldNotBeNil)
			So(ret.UpdatedAt, ShouldNotBeNil)
		})

		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		Convey("更新账户余额", t, func() {
			balance := decimal.NewFromFloat(100)
			account := &Account{
				Balance:     balance,
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试账户",
				UserId:      ksuid.New().Next().String(),
				Username: sql.NullString{
					String: "测试用户",
					Valid:  true,
				},
			}
			id, err := dao.Insert(account)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			fmt.Println("id", id)

			Convey("增加余额", func() {
				amount := decimal.NewFromFloat(10)
				rows, err := dao.UpdateBalance(account.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)

				newAccount := dao.GetOne(account.AccountNo)
				newBalance := balance.Add(amount)
				So(newAccount.Balance, ShouldEqual, newBalance)
			})
			Convey("减少余额-余额够", func() {
				amount := decimal.NewFromFloat(-10)
				rows, err := dao.UpdateBalance(account.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)

				newAccount := dao.GetOne(account.AccountNo)
				newBalance := balance.Add(amount)
				So(newAccount.Balance, ShouldEqual, newBalance)
			})
			Convey("减少余额-余额不够", func() {
				account1 := dao.GetOne(account.AccountNo)
				So(account1, ShouldNotBeNil)
				amount := decimal.NewFromFloat(-101)
				rows, err := dao.UpdateBalance(account.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 0)

				account2 := dao.GetOne(account.AccountNo)
				So(account2, ShouldNotBeNil)

				So(account1.Balance, ShouldEqual, account2.Balance)
			})
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
