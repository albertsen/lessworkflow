package conn

import (
	"os"

	"github.com/go-pg/pg"
)

var (
	db *pg.DB
)

func Connect() {
	addr := os.Getenv("DB_ADDR")
	if addr == "" {
		addr = "localhost:5432"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "lessworkflow"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "lwadmin"
	}
	password := os.Getenv("DB_PASSWORD")
	db = pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbName,
		Addr:     addr,
	})
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func DB() *pg.DB {
	return db
}
