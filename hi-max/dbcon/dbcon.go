package dbcon

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", os.Getenv("DB_INFO"))
		if err != nil {
			log.Fatal(err.Error())
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	return db
}
