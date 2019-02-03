package handler

import (
	"fasthttptest/log"
	"fasthttptest/repository"
	"fasthttptest/service"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
)

type MainHandler struct {
	user user
	article article
}

func NewMainHandler(db *pgx.ConnPool, redis *redis.Client) MainHandler {
	userRepository := repository.NewUserRepository(db, redis)
	articleRepository := repository.NewArticle(db)
	userService := service.User{Repository: &userRepository}
	articleService := service.Article{Repository: &articleRepository}
	return MainHandler {
		user: user{&userService},
		article: article{&articleService, &userService},
	}
}

func (mh *MainHandler) Handle(logger log.Logger, ctx *fasthttp.RequestCtx) (interface{}, error) {
	if isPrefix(ctx, "user") {
		return mh.user.handle(logger, ctx)
	} else if isPrefix(ctx, "article") {
		return mh.article.handle(logger, ctx)
	} else {
		return notFound(ctx)
	}
}