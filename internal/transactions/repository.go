package transactions

import (
	"context"
	"database/sql"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

type Repository interface {
	GetAllByIDLimit(ctx context.Context, id, limit int) ([]domain.TransactionInfo, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllByIDLimit(ctx context.Context, id, limit int) ([]domain.TransactionInfo, error) {
	query := "SELECT * FROM transactions WHERE account_id = ? LIMIT ?;"
	rows, err := r.db.Query(query, id, limit)
	if err != nil {
		return []domain.TransactionInfo{}, err
	}

	var transactions []domain.TransactionInfo

	for rows.Next() {
		trx := domain.TransactionInfo{}
		err = rows.Scan(&trx.ID, &trx.AccountID, &trx.DestinationCVU, &trx.Description, &trx.Amount, &trx.DateTime, &trx.Type, &trx.OriginCVU)
		if err != nil {
			return []domain.TransactionInfo{}, err
		}

		transactions = append(transactions, trx)
	}

	return transactions, nil
}
