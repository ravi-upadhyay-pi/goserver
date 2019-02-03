package handler

import (
	"encoding/json"
	"fasthttptest/log"
	"fasthttptest/model"
	"fasthttptest/service"
	"github.com/valyala/fasthttp"
)

type user struct {
	userService *service.UserService
}

func (h *user) handle(logger log.Logger, ctx *fasthttp.RequestCtx) (interface{}, error) {
	if match(ctx, "POST", ":username", ":password") {
		username := ctx.UserValue(":username").(string)
		password := ctx.UserValue(":password").(string)
		return h.userService.CreateSession(logger, username, password)
	} else if match(ctx, "GET") {
		sessionId := string(ctx.Request.Header.Peek("Session-Id"))
		return h.userService.GetProfile(logger, sessionId)
	} else if match(ctx, "DELETE") {
		sessionId := string(ctx.Request.Header.Peek("Session-Id"))
		return nil, h.userService.RemoveSession(logger, sessionId)
	} else if match(ctx, "DELETE", "all"){
		sessionId := string(ctx.Request.Header.Peek("Session-Id"))
		return nil, h.userService.RemoveAllSession(logger, sessionId)
	} else if match(ctx, "POST") {
		user := &model.User{}
		err := json.Unmarshal(ctx.PostBody(), user)
		if err != nil {
			return nil, err
		}
		return nil, h.userService.Add(logger, user)
	} else {
		return notFound(ctx)
	}
}
