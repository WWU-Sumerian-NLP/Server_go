package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type InternalDB struct {
	mu sync.Mutex
	db *sql.DB
}

func NewInternalDB() (*InternalDB, error) {
	db, err := sql.Open("sqlite3", "internalDB.sqlt")
	fmt.Printf("db: %v\n", db)
	if err != nil {
		log.Fatal(err)
	}

	return &InternalDB{
		db: db,
	}, nil
}

func (i *InternalDB) CloseDB() {
	i.db.Close()
}
