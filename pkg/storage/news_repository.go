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
	const query = `	INSERT INTO news (title, descriptions, link) VALUES ($1, $2, $3)`

	rows, err := nr.db.QueryContext(ctx, query, news.Title, news.Description, news.Link)
	if err != nil {
		return fmt.Errorf("create news: %w", err)
	}
	defer rows.Close()

	return nil
}

func (nr *newsRepository) GetNews(ctx context.Context, search string) ([]news.News, error) {
	const query = ` SELECT id, title, descriptions, link FROM news WHERE title ilike  '%' || $1 || '%'`

	rows, err := nr.db.QueryContext(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf("get news: %w", err)
	}
	defer rows.Close()

	ns := []news.News{}

	for rows.Next() {
		n := news.News{}

		if err := rows.Scan(&n.ID, &n.Title, &n.Description, &n.Link); err != nil {
			return nil, fmt.Errorf("scan events: %w", err)
		}

		ns = append(ns, n)
	}

	return ns, nil
}