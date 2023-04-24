package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"shop/backend-api/pkg/config"
)

func NewPostgresDB(cfg *config.DB) (*sql.DB, error) {
	postgresDbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBhost, cfg.DBport, cfg.DBuser, cfg.DBpass, cfg.DBname)

	db, err := sql.Open("postgres", postgresDbInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}
	err = createAdmin(db)
	if err != nil {
		log.Println(err)
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	path := "./internal/repository/tables/"
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, v := range dir {
		query, err := os.ReadFile(path + v.Name())
		if err != nil {
			return err
		}
		_, err = db.Exec(string(query))
		if err != nil {
			return err
		}

	}
	return nil
}

func createAdmin(db *sql.DB) error {
	query := `INSERT INTO users(nickname, age, gender, first_name, last_name, email, password, role, created, updated) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := db.Exec(query, "admin", 18, "Male", "master", "master", "admin@mail.ru", "pass", "admin", time.Now(), time.Now())
	return err
}
