package featureUser

import (
	cfg "backend/config"
	logg "backend/logger"
	"fmt"
)

func GetUserBalance(logger *logg.Logger, config *cfg.Config, userId int) (User, error) {
	// todo GetUser должен принимать context первым аргументом (context.Context)
	fmt.Println("[START][GetUserBalance]")
	//time.Sleep(time.Second * 3)
	return GetUser(logger, config, userId)
}
