package envelope

import (
	"context"
	"errors"
	"fmt"
	"resk/infra/base"
	"resk/services"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type itemDomain struct {
	RedEnvelopeItem
}

// 生成 itemNo
func (d *itemDomain) createItemNo() {
	d.ItemNo = ksuid.New().Next().String()
}

// 创建 Item
func (d *itemDomain) Create(item services.RedEnvelopeItemDTO) {
	d.RedEnvelopeItem.FromDTO(&item)
	d.RecvUsername.Valid = true
	d.createItemNo()
}

// 保存 Item 数据
func (d *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&d.RedEnvelopeItem)
		return err
	})
	return id, err
}

// 通过 itemNo 查询抢红包明细数据
func (d *itemDomain) GetOne(ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetOne(itemNo)
		if po == nil {
			return errors.New(fmt.Sprintf("查询红包失败,红包编号 %s 不存在", itemNo))
		}
		dto = po.ToDTO()
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
	return dto
}

// 通过 envelopeNo 查询已抢到红包列表
func (d *itemDomain) FindItems(envelopeNo string) (itemDtos []*services.RedEnvelopeItemDTO) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		items := dao.FindItems(envelopeNo)
		itemDtos = make([]*services.RedEnvelopeItemDTO, 0)
		for _, po := range items {
			itemDtos = append(itemDtos, po.ToDTO())
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
	return itemDtos
}
