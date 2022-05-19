package account

import (
	"context"
	"errors"
	"resk/infra/base"
	"resk/services"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

// 有状态的,每次使用时都要实例化

type accountDomain struct {
	account    Account
	accountLog AccountLog
}

func NewAccountDomain() *accountDomain {
	return new(accountDomain)
}

// 创建logNo
func (domain *accountDomain) createAccountLogNo() {
	// 暂时采用 ksuid
	// TODO: 改为分布式 ID
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

// 创建 accountNo
func (domain *accountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

// 创建流水的记录
func (domain *accountDomain) createAccountLog() {
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo
	// 交易主体信息
	domain.accountLog.AccountNo = domain.account.UserId
	domain.accountLog.UserId = domain.account.UserId
	domain.accountLog.Username = domain.account.Username.String
	// 交易对象信息
	domain.accountLog.TargetAccountNo = domain.account.UserId
	domain.accountLog.TargetUserId = domain.account.UserId
	domain.accountLog.TargetUsername = domain.account.Username.String
	if domain.accountLog.ChangeType == services.AccountCreated {
		// 交易金额
		domain.accountLog.Amount = domain.account.Balance
		domain.accountLog.Balance = domain.account.Balance
		// 交易变化属性
		domain.accountLog.Decs = "账户创建"
		domain.accountLog.ChangeType = services.AccountCreated
		domain.accountLog.ChangeFlag = services.FlagAccountCreated
	}
}

// 创建账户
func (domain *accountDomain) Create(dto services.AccountDTO) (*services.AccountDTO, error) {
	// 创建账户持久化对象
	domain.account = Account{}
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	domain.account.Username.Valid = true
	// 创建账户流水持久化对象
	domain.accountLog = AccountLog{}
	domain.createAccountLog()

	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	// 插入数据库数据
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner
		id, err := accountDao.Insert(&domain.account)

		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户失败")
		}
		id, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户流水失败")
		}
		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		return nil
	})
	retDto := domain.account.ToDTO()

	return retDto, err
}
func (domain *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferStatus, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		status, err = domain.TransferWithContext(ctx, dto)
		return err
	})
	return status, err
}

// 转账/充值
func (domain *accountDomain) TransferWithContext(ctx context.Context, dto services.AccountTransferDTO) (status services.TransferStatus, err error) {
	// 修正 amount
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagAccountOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}
	// 创建账户流水记录
	domain.accountLog = AccountLog{}
	domain.accountLog.FromTransferDTO(&dto)
	domain.createAccountLog()
	// 检查余额是否足够和更新余额:通过乐观锁来验证,更新余额的同时验证余额是否足够你
	// 写入流水记录
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		accountLogDao := AccountLogDao{runner: runner}
		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil {
			status = services.TransferStatusFailure
			return err
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagAccountOut {
			status = services.TransferStatusInsufficient
			return errors.New("余额不足")
		}
		a := accountDao.GetOne(dto.TradeBody.AccountNo)
		if a == nil {
			return errors.New("账户不存在")
		}
		domain.account = *a
		domain.accountLog.Balance = domain.account.Balance
		id, err := accountLogDao.Insert(&domain.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferStatusFailure
			return errors.New("账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	} else {
		status = services.TransferStatusSuccess
	}
	return status, err
}

// 根据账号查询账户
func (domain *accountDomain) GetAccount(accountNo string) *services.AccountDTO {
	var account *Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := &AccountDao{
			runner,
		}
		account = accountDao.GetOne(accountNo)
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据用户 ID 查询账户信息
func (domain *accountDomain) GetAccountByUserId(userId string) *services.AccountDTO {
	var account *Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{
			runner,
		}
		account = accountDao.GetByUserId(userId, int8(services.EnvelopeAccountType))

		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据流水ID来查询账户流水
func (a *accountDomain) GetAccountLog(logNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetOne(logNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}

// 根据交易编号来查询账户流水
func (a *accountDomain) GetAccountLogByTradeNo(tradeNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetByTradeNo(tradeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}
