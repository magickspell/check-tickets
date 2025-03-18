package trash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cfg "backend/config"
	db "backend/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func CreateTrash(c *gin.Context, config *cfg.Config) {
	ctxt, _ := context.WithTimeout(c, time.Second*3)
	var tracer = otel.Tracer("trash.repo.go")
	_, span := tracer.Start(ctxt, "SendTrash")
	span.AddEvent("SendTrash init")
	defer span.End()

	dbConn := db.Conn(config)
	defer dbConn.Close()

	var trash Trash
	if err := c.ShouldBindJSON(&trash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("[CREATE][[trash]]")
	fmt.Println(trash)

	_, err := dbConn.ExecContext(
		ctxt,
		"INSERT INTO trashes (trash_id, trash_name, trash_code, trash_json) VALUES ($1, $2, $3, $4)",
		trash.TrashID,
		trash.TrashName,
		trash.TrashCode,
		trash.TrashJSON,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.AddEvent("SendTrash end")

	c.JSON(http.StatusCreated, trash)
}

func SendTrash(trash Trash) {
	var tracer = otel.Tracer("trash.service.go")
	_, span := tracer.Start(context.Background(), "SendTrash")
	defer span.End()

	span.SpanContext()
	span.AddEvent("SendTrash event")
	span.SetAttributes(attribute.KeyValue{Key: "key"})
	span.SetStatus(codes.Ok, "SendTrash ok status")
	span.SetName("SendTrash name setName")

	// todo from config
	url := "http://localhost:8080/trash"
	jsonData, err := json.Marshal(trash)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println("[jsonData]")
	fmt.Println(jsonData)
	fmt.Println(bytes.NewBuffer(jsonData))
	span.AddEvent("SendTrash json parsed")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	span.AddEvent("SendTrash json posted")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("Successfully created trash: %+v\n", trash)
	} else {
		fmt.Printf("Failed to create trash: %s\n", resp.Status)
	}
}

func GenerateUUID() string {
	// Простой генератор UUID (можно использовать библиотеку для более надежного генератора)
	return uuid.NewString()
}
