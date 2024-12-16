package repository

import (
	"errors"
	"fmt"
	"github.com/Quszlet/user_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserMySql struct {
	db *sqlx.DB
}

func NewUserMySql(db *sqlx.DB) *UserMySql {
	return &UserMySql{db: db}
}

func (up *UserMySql) Create(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, email, password, status) values ($1, $2, $3, $4)", usersTable)

	row := up.db.QueryRow(query, user.Login, user.Email, user.Password, user.Is_banned)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// TODO
func (up *UserMySql) Update(userId int) error {
	return nil
}

func (up *UserMySql) Get(userId int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", usersTable)
	err := up.db.Get(&user, query, userId)
	return user, err
}

func (up *UserMySql) GetAll() ([]models.User, error) {
	var users []models.User
	query := fmt.Sprintf("SELECT * FROM %s", usersTable)
	err := up.db.Select(&users, query)
	return users, err
}

func (up *UserMySql) Delete(userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", usersTable)
	res, err := up.db.Exec(query, userId)
	if err != nil {
		return err
	}

	affRows, err := res.RowsAffected()
	if affRows == 0 {
		return errors.New("User with this ID does not exist")
	}
	
	return err
}