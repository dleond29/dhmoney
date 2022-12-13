package cards

import (
	"context"
	"database/sql"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

type Repository interface {
	SaveCard(ctx context.Context, id int, card domain.CardDto) (int, error)
	GetAll(ctx context.Context, accountID int) ([]domain.Card, error)
	GetByID(ctx context.Context, accountID, cardID int) (domain.Card, error)
	Exists(ctx context.Context, pan string) (bool, error)
	DeleteByCardID(ctx context.Context, cardID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveCard(ctx context.Context, id int, card domain.CardDto) (int, error) {
	query := "INSERT INTO cards (account_id, pan, holder_name, expiration_date, cid, type) VALUES (?, ?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&id, &card.PAN, &card.HolderName, &card.ExpirationDate, &card.CID, &card.Type)
	if err != nil {
		return 0, err
	}

	cardID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(cardID), nil
}

func (r *repository) GetAll(ctx context.Context, accountID int) ([]domain.Card, error) {
	query := "SELECT * FROM cards WHERE account_id = ?;"
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return []domain.Card{}, err
	}

	var cards []domain.Card

	for rows.Next() {
		card := domain.Card{}
		err = rows.Scan(&card.ID, &card.AccountID, &card.PAN, &card.HolderName, &card.ExpirationDate, &card.CID, &card.Type)
		if err != nil {
			return []domain.Card{}, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (r *repository) GetByID(ctx context.Context, accountID, cardID int) (domain.Card, error) {
	query := "SELECT * FROM cards WHERE id = ? and account_id = ?;"
	rows := r.db.QueryRow(query, cardID, accountID)

	var card domain.Card

	err := rows.Scan(&card.ID, &card.AccountID, &card.PAN, &card.HolderName, &card.ExpirationDate, &card.CID, &card.Type)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return domain.Card{}, ErrCardNotFound
		}
		return domain.Card{}, err
	}

	return card, nil
}

func (r *repository) Exists(ctx context.Context, pan string) (bool, error) {
	query := "SELECT id FROM cards WHERE pan = ?;"
	rows := r.db.QueryRow(query, pan)

	var id int

	err := rows.Scan(&id)
	if err != nil {
		// todo refactor
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}

	return id != 0, nil
}

func (r *repository) DeleteByCardID(ctx context.Context, cardID int) error {
	query := "DELETE FROM cards WHERE id=?;"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(cardID)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrCardNotFound
	}

	return nil
}
