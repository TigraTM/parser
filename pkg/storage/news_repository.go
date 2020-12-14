package storage

import (
	"github.com/jmoiron/sqlx"

	"parser/pkg/news"
)

type NewsRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) news.Repository {
	return NewsRepository{
		db: db,
	}
}
