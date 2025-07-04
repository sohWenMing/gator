package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

// ############################ Connections ###################################//

func ConnectToDB(connectionString string) (dbQueries *Queries, err error) {
	for i := range 60 {
		db, err := attemptOpenConnection(connectionString)
		if err != nil {
			fmt.Println("attempt to connect to database failed: ", i)
		} else {
			dbQueries := New(db)
			return dbQueries, nil
		}
	}
	return nil, errors.New("db did not connect in time")
}

// this function will try continually to connect to the database, but
// will fail after 60 times

func attemptOpenConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
