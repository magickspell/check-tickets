package testfuncs

import (
	"fmt"
	syslog "log"
	"time"

	cntx "backend/context"
	logg "backend/logger"

	"github.com/gin-gonic/gin"
)

func LongOperation(gc *gin.Context) string {
	isCancelled, err := processLongOperation(gc)
	if err != nil {
		return "long bad"
	}

	if isCancelled {
		/**/
		fmt.Println("Long NOT complete (isCancelled)")
		return "long bad"
	} else {
		fmt.Println("Long complete (isCancelled)")
		return "long ok"
	}
}

func processLongOperation(gc *gin.Context) (bool, error) {
	fmt.Println("[GC]")
	fmt.Println("[gc][pointer]")
	fmt.Println(gc)
	fmt.Println("[gc][value]")
	fmt.Println(&gc)

	fmt.Println("[gc][config]")
	fmt.Println(gc.Value("Config"))
	fmt.Println(gc.Value("Timeout"))
	fmt.Println(gc.Value("Logger"))
	fmt.Println(gc.Value("IsCancelled"))
	fmt.Println(gc.Value("ctx"))

	context, err := getAppContext(gc)
	if err != nil {
		syslog.Println("[error occured]")
		syslog.Println(err.Error())
		panic(err.Error())
	}

	// Теперь ctx имеет тип *Context, и вы можете работать с ним
	fmt.Println("[ctx.Timeout]")
	fmt.Println(context.Timeout)
	fmt.Println("[ctx.Config]")
	fmt.Println(context.Config)
	fmt.Println(&context.Config)
	fmt.Println("[ctx.Logger]")
	fmt.Println(context.Logger)
	fmt.Println("[ctx.IsCancelled]")
	fmt.Println(context.IsCancelled)

	isCancelled := &context.IsCancelled

	context.Logger.OuteputLog(logg.LogPayload{Info: "CONTEXT GRANTED"})
	// logger := context.Logger
	// logger.OuteputLog(logg.LogPayload{Info: "CONTEXT GRANTED"})

	fmt.Println("[ctx]")
	fmt.Println(context)
	fmt.Println(&context)

	time.Sleep(time.Second * 7)
	fmt.Println("isCancelled = ", isCancelled)
	fmt.Println("*isCancelled = ", *isCancelled)
	fmt.Println("&isCancelled = ", &isCancelled)
	if *isCancelled {
		/**/
		fmt.Println("Long NOT complete (isCancelled)")
		return *isCancelled, nil
	} else {
		fmt.Println("Long complete (isCancelled)")
		return *isCancelled, nil
	}
}

// унести отдельно в хелперы
func getAppContext(gc *gin.Context) (*cntx.Context, error) {
	// обязательно копируем контекст (гин дектиует правило)
	gccp := gc.Copy()
	// получаем свой контекст приложения
	appCtx, ex := gccp.Get("ctx")
	if !ex {
		err := fmt.Errorf("no context provided")
		syslog.Println(err.Error())
		return nil, err
	}
	// Приведеним типа контекстаы
	context, ok := appCtx.(*cntx.Context)
	if !ok {
		err := fmt.Errorf("приведение типа не удалось, unknown не является *cntx.Context")
		syslog.Println(err.Error())
		return nil, err
	}
	return context, nil
}
