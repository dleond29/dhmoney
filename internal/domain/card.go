package domain

type Card struct {
	ID             int    `json:"card_id"`
	AccountID      int    `json:"account_id"`
	PAN            string `json:"pan"`
	HolderName     string `json:"holder_name"`
	ExpirationDate string `json:"expiration_date"`
	CID            string `json:"cvv"`
	Type           string `json:"type"`
}

type CardDto struct {
	PAN            string `json:"pan"`
	HolderName     string `json:"holder_name"`
	ExpirationDate string `json:"expiration_date"`
	CID            string `json:"cvv"`
	Type           string `json:"type"`
}
