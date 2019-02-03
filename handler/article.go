package handler

import (
	"encoding/json"
	"fasthttptest/log"
	"fasthttptest/model"
	"fasthttptest/service"
	"github.com/valyala/fasthttp"
	"strconv"
)

type article struct {
	service *service.Article
	userService *service.User
}

func (a *article) handle(logger log.Logger, ctx *fasthttp.RequestCtx) (interface{}, error) {
	if match(ctx, "POST") {
		sessionId := string(ctx.Request.Header.Peek("Session-Id"))
		username, err := a.userService.Repository.GetSession(sessionId)
		if err != nil {return nil, err}
		article := model.Article{}
		err = json.Unmarshal(ctx.PostBody(), &article)
		if err != nil {return nil, err}
		article.Username = username
		return nil, a.service.Insert(&article)
	} else if match(ctx, "GET") {
		sessionId := string(ctx.Request.Header.Peek("Session-Id"))
		username, err := a.userService.Repository.GetSession(sessionId)
		if err != nil {return nil, err}
		offset, err := strconv.ParseUint(string(ctx.QueryArgs().Peek("offset")), 10, 64)
		if err != nil {return nil, err}
		len, err := strconv.ParseUint(string(ctx.QueryArgs().Peek("len")), 10, 64)
		if err != nil {return nil, err}
		return a.service.Repository.GetByUsername(username, offset, len)
	} else if match(ctx, "GET", ":id") {
		sessionId := string(ctx.Request.Header.Peek("Session-Id"))
		username, err := a.userService.Repository.GetSession(sessionId)
		if err != nil {return nil, err}
		id := ctx.UserValue(":id").(string)
		return a.service.Repository.GetByUsernameAndId(username, id)
	} else {
		return notFound(ctx)
	}
}
