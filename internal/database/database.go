package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToDB(connectionString string) (dbQueries *Queries, err error) {
	fmt.Println("connectionString: ", connectionString)
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

func attemptOpenConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
