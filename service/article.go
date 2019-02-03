package service

import (
	"fasthttptest/model"
	"fasthttptest/repository"
	"time"
)

type Article struct {
	Repository *repository.Article
}

func (a *Article) Insert(article *model.Article) error {
	article.CreatedOn = time.Now()
	article.UpdatedOn = article.CreatedOn
	return a.Repository.Insert(article)
}