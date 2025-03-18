package storage

import (
	"context"
	"tg-bot-rss/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {

}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit int64) ([]model.Article, error) {

}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {

}

type dbArticle struct {
	ID          int64     `db:"id"`
	SourceID    int64     `db:"source_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
}
