package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	rzdFeatures "backend/app/feature/rzd"
	transactionFeature "backend/app/feature/transaction"
	userFeature "backend/app/feature/user"
	testfuncs "backend/app/test-funcs"
	cfg "backend/config"
	cntx "backend/context"
	db "backend/database"
	logg "backend/logger"
)

// RabbitMQ, Kafka, NATS
// redis
// grpc protobuff
// clickhouse
// JAEGER
// sentry
// grafana
// elk search
// tgbot
// todo GORM для конекшена БД
// todo add logger
// todo add прерывания операций по timeout и cancel
// todo подрубить контекст
// todo сделать DTO
// todo разложить по папочкам красиво все
// todo add auth
// подрубить ормку GORM
// todo сделать нормальные тесты

func startTracing() (*trace.TracerProvider, error) {
	serviceName := "product-app"
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("t-jaeger:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("environment", "testing"),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)

	return tracerprovider, nil
}

func main() {
	logger := logg.NewLogger()
	config := cfg.GetConfig(logger)
	fmt.Println("Type of logger:", reflect.TypeOf(logger))
	fmt.Println("Type of config:", reflect.TypeOf(config))

	dbConn := db.Conn(config)
	// defer dbConn.Close()

	if err := db.RunMigrations(logger, dbConn); err != nil {
		logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("migration failed: %v", err)})
	}
	// else {
	// 	dbConn.Close()
	// }

	// gin logging into file
	// err := os.MkdirAll("./log", 0777)
	// if err != nil {
	// 	logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("log file created failed: %v", err)})
	// }
	// logFile, err := os.Create("./log/gin.log" + time.Local.String())
	// if err != nil {
	// 	logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("log file created failed: %v", err)})
	// }
	// gin.DefaultWriter = io.MultiWriter(logFile)
	traceProvider, err := startTracing()
	if err != nil {
		log.Fatalf("traceprovider: %v", err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("traceprovider: %v", err)
		}
	}()

	tracerJaeger := traceProvider.Tracer("my-app")
	fmt.Println("[tracerJaeger]")
	fmt.Println(tracerJaeger)
	// run gin app
	router := gin.Default()
	router.Use(cntx.ContextMiddleware(config, logger))
	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	router.Use(gin.Recovery()) // panic recover middleware
	router.GET("/hi", func(gc *gin.Context) {
		var tracer = otel.Tracer("github.com/Salaton/tracing/pkg/usecases/product")
		fmt.Println("[tracer]")
		fmt.Println(tracer)
		_, span := tracer.Start(gc, "CreateProduct")
		defer span.End()
		span.SpanContext()
		span.AddEvent("event")
		span.SetAttributes(attribute.KeyValue{Key: "key"})
		span.SetStatus(codes.Ok, "ok status")
		span.SetName("name setName")

		fmt.Println("[span]")
		fmt.Println(span)
		span.AddEvent("after fmt")
		gc.String(http.StatusOK, "Holla! Welcome Gin Serer!\n")
		span.AddEvent("answer")
	})
	router.GET("/panic", func(gc *gin.Context) { panic("panic") })
	router.GET("/long", func(gc *gin.Context) { testfuncs.LongOperation(gc) })
	router.GET("/user-balance", func(gc *gin.Context) { userFeature.HandleUserBalance(logger, config, gc) })
	router.GET("/transactions", func(gc *gin.Context) { transactionFeature.HandleUserTransactions(logger, config, gc) })
	router.POST("/transactions", func(gc *gin.Context) { transactionFeature.HandleCreateTransaction(logger, config, gc) })
	router.POST("/rzd-stations", func(gc *gin.Context) { rzdFeatures.UploadStations(logger, config, gc) })
	// валидация, параметры (от до станции по коду)
	router.GET("/rzd-stations", func(gc *gin.Context) { rzdFeatures.GetStations(logger, config, gc) })
	// router.Run(config.Host)
	server := &http.Server{
		Addr:    config.Host,
		Handler: router.Handler(),
		// WriteTimeout: 3 * time.Second,
		// ReadTimeout:  30 * time.Second,
	}
	// _ = endless.ListenAndServe(config.Host, server.Handler) // facebok библа для gracfulshutdown
	go func() {
		logger.OuteputLog(logg.LogPayload{Info: "Starting server..."})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("error listening: %s", err)})
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM, kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	// docker exec test-spb-go-bank-bank-go-app-1 kill -SIGTERM 992 (на процесс main)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-timeoutCtx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")

	//  пример проверки завершения запроса
	/*
		r := gin.Default()
		// Middleware с таймаутом
		r.Use(func(c *gin.Context) {
			// Создаем контекст с таймаутом
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			// Заменяем контекст запроса на новый с таймаутом
			c.Request = c.Request.WithContext(ctx)
			// Передаем управление следующему middleware или обработчику
			c.Next()
		})
		// Обработчик, который может выполняться долго
		r.GET("/long-operation", func(c *gin.Context) {
			// Симуляция долгой операции
			select {
			case <-time.After(5 * time.Second):
				c.JSON(http.StatusOK, gin.H{"message": "Операция завершена"})
			case <-c.Request.Context().Done():
				// Если контекст отменен (например, по таймауту)
				c.JSON(http.StatusRequestTimeout, gin.H{"error": "Операция прервана по таймауту"})
			}
		})
		r.Run(":8080")
	*/
}
