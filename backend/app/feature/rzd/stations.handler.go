package rzd

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"sync"

	cfg "backend/config"
	// cntx "backend/context"

	logg "backend/logger"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TrainStation struct {
	StationID   string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();"`
	StationName string `gorm:"type:varchar(100)"`
	StationCode string `gorm:"type:varchar(50)"`
}

// todo сделать DTO
// todo унести куда надо все
// todo GORM для конекшена БД
func UploadStations(logger *logg.Logger, config *cfg.Config, gc *gin.Context) {
	// Шаг 1: Скачать ZIP файл
	zipURL := "https://support.travelpayouts.com/hc/article_attachments/360031345731/tutu_routes.csv.zip"
	resp, err := http.Get(zipURL)
	if err != nil {
		fmt.Println("Ошибка при скачивании файла:", err)
		return
	}
	defer resp.Body.Close()

	// Чтение содержимого ZIP файла в память
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении содержимого ZIP файла:", err)
		return
	}

	// Шаг 2: Распаковать ZIP файл
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		fmt.Println("Ошибка при распаковке ZIP файла:", err)
		return
	}

	var csvFile io.ReadCloser
	for _, file := range zipReader.File {
		if file.Name == "tutu_routes.csv" {
			csvFile, err = file.Open()
			if err != nil {
				fmt.Println("Ошибка при открытии CSV файла:", err)
				return
			}
			defer csvFile.Close()
			break
		}
	}

	// Шаг 3: Прочитать CSV файл
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Ошибка при чтении CSV файла:", err)
		return
	}

	// Шаг 4: Подключение к базе данных PostgreSQL
	dsn := config.DbURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return
	}

	// Шаг 5: Сохранение данных в базу данных
	const workersCount int = 8
	// делаем канал
	newRecords := make(chan TrainStation, workersCount)
	// делаем вэитгруппу
	var wg sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go workerSaver(newRecords, &wg, db)
	}

	for _, record := range records[1:] { // Пропускаем заголовок\
		departureId := record[0]
		departureName := record[1]
		arrivalId := record[2]
		arrivalName := record[3]

		// db.Create(&TrainStation{StationCode: arrivalId, StationName: arrivalName})
		// db.Create(&TrainStation{StationCode: departureId, StationName: departureName})
		newRecords <- TrainStation{StationCode: departureId, StationName: departureName}
		newRecords <- TrainStation{StationCode: arrivalId, StationName: arrivalName}
	}
	close(newRecords)
	// Ждем завершения всех воркеров
	wg.Wait()
	fmt.Println("Данные успешно сохранены в базу данных.")
}

func workerSaver(recordsToSave <-chan TrainStation, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()

	for newRecord := range recordsToSave {
		if err := db.Create(&newRecord).Error; err != nil {
			fmt.Println("Ошибка при сохранении:", err)
		}
	}
}

func GetStations(logger *logg.Logger, config *cfg.Config, gc *gin.Context) {
	dsn := config.DbURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("ошибка при подключении к базе данных: ", err)})
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при подключении к базе данных"})
		return
	}

	var stations []TrainStation
	if err := db.Find(&stations).Error; err != nil {
		logger.OuteputLog(logg.LogPayload{Error: fmt.Errorf("ошибка при получении станций:", err)})
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении станций"})
		return
	}

	gc.JSON(http.StatusOK, stations)
}
