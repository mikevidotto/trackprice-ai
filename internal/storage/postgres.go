package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mikevidotto/blogfolio/internal/article"
)

type MypostgresStorage struct {
	db *sql.DB
}

type MypostgresConfig struct {
	Username string
	Password string
	DbName   string
	Port     uint
	Host     string
}

func NewMypostgresStorage(conf MypostgresConfig) (MypostgresStorage, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.Username, conf.Password, conf.Host, conf.Port, conf.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return MypostgresStorage{}, fmt.Errorf("cannot open postgres connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return MypostgresStorage{}, fmt.Errorf("cannot ping postgres connection: %w", err)
	}
	return MypostgresStorage{
		db: db,
	}, nil
}

func (s MypostgresStorage) Create(ctx context.Context, a article.Article) (article.Article, error) {
	query := `insert into articles (author, title, body, markdown, html, published)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	err := s.db.QueryRowContext(ctx, query, a.Author, a.Title, a.Body, a.Markdown, a.Html, a.Published).Scan(&a.Id)
	if err != nil {
		fmt.Errorf("cannot create article: %w", err)
	}
	return a, nil
}

func (s MypostgresStorage) Read(ctx context.Context) ([]article.Article, error) {
	var articles []article.Article
	var a article.Article
	query := "SELECT id, author, title, body, markdown, html, published FROM articles"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		fmt.Errorf("cannot retrieve articles from db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&a.Id, &a.Author, &a.Title, &a.Body, &a.Markdown, &a.Html, &a.Published)
		if err != nil {
			fmt.Errorf("cannot scan from rows: %w", err)
		}
		articles = append(articles, a)
		// fmt.Printf("%d - %s - %s - %s\n", a.Id, a.Author, a.Title, a.Body, a.Published)
	}
	err = rows.Err()
	if err != nil {
		fmt.Errorf("cannot check for errors in rows.... %w", err)
	}
	return articles, err
}

func (s MypostgresStorage) FindArticle(ctx context.Context, id int64) (article.Article, error) {
	var a article.Article
	query := "SELECT id, author, title, body, markdown, html, published FROM articles WHERE id = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&a.Id, &a.Author, &a.Title, &a.Body, &a.Markdown, &a.Html, &a.Published)
	if err != nil {
		fmt.Errorf("cannot scan row.. %w", err)
	}
	return a, nil
}

func (s MypostgresStorage) Update(ctx context.Context, a article.Article) (article.Article, error) {
	query := "UPDATE articles SET body = $1 WHERE id = $2"
	err := s.db.QueryRowContext(ctx, query, a.Body, a.Id)
	if err != nil {
		fmt.Errorf("cannot update article: %w", err)
	}
	return a, nil
}

func (s MypostgresStorage) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM articles WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, id)
	if err != nil {
		fmt.Errorf("cannot delete record from table: %w", err)
	}
	return nil
}
