package domain

import "github.com/shopspring/decimal"

type Account struct {
	ID      int
	AuthID  string
	User    User
	CVU     string
	Alias   string
	Balance decimal.Decimal
}

type AccountInfo struct {
	UserID    int             `json:"user_id"`
	AccountID int             `json:"account_id"`
	CVU       string          `json:"cvu"`
	Alias     string          `json:"alias"`
	Balance   decimal.Decimal `json:"balance"`
}

type AccountDto struct {
	UserID  int
	AuthID  string
	DNI     int
	Phone   int
	CVU     string
	Alias   string
	Balance decimal.Decimal
}

type Alias struct {
	Alias string `json:"alias"`
}
