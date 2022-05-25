package envelope

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

// 查询
func (dao *RedEnvelopeItemDao) GetOne(itemNo string) *RedEnvelopeItem {
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

func (dao *RedEnvelopeItemDao) FindItems(envelopeNo string) (items []*RedEnvelopeItem) {
	sql := `
		SELECT
			*
		FROM
			red_envelope_item
		WHERE
			envelope_no=?
	`
	err := dao.runner.Find(items, sql, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}
