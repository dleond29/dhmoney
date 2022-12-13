package accounts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/internal/users"
)

type Repository interface {
	SaveAccount(ctx context.Context, accountDto domain.AccountDto) (*users.UserDto, error)
	GetAccountByID(ctx context.Context, id int) (domain.Account, error)
	GetAccountByUserID(ctx context.Context, userID int) (domain.Account, error)
	CVUExist(ctx context.Context, cvu string) bool
	AliasExist(ctx context.Context, alias string) bool
	IsAuthorized(ctx context.Context, accountID int, isUserID bool, authID string) (bool, error)
	UpdateAlias(ctx context.Context, accountID int, alias string) error
}

type repository struct {
	db *sql.DB
}

type accountDB struct {
	userID  int
	authID  string
	cvu     string
	alias   string
	balance decimal.Decimal
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAccountByID(ctx context.Context, id int) (domain.Account, error) {
	query := "SELECT * FROM accounts WHERE id = ?;"
	rows := r.db.QueryRow(query, id)

	var account domain.Account

	err := rows.Scan(&account.ID, &account.User.ID, &account.AuthID, &account.CVU, &account.Alias, &account.Balance)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return domain.Account{}, ErrAccountNotFound
		default:
			return domain.Account{}, err
		}
	}

	return account, nil
}

func (r *repository) GetAccountByUserID(ctx context.Context, userID int) (domain.Account, error) {
	query := "SELECT * FROM accounts WHERE user_id = ?;"
	rows := r.db.QueryRow(query, userID)

	var account domain.Account

	err := rows.Scan(&account.ID, &account.User.ID, &account.AuthID, &account.CVU, &account.Alias, &account.Balance)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return domain.Account{}, ErrAccountNotFound
		default:
			return domain.Account{}, err
		}
	}

	return account, nil
}

func (r *repository) IsAuthorized(ctx context.Context, accountID int, isUserID bool, authID string) (bool, error) {
	var idName string
	if isUserID {
		idName = "user_id"
	} else {
		idName = "id"
	}
	query := fmt.Sprintf("SELECT id FROM accounts WHERE %s = ? and auth_id = ?;", idName)
	row := r.db.QueryRow(query, accountID, authID)

	var id int

	err := row.Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}

	return id != 0, nil
}

func Save(ctx context.Context, tx *sql.Tx, account accountDB) (int, error) {
	query := "INSERT INTO accounts(user_id, auth_id, cvu, alias, balance) VALUES(?, ?, ?, ?, 0)"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(account.userID, account.authID, account.cvu, account.alias)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) SaveAccount(ctx context.Context, accountDto domain.AccountDto) (*users.UserDto, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &users.UserDto{}, err
	}

	userDB := domain.UserDB{
		DNI:   accountDto.DNI,
		Phone: accountDto.Phone,
	}
	userID, err := users.Save(ctx, tx, userDB)
	if err != nil {
		return &users.UserDto{}, err
	}

	accountToSave := accountDB{
		userID: userID,
		authID: accountDto.AuthID,
		cvu:    accountDto.CVU,
		alias:  accountDto.Alias,
	}

	_, err = Save(ctx, tx, accountToSave)
	if err != nil {
		return &users.UserDto{}, err
	}

	err = tx.Commit()
	if err != nil {
		return &users.UserDto{}, err
	}

	account := users.UserDto{
		ID:    userID,
		DNI:   accountDto.DNI,
		Phone: accountDto.Phone,
		CVU:   accountDto.CVU,
		Alias: accountDto.Alias,
	}
	return &account, nil
}

func (r *repository) CVUExist(ctx context.Context, cvu string) bool {
	query := "SELECT id FROM accounts WHERE cvu = ?"
	row := r.db.QueryRow(query, cvu)

	var id int
	err := row.Scan(&id)

	return err == nil
}

func (r *repository) AliasExist(ctx context.Context, alias string) bool {
	query := "SELECT id FROM accounts WHERE alias = ?"
	row := r.db.QueryRow(query, alias)

	var id int
	err := row.Scan(&id)

	return err == nil
}

func (r *repository) UpdateAlias(ctx context.Context, accountID int, alias string) error {
	query := "UPDATE accounts SET alias = ? WHERE id = ?;"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&alias, &accountID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
