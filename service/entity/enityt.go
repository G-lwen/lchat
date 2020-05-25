package entity

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"lchat/service/conf"
)

const (
	UserStore = "user"
	PostStore = "post"
)

type Entity struct {
	gorm.Model
}

var store map[string]*gorm.DB

func NewStore() error {
	store = make(map[string]*gorm.DB)
	db := conf.Get().DB
	for i := 0; i < len(db); i++ {
		switch db[i].Name {
		case UserStore:
			newDB(db[i].Name, db[i].Dialect, db[i].DSN, &User{})
		case PostStore:
			newDB(db[i].Name, db[i].Dialect, db[i].DSN, &Post{}, &Tag{}, &PostTag{})
		default:

		}
	}
	return nil
}

func newDB(name, dialect, dsn string, values ...interface{}) error {
	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		return err
	}
	db.LogMode(true)
	db.AutoMigrate(values...)
	store[name] = db
	return nil
}

func Close() {
	for key, value := range store {
		value.Close()
		delete(store, key)
	}
}