package storage

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

type UserExports struct {
	ID        string    `json:"id"`
	Cookies   string    `json:"cookies"`
	CreatedAt time.Time `json:"created_at"`
}

type SQLiteDB struct {
	DB *sql.DB
}

func (s *SQLiteDB) Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/exports.db"
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
       CREATE TABLE IF NOT EXISTS exports (
           id TEXT PRIMARY KEY,
           data TEXT NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
       )
   `)
	s.DB = db
}

func (s *SQLiteDB) SaveExport(data string) (string, error) {
	id := uuid.New().String()
	_, err := s.DB.Exec(`INSERT INTO exports (id, data, created_at) VALUES (?, ?, ?)`,
		id, data, time.Now())
	if err != nil {
		log.Printf("Error inserting data into DB: %v", err)
	}
	return id, err
}

func (s *SQLiteDB) GetExport(id string) (string, error) {
	var data string
	err := s.DB.QueryRow(`SELECT data FROM  exports WHERE id = ?`, id).Scan(&data)

	return data, err
}

func (s *SQLiteDB) StartCleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			_, err := s.DB.Exec(`DELETE FROM exports WHERE created_at < ?`,
				time.Now().Add(-24*time.Hour))
			if err != nil {
				log.Printf("Cleanup error: %v", err)
			}
		}
	}()
}
