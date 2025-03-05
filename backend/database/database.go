package database

import (
	"database/sql"
	"fmt"
	"os"

	config "backend/config"
)

func Conn(c *config.Config) *sql.DB {
	// todo env должно быть только в main.go
	dbconn, err := sql.Open("pgx", c.DbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbconn
}
