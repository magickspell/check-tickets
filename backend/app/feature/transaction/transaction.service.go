package featureTransaction

import (
	"fmt"

	user "backend/app/feature/user"
	cfg "backend/config"

	// cntx "backend/context"
	logg "backend/logger"

	"github.com/gin-gonic/gin"
)

type transactionRequest struct {
	Amount     int  `json:"amount"`
	UserFromId *int `json:"userFromId"`
	UserToId   int  `json:"userToId"`
}

func GetTransactions(logger *logg.Logger, config *cfg.Config, userId string) ([]Transaction, error) {
	return GetUserTransactions(logger, config, userId)
}

func CreateTransaction(logger *logg.Logger, config *cfg.Config, gc *gin.Context) error {
	var req transactionRequest
	if err := gc.BindJSON(&req); err != nil {
		return fmt.Errorf("error: cant process transaction request")
	}

	if req.UserToId == 0 || req.Amount == 0 {
		return fmt.Errorf("error: some properties wasnot provided")
	}

	_, err := user.GetUserBalance(logger, config, req.UserToId)
	if err != nil {
		return fmt.Errorf("cant find user [UserToId]: '%v'", err)
	}

	if req.UserFromId != nil {
		userFrom, err := user.GetUserBalance(logger, config, *req.UserFromId)
		if err != nil {
			return fmt.Errorf("cant find user [UserFromId]: '%v'", err)
		}
		if userFrom.Balance < req.Amount {
			return fmt.Errorf("баланс меньше суммы перевода")
		}
	}

	err = InsertTransaction(logger, config, req.Amount, req.UserFromId, req.UserToId)
	if err != nil {
		return fmt.Errorf("cant InsertTransaction: '%v'", err)
	}
	return nil
}
