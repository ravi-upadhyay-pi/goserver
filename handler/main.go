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
}

func NewMainHandler(db *pgx.ConnPool, redis *redis.Client) MainHandler {
	userRepository := repository.NewUserRepository(db, redis)
	userService := service.NewUserService(&userRepository)
	return MainHandler {
		user: user{&userService},
	}
}

func (mh *MainHandler) Handle(logger log.Logger, ctx *fasthttp.RequestCtx) (interface{}, error) {
	if isPrefix(ctx, "user") {
		return mh.user.handle(logger, ctx)
	} else {
		return notFound(ctx)
	}
}