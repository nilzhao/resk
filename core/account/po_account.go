package account

import (
	"database/sql"
	"resk/services"
	"time"

	"github.com/shopspring/decimal"
)

// Account 账户持久化对象
type Account struct {
	Id           int64           `db:"id,omitempty"`
	AccountNo    string          `db:"account_no,uni"`
	AccountName  string          `db:"account_name"`
	AccountType  int8            `db:"account_type"`
	CurrencyCode string          `db:"currency_code"`
	UserId       string          `db:"user_id"`
	Username     sql.NullString  `db:"username"`
	Balance      decimal.Decimal `db:"balance"`
	Status       int8            `db:"status"`
	CreatedAt    time.Time       `db:"created_at,omitempty"`
	UpdatedAt    time.Time       `db:"updated_at,omitempty"`
}

//,omitempty
func (po *Account) ToDTO() *services.AccountDTO {
	dto := &services.AccountDTO{}
	dto.AccountNo = po.AccountNo
	dto.AccountName = po.AccountName
	dto.AccountType = po.AccountType
	dto.CurrencyCode = po.CurrencyCode
	dto.UserId = po.UserId
	dto.Username = po.Username.String
	dto.Balance = po.Balance
	dto.Status = po.Status
	dto.CreatedAt = po.CreatedAt
	dto.UpdatedAt = po.UpdatedAt
	return dto
}

func (po *Account) FromDTO(dto *services.AccountDTO) {
	po.AccountNo = dto.AccountNo
	po.AccountName = dto.AccountName
	po.AccountType = dto.AccountType
	po.CurrencyCode = dto.CurrencyCode
	po.UserId = dto.UserId
	po.Username = sql.NullString{Valid: true, String: dto.Username}
	po.Balance = dto.Balance
	po.Status = dto.Status
	po.CreatedAt = dto.CreatedAt
	po.UpdatedAt = dto.UpdatedAt
}
