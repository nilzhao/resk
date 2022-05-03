package account

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

// GetOne 查询数据库持久化对象的单实例，获取一行数据
func (dao *AccountDao) GetOne(accountNo string) *Account {
	account := &Account{
		AccountNo: accountNo,
	}
	ok, err := dao.runner.GetOne(account)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return account
}

// GetByUserId 通过用户 ID 和账户类型来查询账户信息
func (dao *AccountDao) GetByUserId(userId string, accountType int8) *Account {
	account := &Account{}
	sql := "select * from  account where user_id=? and account_type=?"
	ok, err := dao.runner.Get(account, sql, userId, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return account
}

// Insert 账户数据的插入
func (dao *AccountDao) Insert(account *Account) (id int64, err error) {
	ret, err := dao.runner.Insert(account)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

// UpdateBalance 账户余额的更新
func (dao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (rows int64, err error) {
	sql := `
		update account
		set balance=balance+CAST(? AS DECIMAL(30,6))
		where account_no=?
		and balance>=-1*CAST(? AS DECIMAL(30,6)) 
	`
	ret, err := dao.runner.Exec(sql, amount.String(), accountNo, amount.String())
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

// UpdateStatus 账户状态更新
func (dao AccountDao) UpdateStatus(accountNo string, status int) (rows int64, err error) {
	sql := `
		update account
		set status=?
		where account_no=?
	`
	ret, err := dao.runner.Exec(sql, status, accountNo)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}
