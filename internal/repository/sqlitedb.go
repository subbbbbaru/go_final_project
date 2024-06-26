package repository

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	todoTaskTable = "scheduler"
)

// Database represents a SQLite database.
// type Database struct {
// 	db *sql.DB
// }

// Создание новой базы данных или подключение к существующей БД
func NewSQLite3DB(dbName string) (*sql.DB, error) {

	if !checkDB(dbName) {
		if err := createDB(dbName); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	if err = createTable(db); err != nil {
		return nil, err
	}
	return db, nil // &Database{db: db}, nil
}

// Создание файла базы данных
func createDB(dbName string) error {
	log.Println("Create database file")
	appPath, err := os.Executable()
	if err != nil {
		return err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)

	_, err = os.Create(dbFile)
	if err != nil {
		return err
	}
	return nil
}

// Проверка существования файла базы данных
func checkDB(dbName string) bool {
	log.Println("Check exist database file")
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		log.Println("Database file not exist")
		return false
	}
	log.Println("Database file exist")
	return true
}

// // Close closes the database connection.
// func (db *Database) Close() error {
// 	return db.db.Close()
// }

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
