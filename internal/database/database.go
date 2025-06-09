package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func ConnectToDB(connectionString string) (conn *pgx.Conn, err error) {
	for i := range 60 {
		fmt.Println("Attempting to connect to database, attempt: ", i)
		conn, err = pgx.Connect(context.Background(), connectionString)
		if err == nil {
			return conn, nil
		} else {
			time.Sleep(500 * time.Millisecond)
			continue
		}
	}
	return nil, errors.New("connection to database could not be established in time")
}
