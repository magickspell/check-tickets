package featureTransaction

import (
	"database/sql"
	"net/http"

	cfg "backend/config"
	// cntx "backend/context"
	logg "backend/logger"

	"github.com/gin-gonic/gin"
)

func HandleUserTransactions(logger *logg.Logger, config *cfg.Config, gc *gin.Context) {
	transactions, err := GetTransactions(logger, config, gc.Query("userId"))
	if err != nil {
		if err == sql.ErrNoRows {
			gc.JSON(http.StatusNotFound, gin.H{"messege": "user not found"})
		} else {
			gc.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}

func HandleCreateTransaction(logger *logg.Logger, config *cfg.Config, gc *gin.Context) {
	err := CreateTransaction(logger, config, gc)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, gin.H{"messege": "transaction complete"})
}
