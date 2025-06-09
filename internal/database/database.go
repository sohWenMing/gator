package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func ConnectToDB(connectionString string) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
