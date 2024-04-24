package httpclient

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type IHandler interface {
	HandleResponse(ctx context.Context, res *gclient.Response, out interface{}) error
}

// DefaultHandler handle ghttp.DefaultHandlerResponse of json format.
type DefaultHandler struct {
	Logger  *glog.Logger
	RawDump bool
}

func NewDefaultHandler(config Config) *DefaultHandler {
	return &DefaultHandler{
		Logger:  config.Logger,
		RawDump: config.RawDump,
	}
}

func (h DefaultHandler) HandleResponse(ctx context.Context, res *gclient.Response, out interface{}) error {
	defer res.Close()
	if h.RawDump {
		h.Logger.Debugf(ctx, "raw request&response:\n%s", res.Raw())
	}
	var (
		responseBytes = res.ReadAll()
		result        = ghttp.DefaultHandlerResponse{
			Data: out,
		}
	)
	if !json.Valid(responseBytes) {
		return gerror.Newf(`invalid response content: %s`, responseBytes)
	}
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return gerror.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	if result.Code != gcode.CodeOK.Code() {
		return gerror.NewCode(
			gcode.New(result.Code, result.Message, nil),
			result.Message,
		)
	}
	return nil
}
