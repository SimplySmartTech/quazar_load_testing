package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"simply_smart_mqtt_load_testing/pkg/model"

	_ "github.com/lib/pq"
)

var db *sql.DB

type DatabaseOps interface {
	FetchThingies(ctx context.Context, limit int32, sql string) (model.Thingies, error)
	TotalGateway(ctx context.Context, meterpergateway int32) (total int32)
	FetchCountMeterReading(ctx context.Context, tabel string) (int, error)
	FetchTemplateKey(ctx context.Context) ([]string, error)
	GetTemplateById(ctx context.Context, id int) ([]string, error)
	RemoveTemplateThingData(ctx context.Context, templateKey string) (int, error)
}

type databaseOps struct {
	db *sql.DB
}

func InitDatabase() {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	// open database
	fmt.Println(psqlconn, "url")
	dbconn, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	dbconn.SetMaxOpenConns(60)
	db = dbconn
	log.Println("Databases initialized successfully")
}

func NewDatabase() DatabaseOps {
	return &databaseOps{
		db: db,
	}
}
