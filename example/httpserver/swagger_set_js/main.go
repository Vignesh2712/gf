package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HelloReq struct {
	g.Meta `path:"/hello" method:"get" sort:"1"`
	Name   string `v:"required" dc:"Your name"`
}

type HelloRes struct {
	Reply string `dc:"Reply content"`
}

type Hello struct{}

func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	g.Log().Debugf(ctx, `receive say: %+v`, req)
	res = &HelloRes{
		Reply: fmt.Sprintf(`Hi %s`, req.Name),
	}
	return
}

func main() {
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(Hello),
		)
	})

	s.SetSwaggerJsURL("https://unpkg.com/redoc@2.0.0/bundles/redoc.standalone.js")
	// or with static files
	// s.AddStaticPath("/js", "js")
	// s.SetSwaggerJsURL("/js/redoc.standalone.js")

	s.Run()
}
