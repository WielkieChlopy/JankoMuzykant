package testutils

import (
	"encoding/json"

	"backend/api/v1"
	"backend/api/v1/user"
	"backend/db"
	"backend/models"
	"backend/router"
	"backend/store"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Test struct {
	Database  *sqlx.DB
	UserStore *store.UserStore
	SongStore *store.SongStore
	Handler   *v1.Handler
	Router    *echo.Echo
}

func SetupTest() (*Test, error) {
	t := Test{}

	db, err := db.NewTestDb()
	if err != nil {
		return nil, err
	}

	t.Database = db
	t.UserStore = store.NewUserStore(db)
	t.SongStore = store.NewSongStore(db)
	h, err := v1.NewHandler(t.UserStore, t.SongStore)
	if err != nil {
		return nil, err
	}
	t.Handler = h
	e := router.New()
	t.Router = e
	t.loadFixtures()

	return &t, nil
}

func ResponseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func (t Test) loadFixtures() error {
	u1 := &models.User{
		Username: "user1",
	}
	u1.Password, _ = user.HashPassword("secret")
	u1, err := t.UserStore.Create(u1)
	if err != nil {
		return err
	}

	u2 := &models.User{
		Username: "user2",
	}
	u2.Password, _ = user.HashPassword("secret")
	u2, err = t.UserStore.Create(u2)
	if err != nil {
		return err
	}

	return nil
}
