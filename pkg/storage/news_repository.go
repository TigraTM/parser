package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"


	"parser/pkg/news"
)

type newsRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) news.Repository {
	return &newsRepository{
		db: db,
	}
}

func (nr *newsRepository) CreateNews(ctx context.Context, news news.News) error {
	const query = `	INSERT INTO news (title, descriptions, link) VALUES($1, $2, $3)`

	rows, err := nr.db.QueryContext(ctx, query, news.Title, news.Description, news.Link)
	if err != nil {
		return fmt.Errorf("create news: %w", err)
	}
	defer rows.Close()

	return nil
}

func (nr *newsRepository) GetNews(ctx context.Context, search string) ([]news.News, error) {
	return nil, nil
}