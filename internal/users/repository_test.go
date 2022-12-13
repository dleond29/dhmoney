package users

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

/* -------------------------------------------------------------------------- */
/*                                 TEST CREATE                                */
/* -------------------------------------------------------------------------- */
/* ----------------------------- TEST Save_Ok ----------------------------- */
func TestRepositorySave(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO users")
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewRepository(db)
	user := domain.User{
		ID:     1,
		IDAuth: "id1",
		DNI:    12345678,
		Phone:  12345678,
		CVU:    "1111111111111111111111",
		Alias:  "casa.perro.pelota",
	}

	us, err := repo.Save(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, us)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ----------------------------- TEST Save_Conflict_Exec ----------------------------- */

func TestRepositorySaveConflictExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO users")
	mock.ExpectExec("INSERT INTO users").WillReturnError(errors.New(""))
	repo := NewRepository(db)
	user := domain.User{
		ID:     1,
		IDAuth: "id1",
		DNI:    12345678,
		Phone:  12345678,
		CVU:    "1111111111111111111111",
		Alias:  "casa.perro.pelota",
	}
	us, err := repo.Save(context.Background(), user)
	assert.Error(t, err)
	assert.Zero(t, us)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ----------------------------- TEST Save_Conflict_Prepare ----------------------------- */
func TestRepositorySaveErrPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	user := domain.User{
		ID:     1,
		IDAuth: "id1",
		DNI:    12345678,
		Phone:  12345678,
		CVU:    "1111111111111111111111",
		Alias:  "casa.perro.pelota",
	}
	mock.ExpectPrepare("INSERT INTO users").WillReturnError(errors.New("err"))

	repo := NewRepository(db)

	_, errSave := repo.Save(context.Background(), user)
	assert.NotNil(t, errSave)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* -------------------------------------------------------------------------- */
/*                                 TEST EXISTS                                */
/* -------------------------------------------------------------------------- */
/* ------------------------------- TEST Exists CVU------------------------------ */

func TestCVUExistsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id"}

	row := sqlmock.NewRows(column)

	//no tiene 22 digitos
	row.AddRow("1")

	cvu := "1111111111111111111111"

	mock.ExpectQuery("SELECT id FROM users WHERE cvu = ?").WithArgs(cvu).WillReturnRows(row)

	repo := NewRepository(db)
	boolean := repo.CVUExist(context.Background(), cvu)

	assert.True(t, boolean)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ------------------------------- TEST Exists Alias------------------------------ */

func TestAliasExistsOk(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id"}

	row := sqlmock.NewRows(column)

	row.AddRow("1")

	alias := "aa.aa.aa"
	mock.ExpectQuery("SELECT id FROM users WHERE alias = ?").WithArgs(alias).WillReturnRows(row)

	repo := NewRepository(db)
	boolean := repo.AliasExist(context.Background(), alias)

	assert.True(t, boolean)
	assert.NoError(t, mock.ExpectationsWereMet())
}
