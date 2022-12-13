package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID             int
	Account        Account
	DestinationCVU string
	Description    string
	Amount         decimal.Decimal
	DateTime       time.Time
	Type           string
}

type TransactionInfo struct {
	ID             int             `json:"transaction_id"`
	AccountID      int             `json:"account_id"`
	OriginCVU      string          `json:"origin_cvu"`
	DestinationCVU string          `json:"destination_cvu"`
	Description    string          `json:"description"`
	Amount         decimal.Decimal `json:"amount"`
	DateTime       time.Time       `json:"date_time"`
	Type           string          `json:"type"`
}
