package repository

import (
	"database/sql"

	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/subbbbbaru/first-sample/pkg/log"
)

const (
	todoTaskTable = "scheduler"
)

// Создание новой базы данных или подключение к существующей БД
func NewSQLite3DB(dbName string) (*sql.DB, error) {

	if !checkDB(dbName) {
		if err := createDB(dbName); err != nil {
			return nil, err
		}
	}

	db, err := openDB(dbName)
	if err != nil {
		return nil, err
	}
	if err = createTable(db); err != nil {
		return nil, err
	}
	return db, nil // &Database{db: db}, nil
}
func openDB(dbName string) (*sql.DB, error) {
	log.Info().Println("Open database file")
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "db", dbName)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	return db, err
}

// Создание файла базы данных
func createDB(dbName string) error {
	log.Info().Println("Create database file")
	appPath, err := os.Executable()
	if err != nil {
		return err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "db", dbName)

	_, err = os.Create(dbFile)
	if err != nil {
		return err
	}
	return nil
}

// Проверка существования файла базы данных
func checkDB(dbName string) bool {
	log.Info().Println("Check exist database file")
	appPath, err := os.Executable()
	if err != nil {
		log.Error().Fatal(err)
		return false
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "db", dbName)

	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		log.Error().Println("Database file not exist")
		return false
	}
	log.Error().Println("Database file exist")
	return true
}

// createTable создает таблицу `scheduler` в базе данных
func createTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
            date VARCHAR(8) NOT NULL,
            comment TEXT,
            repeat VARCHAR(128)
        );
        CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
    `)
	return err
}
