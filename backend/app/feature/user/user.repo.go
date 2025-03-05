package featureUser

import (
	cfg "backend/config"
	db "backend/database"
	logg "backend/logger"
	"database/sql"
	"fmt"
	"os"
)

type User struct {
	Id      int
	Balance int
}

type Operation string

const (
	OPERATION_PLUS  Operation = "plus"
	OPERATION_MINUS Operation = "minus"
)

func GetUser(logger *logg.Logger, config *cfg.Config, userId int) (User, error) {
	fmt.Println("[GetUser][START]")
	dbConn := db.Conn(config)
	defer dbConn.Close()

	var id, balance int

	err := dbConn.QueryRow("SELECT * FROM users u WHERE u.id=$1;", userId).Scan(&id, &balance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return User{}, err
	}

	return User{
		Id:      id,
		Balance: balance,
	}, nil
}

func UpdateUserBalance(logger *logg.Logger, config *cfg.Config, userId int, amount int, operation Operation, tran *sql.Tx) error {
	dbConn := db.Conn(config)
	defer dbConn.Close()

	var query string
	if operation == OPERATION_PLUS {
		query = "UPDATE users SET balance = balance + $1 WHERE id = $2"
	}
	if operation == OPERATION_MINUS {
		query = "UPDATE users SET balance = balance - $1 WHERE id = $2"
	}

	if len(query) == 0 {
		return fmt.Errorf("query is empty")
	}

	_, err := tran.Exec(query, amount, userId)
	if err != nil {
		return fmt.Errorf("cant process transaction: '%v'", err)
	}

	return nil
}
