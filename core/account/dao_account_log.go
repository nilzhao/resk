package account

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountLogDao struct {
	runner *dbx.TxRunner
}

// GetOne 查询数据库持久化对象的单实例，获取一行数据
func (dao *AccountLogDao) GetOne(logNo string) *AccountLog {
	accountLog := &AccountLog{
		LogNo: logNo,
	}
	ok, err := dao.runner.GetOne(accountLog)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return accountLog
}

// GetByTradeNo 交易编号获取账户日志
func (dao *AccountLogDao) GetByTradeNo(tradeNo string) *AccountLog {
	accountLog := &AccountLog{}
	sql := "select * from  account_log where trade_no=?"
	ok, err := dao.runner.Get(accountLog, sql, tradeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return accountLog
}

// Insert 账户日志数据的插入
func (dao *AccountLogDao) Insert(accountLog *AccountLog) (id int64, err error) {
	ret, err := dao.runner.Insert(accountLog)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}
