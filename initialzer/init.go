package initialzer

import (
	"fmt"
	"log"
	"project/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitSetup() {
	DbInit()
}
func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{},&models.Admin{})
}

func DbInit() {
	var err error
	DSN := "host=localhost user=postgres password=7009 dbname=test port=5432"
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	autoMigrate(db)
	DB = db
	fmt.Println("============================ CONNECTED TO DB =====================================")
}
func MockDbConfig(t *testing.T) *sqlmock.Sqlmock {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: mockDB,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	DB = db
	return &mock
}
