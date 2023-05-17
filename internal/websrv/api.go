package websrv

import (
	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"

	"github.com/valuedig/valueui/internal/status"
)

func NewApiModule() httpsrv.Module {

	module := httpsrv.NewModule("valueui_api")

	module.ControllerRegister(new(Layout))

	return module
}

type Layout struct {
	*httpsrv.Controller
}

func (c Layout) FetchAction() {
	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	item := status.Layouts.Get(c.Params.Get("name"))
	if item == nil {
		hlog.Printf("info", "layout fetch %s fail", c.Params.Get("name"))
		return
	}

	c.RenderJson(item)
}
