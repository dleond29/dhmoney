package cards

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

func TestRepositorySaveSuccesfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO cards ")
	mock.ExpectExec("INSERT INTO cards").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(db)
	card := domain.CardDto{
		PAN:            "8543345689903456",
		HolderName:     "Marco Suarez",
		ExpirationDate: "2022-12-10",
		CID:            "567",
	}

	us, err := repo.SaveCard(context.Background(), 1, card)
	assert.NoError(t, err)
	assert.NotZero(t, us)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetAllSuccesfully(t *testing.T) {
	input := []domain.Card{
		{
			ID:             1,
			AccountID:      1,
			PAN:            "123456",
			HolderName:     "Miguel",
			ExpirationDate: "2022-10-10",
			CID:            "123",
			Type:           "Visa",
		},
		{
			ID:             2,
			AccountID:      1,
			PAN:            "987643",
			HolderName:     "Marcos",
			ExpirationDate: "2022-02-10",
			CID:            "456",
			Type:           "MasterCard",
		},
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.TODO()

	columns := []string{"card_id", "account_id", "pan", "holder_name", "expiration_date", "cvv", "type"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(input[0].ID, input[0].AccountID, input[0].PAN, input[0].HolderName, input[0].ExpirationDate, input[0].CID, input[0].Type).AddRow(input[1].ID, input[1].AccountID, input[1].PAN, input[1].HolderName, input[1].ExpirationDate, input[1].CID, input[1].Type)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM cards WHERE account_id = ?;")).WithArgs(1).WillReturnRows(rows)

	repo := NewRepository(db)

	cards, err := repo.GetAll(ctx, 1)
	assert.NoError(t, err)
	assert.NotZero(t, cards)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryGetByIDSuccesfully(t *testing.T) {
	input := domain.Card{
		ID:             1,
		AccountID:      1,
		PAN:            "123456",
		HolderName:     "Miguel",
		ExpirationDate: "2022-10-10",
		CID:            "123",
		Type:           "Visa",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.TODO()

	columns := []string{"card_id", "account_id", "pan", "holder_name", "expiration_date", "cvv", "type"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(input.ID, input.AccountID, input.PAN, input.HolderName, input.ExpirationDate, input.CID, input.Type)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM cards WHERE id = ? and account_id = ?;")).WithArgs(1, 1).WillReturnRows(rows)

	repo := NewRepository(db)
	card, err := repo.GetByID(ctx, 1, 1)
	assert.NoError(t, err)
	assert.NotZero(t, card)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryExistsSuccesfully(t *testing.T) {
	input := domain.Card{
		ID:             1,
		AccountID:      1,
		PAN:            "123456",
		HolderName:     "Miguel",
		ExpirationDate: "2022-10-10",
		CID:            "123",
		Type:           "Visa",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.TODO()

	columns := []string{"card_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(input.ID)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM cards WHERE pan = ?;")).WithArgs("123456").WillReturnRows(rows)

	repo := NewRepository(db)
	result, err := repo.Exists(ctx, "123456")
	assert.NoError(t, err)
	assert.NotZero(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryDeleteSuccesfully(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.TODO()

	mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM cards WHERE id=?;")).ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewRepository(db)
	err = repo.DeleteByCardID(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
