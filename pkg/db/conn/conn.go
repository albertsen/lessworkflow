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
	db = pg.Connect(&pg.Options{
		User:     "lwadmin",
		Database: "lessworkflow",
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
