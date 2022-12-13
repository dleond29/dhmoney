package users

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/leorodriguez/grupo-04/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func Save(ctx context.Context, tx *sql.Tx, user domain.UserDB) (int, error) {
	query := "INSERT INTO users(dni, phone) VALUES(?, ?)"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(user.DNI, user.Phone)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (domain.UserDB, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var user domain.UserDB
	err := row.Scan(&user.ID, &user.DNI, &user.Phone)
	if err != nil {
		return domain.UserDB{}, err
	}

	return user, nil
}

func (r *Repository) CVUExist(ctx context.Context, cvu string) bool {
	query := "SELECT id FROM users WHERE cvu = ?"
	row := r.db.QueryRow(query, cvu)

	var id int
	err := row.Scan(&id)

	return err == nil
}

func (r *Repository) AliasExist(ctx context.Context, alias string) bool {
	query := "SELECT id FROM users WHERE alias = ?"
	row := r.db.QueryRow(query, alias)

	var id int
	err := row.Scan(&id)

	return err == nil
}

func (r *Repository) UpdateUser(ctx context.Context, accountDto domain.AccountDto, id int) error {
	var sets string
	var args []interface{}
	if accountDto.DNI != 0 {
		sets += "dni = ?"
		args = append(args, &accountDto.DNI)
	}
	if accountDto.Phone != 0 {
		if sets != "" {
			sets += ","
		}
		sets += "phone = ?"
		args = append(args, &accountDto.Phone)
	}
	args = append(args, &id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?;", sets)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
