package conn

import "github.com/go-pg/pg"

var (
	db *pg.DB
)

func Connect() {
	db = pg.Connect(&pg.Options{
		User:     "lwadmin",
		Database: "lessworkflow",
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
