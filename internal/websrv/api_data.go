package websrv

import (
	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"

	// "google.golang.org/protobuf/types/known/structpb"

	"github.com/valuedig/valueui/internal/status"
	"github.com/valuedig/valueui/pkg/uiapi"
)

type Datalet struct {
	*httpsrv.Controller
}

func (c Datalet) RunAction() {
	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	var resp uiapi.DataletRunResponse
	defer c.RenderJson(&resp)

	item := status.Assets.Get(c.Params.Get("name"))
	if item == nil {
		hlog.Printf("info", "viewlet fetch %s fail", c.Params.Get("name"))
		return
	}

	datalet, ok := item.(*uiapi.DataletSpec)
	if !ok {
		hlog.Printf("info", "viewlet (%s) fetch fail : object type error", c.Params.Get("name"))
		return
	}

	resp.Kind = "DataletRunResponse"

	if datalet.Action != nil {
		switch {
		case datalet.Action.Search != nil:
			item = status.Assets.Get("model/" + datalet.Action.Search.Model + ".toml")
			if item == nil {
				return
			}
			model, ok := item.(*uiapi.ModelSpec)
			if !ok {
				return
			}
			if len(model.DefaultRows) > 0 {
				resp.Datasets = append(resp.Datasets, &uiapi.DataResult{
					Name: datalet.Name,
					Spec: &uiapi.ModelSpec{
						Fields: model.Fields,
					},
					Rows: model.DefaultRows,
				})
			}
		}
	}
}
