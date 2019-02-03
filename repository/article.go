package repository

import (
	"fasthttptest/ferror"
	"fasthttptest/model"
	"github.com/jackc/pgx"
)

type Article struct {
	db *pgx.ConnPool
}

func NewArticle(db *pgx.ConnPool) Article {
	return Article{db: db}
}

func (a *Article) Insert(article *model.Article) error {
	sql := `
		insert into article
		(id, username, title, body, tags, format, next, previous, private, created_on, updated_on) values 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := a.db.Exec(sql, article.Id, article.Username, article.Title, article.Body, article.Tags, article.Format,
		article.Next, article.Previous, article.Private, article.CreatedOn, article.UpdatedOn)
	if pgerr, ok := err.(pgx.PgError); ok {
		if pgerr.ConstraintName == "article_pkey" {
			return ferror.New(400, ferror.ArticleIdNotAvailable, "Article id is already taken. Change it.")
		}
	}
	return err
}

func (a *Article) GetByUsername(username string, offset uint64, len uint64) ([]model.Article, error) {
	sql := `select id, title, tags, created_on, updated_on from article where username = $1 offset $2 limit $3`
	rows, err := a.db.Query(sql, username, offset, len)
	if err != nil {return nil, err}
	var articles []model.Article
	for rows.Next() {
		article := model.Article{Username: username}
		rows.Scan(&article.Id, &article.Title, &article.Tags, &article.CreatedOn, &article.UpdatedOn)
		articles = append(articles, article)
	}
	return articles, rows.Err()
}

func (a *Article) GetByUsernameAndId(username string, id string) (model.Article, error) {
	sql := `
		select id, username, title, body, tags, format, next, previous, private, created_on, updated_on
		from article where username = $1 and id = $2`
	var article model.Article
	err := a.db.QueryRow(sql, username, id).Scan(&article.Id, &article.Username, &article.Title, &article.Body,
		&article.Tags, &article.Format, &article.Next, &article.Previous, &article.Private, &article.CreatedOn,
		&article.UpdatedOn)
	return article, err
}