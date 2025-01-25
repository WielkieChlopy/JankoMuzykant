package store

import (
	"backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us UserStore) GetByID(id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	if err := us.db.Get(u, "SELECT * FROM \"user\" where id=$1", id); err != nil {
		return nil, err
	}

	return u, nil
}

func (us UserStore) GetByUsername(username string) (*models.User, error) {
	u := &models.User{}
	if err := us.db.Get(u, "SELECT * FROM \"user\" where username=$1", username); err != nil {
		return nil, err
	}

	return u, nil
}

func (us UserStore) Create(u *models.User) (*models.User, error) {
	if _, err := us.db.Exec("insert into \"user\" (username, password) values ($1, $2)", u.Username, u.Password); err != nil {
		return nil, err
	}

	if err := us.db.Get(u, "SELECT * FROM \"user\" where username=$1", u.Username); err != nil {
		return nil, err
	}

	return u, nil

}

func (us UserStore) Update(u *models.User) (*models.User, error) {
	if _, err := us.db.Exec("update \"user\" set username=$2, password=$3 where id = $1", u.Id, u.Username, u.Password); err != nil {
		return nil, err
	}

	new_u := &models.User{}
	if err := us.db.Get(new_u, "SELECT * FROM \"user\" where id=$1", u.Id); err != nil {
		return nil, err
	}

	return new_u, nil
}

func (us UserStore) Delete(id uuid.UUID) error {
	if _, err := us.db.Exec("delete from \"user\" where id = $1", id); err != nil {
		return err
	}

	return nil
}
