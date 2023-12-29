package article

import (
	"context"
	v1 "github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genpb/api/article/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedArticleRpcServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterArticleRpcServer(s.Server, &Controller{})
}

func (*Controller) GetArticle(ctx context.Context, req *v1.GetArticleReq) (res *v1.GetArticleRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CreateArticle(ctx context.Context, req *v1.CreateArticleReq) (res *v1.CreateArticleRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
