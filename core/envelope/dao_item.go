package envelope

import (
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

// 查询
func (dao *RedEnvelopeItemDao) GetOne(itemNo int64) *RedEnvelopeItem {
	item := &RedEnvelopeItem{
		ItemNo: itemNo,
	}
	ok, err := dao.runner.GetOne(item)
	if err != nil || !ok {
		return nil
	}
	return item
}

// 插入
func (dao *RedEnvelopeItemDao) Insert(item *RedEnvelopeItem) (int64, error) {
	ret, err := dao.runner.Insert(item)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}
