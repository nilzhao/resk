package account

import (
	"errors"
	"resk/infra/base"
	"resk/services"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var _ services.AccountService = new(accountService)
var once sync.Once

func init() {
	once.Do(func() {
		services.IAccountService = new(accountService)
	})
}

type accountService struct{}

// CreateAccount 创建账户
func (s *accountService) CreateAccount(dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	// 验证输入参数
	err := base.Validate().Struct(dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("参数验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				logrus.Error(err.Translate(base.Translate()))
			}
		}
		return nil, err
	}
	// 执行账户创建的业务逻辑
	domain := accountDomain{}
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	accountDto := &services.AccountDTO{
		UserId:       dto.UserId,
		Username:     dto.Username,
		AccountName:  dto.AccountName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		Balance:      amount,
	}
	accountDto, err = domain.Create(*accountDto)

	return accountDto, err
}

// Transfer 转账
func (s *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferStatus, error) {
	// 验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("参数验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				logrus.Error(err.Translate(base.Translate()))
			}
		}
		return services.TransferStatusFailure, err
	}
	// 验证转账类型和 flag 是否一致
	if dto.ChangeFlag == services.FlagAccountOut {
		if dto.ChangeType > 0 {
			return services.TransferStatusFailure,
				errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferStatusFailure,
				errors.New("如果changeFlag为收入,那么changeType必须大于0")
		}
	}

	// 执行转账逻辑
	domain := accountDomain{}
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferStatusFailure, err
	}
	dto.Amount = amount
	status, err := domain.Transfer(dto)
	return status, err
}

// StoreValue 存钱
func (s *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferStatus, error) {
	dto.ChangeType = services.AccountStoreValue
	dto.ChangeFlag = services.FlagAccountIn
	dto.TradeTarget = dto.TradeBody
	return s.Transfer(dto)
}

// GetAccountByUserId 通过userId获取红包账户
func (s *accountService) GetAccountByUserId(userId string) *services.AccountDTO {
	domain := accountDomain{}
	accountDto := domain.GetAccountByUserId(userId)
	return accountDto
}

// 通过账户编号查询账户
func (s *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := accountDomain{}
	accountDto := domain.GetAccount(accountNo)
	return accountDto
}
