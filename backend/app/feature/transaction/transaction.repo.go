package featureTransaction

import (
	user "backend/app/feature/user"
	cfg "backend/config"
	db "backend/database"
	"fmt"

	// cntx "backend/context"
	logg "backend/logger"
)

type Transaction struct {
	Id        int
	Amount    int
	ToUser    int
	FromUser  *int
	CreatedAt string
}

func GetUserTransactions(logger *logg.Logger, config *cfg.Config, userId string) ([]Transaction, error) {
	dbConn := db.Conn(config)
	defer dbConn.Close()

	rows, err := dbConn.Query(
		`
		SELECT * FROM transactions t
		WHERE t.user_to=$1
			OR t.user_from=$1
		ORDER BY created_at desc
		LIMIT 10
		;
		`,
		userId,
	)
	if err != nil {
		fmt.Printf("query failed: '%v'", err)
		return []Transaction{}, fmt.Errorf("query failed: '%w'", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.Id, &transaction.ToUser, &transaction.FromUser, &transaction.Amount, &transaction.CreatedAt)
		if err != nil {
			fmt.Printf("scanning row failde: '%v'", err)
			return nil, fmt.Errorf("scanning row failde: '%w'", err)
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("transactions for user not found")
	}
	return transactions, nil
}

func InsertTransaction(logger *logg.Logger, config *cfg.Config, amount int, userFrom *int, userTo int) error {
	dbConn := db.Conn(config)
	defer dbConn.Close()

	// стартуем транзакцию по переводу денег от одного пользователя к другому
	tran, err := dbConn.Begin()
	if err != nil {
		return fmt.Errorf("cant start transaction: '%w'", err)
	}

	if userFrom != nil {
		err = user.UpdateUserBalance(logger, config, *userFrom, amount, user.OPERATION_MINUS, tran)
		if err != nil {
			return fmt.Errorf("cant OPERATION_MINUS: '%w'", err)
		}
	}

	err = user.UpdateUserBalance(logger, config, userTo, amount, user.OPERATION_PLUS, tran)
	if err != nil {
		return fmt.Errorf("cant OPERATION_PLUS: '%w'", err)
	}

	_, err = tran.Exec(
		"INSERT INTO transactions (amount, user_from, user_to) VALUES ($1, $2, $3);",
		amount, userFrom, userTo,
	)
	if err != nil {
		return fmt.Errorf("unable to insert row: '%w'", err)
	}

	// комитим транзакцию
	err = tran.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit row: '%w'", err)
	}

	return nil
}
