package websrv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"

	"github.com/valuedig/valueui/internal/status"
	"github.com/valuedig/valueui/pkg/uiapi"
)

func NewApiModule() httpsrv.Module {

	module := httpsrv.NewModule("valueui_api")

	module.ControllerRegister(new(Viewlet))

	return module
}

type Viewlet struct {
	*httpsrv.Controller
}

func (c Viewlet) FetchAction() {
	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	item := status.Assets.Get(c.Params.Get("name"))
	if item == nil {
		hlog.Printf("info", "viewlet fetch %s fail", c.Params.Get("name"))
		return
	}

	viewlet, ok := item.(*uiapi.Viewlet)
	if !ok {
		hlog.Printf("info", "viewlet (%s) fetch fail : object type error", c.Params.Get("name"))
		return
	}

	if err := viewletPreRender(viewlet); err != nil {
		hlog.Printf("info", "viewlet (%s) pre-render err %s", c.Params.Get("name"), err.Error())
		return
	}

	c.RenderJson(viewlet)
}

func viewletPreRender(item *uiapi.Viewlet) error {
	if item.Template == nil {
		return nil
	}

	switch {

	case item.Template.Html != nil && item.Template.Html.Html == "":
		if tpl := status.Assets.Get("template/" + item.Template.Html.File); tpl != nil {
			if h, ok := tpl.(*uiapi.TemplateHtml); ok {
				item.Template.Html.Html = h.Html
			}
		}

	case item.Template.Layout != nil:
		item.Template.Layout.Refix()
		str := fmt.Sprintf("<!-- id:%s -->\n", item.Template.Name)
		str += fmt.Sprintf("<div class=\"valueui-container\">\n")
		str += fmt.Sprintf("<div class=\"valueui-row valueui-row-%s\">\n", item.Template.Layout.Width)
		for _, v := range item.Template.Layout.Cols {
			str += fmt.Sprintf("  <div class=\"valueui-col valueui-col-%s\" style=\"width:%s\">%s</div>\n",
				v.Width, cssUnitFilter(v.Width), v.Name)
		}
		str += "</div>\n"
		str += "</div>\n"
		item.Template.Html = &uiapi.TemplateHtml{
			Html: str,
		}
	}

	return nil
}

func cssUnitFilter(v string) string {
	switch {
	case strings.HasSuffix(v, "rem"):
		if v, err := strconv.ParseFloat(v[:len(v)-3], 32); err == nil {
			return fmt.Sprintf("%drem", int(v))
		}
	case strings.HasSuffix(v, "px"):
		if v, err := strconv.ParseFloat(v[:len(v)-2], 32); err == nil {
			return fmt.Sprintf("%dpx", int(v))
		}
	case strings.HasSuffix(v, "%"):
		if v, err := strconv.ParseFloat(v[:len(v)-2], 32); err == nil {
			return fmt.Sprintf("%.2f%", v)
		}
	}
	return "auto"
}
