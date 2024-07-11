package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// DB содержит экземпляр базы данных и методы для работы с ней
type Repository struct {
	db *sqlx.DB
}

// NewDB создает новый экземпляр DB
func NewDB(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// DropTables удаляет необходимые таблицы в базе данных
func DropTables(db *sqlx.DB) {
	tables := []string{
		"reactions",
		"friends",
		"posts",
		"users",
		"session",
	}

	for _, table := range tables {
		query := "DROP TABLE IF EXISTS " + table + " CASCADE;"
		if _, err := db.Exec(query); err != nil {
			log.Fatalf("Error dropping table %s: %v", table, err)
		}
		log.Printf("Table %s dropped successfully", table)
	}
}

// CreateTables создает необходимые таблицы в базе данных
func CreateTables(db *sqlx.DB) {

	// Создание таблицы users
	q := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT,
			email TEXT,
			phone TEXT,
			password TEXT,
			createdAt TIMESTAMP
		);
	`

	if _, err := db.Exec(q); err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	// Создание таблицы friends
	q = `
		CREATE TABLE IF NOT EXISTS friends (
			userId INT,
			friendId INT,
			friendLogin TEXT,
			addedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (userId, friendId),
			FOREIGN KEY (userId) REFERENCES users(id),
			FOREIGN KEY (friendId) REFERENCES users(id)
		);
	`

	if _, err := db.Exec(q); err != nil {
		log.Fatalf("Error creating friends table: %v", err)
	}

	// Создание таблицы posts
	q = `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			content TEXT,
			author TEXT,
			tags TEXT[],
			createdAt TIMESTAMP,
			likeCount INT,
			dislikeCount INT
		);
	`

	if _, err := db.Exec(q); err != nil {
		log.Fatalf("Error creating posts table: %v", err)
	}

	// Создание таблицы reactions
	// type: 0 - dislike, 1 - like
	q = `
		CREATE TABLE IF NOT EXISTS reactions (
			id SERIAL PRIMARY KEY,
			userId INT,
			postId INT,
			reactionType INT CHECK (reactionType IN (0, 1)),
			createdAt TIMESTAMP,
			CONSTRAINT fk_user_post FOREIGN KEY (userId) REFERENCES users(id),
			CONSTRAINT fk_post_reaction FOREIGN KEY (postId) REFERENCES posts(id)
		);
	`

	if _, err := db.Exec(q); err != nil {
		log.Fatalf("Error creating reactions table: %v", err)
	}

	// Создание таблицы sessions
	q = `
		CREATE TABLE sessions (
			id SERIAL PRIMARY KEY,
			refresh_token VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`

}
