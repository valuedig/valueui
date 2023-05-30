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
	module.ControllerRegister(new(Datalet))

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

	case item.Template.Layout != nil:
		item.Template.Layout.Refix()
		str := fmt.Sprintf("<!-- tpl:layout:%s -->\n", item.Template.Name)
		str += fmt.Sprintf("<div class=\"valueui-container\">\n")
		str += fmt.Sprintf("<div class=\"valueui-row%s\"%s>\n",
			colUnitFilter("valueui-row-", item.Template.Layout.Width),
			colCssFilter(item.Template.Layout))
		for _, v := range item.Template.Layout.Cols {
			str += fmt.Sprintf("  <div id=\"valueui-%s\" class=\"valueui-col%s\"%s>%s</div>\n",
				v.Name, colUnitFilter("valueui-col-", v.Width), colCssFilter(v), v.Name)
		}
		str += "</div>\n"
		str += "</div>\n"
		item.Template.Html = &uiapi.TemplateHtml{
			Html: str,
		}

	case item.Template.Nav != nil:

		str := fmt.Sprintf("<!-- tpl:nav:%s -->\n", item.Template.Name)

		str += fmt.Sprintf("<nav id=\"nav-{[=it.name]}\" class=\"nav valueui-nav valueui-gap-box%s\"%s>\n",
			navClassFilter(item.Template.Nav),
			navCssFilter(item.Template.Nav),
		)

		str += "{[~it.rows :row]}\n" +
			"<li id=\"nav-item-{[=it.name]}-{[=row.name]}\" class=\"nav-item valueui-nav-item\">\n" +
			"  <a id=\"nav-link-{[=it.name]}-{[=row.name]}\" class=\"nav-link valueui-nav-link\">{[=row.title]}</a>\n" +
			"</li>\n" +
			"{[~]}\n"

		str += "</nav>\n"

		item.Template.Html = &uiapi.TemplateHtml{
			Html: str,
		}

	case item.Template.Html != nil && item.Template.Html.Html == "":
		if tpl := status.Assets.Get("template/" + item.Template.Html.File); tpl != nil {
			if h, ok := tpl.(*uiapi.TemplateHtml); ok {
				item.Template.Html.Html = h.Html
			}
		}
	}

	return nil
}

func navClassFilter(c *uiapi.TemplateNav) string {
	if c.Display == "flex-column" {
		return " " + c.Display
	}
	return ""
}

func navCssFilter(c *uiapi.TemplateNav) string {
	return ""
}

func colUnitFilter(prefix, v string) string {
	if v == "auto" {
		return " " + prefix + v
	}
	return ""
}

func colCssFilter(c *uiapi.TemplateLayout) string {

	unitFilter := func(s string) []string {
		ar := []string{}
		for _, v := range strings.Split(s, ",") {
			switch {
			case strings.HasSuffix(v, "rem"):
				if v, err := strconv.ParseFloat(v[:len(v)-3], 32); err == nil {
					ar = append(ar, fmt.Sprintf("%drem", int(v)))
				}
			case strings.HasSuffix(v, "px"):
				if v, err := strconv.ParseFloat(v[:len(v)-2], 32); err == nil {
					ar = append(ar, fmt.Sprintf("%dpx", int(v)))
				}
			case strings.HasSuffix(v, "vh"):
				if v, err := strconv.ParseFloat(v[:len(v)-2], 32); err == nil {
					ar = append(ar, fmt.Sprintf("%dvh", int(v)))
				}
			case strings.HasSuffix(v, "%"):
				if v, err := strconv.ParseFloat(v[:len(v)-1], 32); err == nil {
					ar = append(ar, fmt.Sprintf("%d%%", int(v)))
				}
			}
		}
		return ar
	}

	css := []string{}

	for _, n := range [][]string{
		{"width", c.Width},
		{"height", c.Height},
	} {
		ar := unitFilter(n[1])
		switch len(ar) {
		case 1:
			css = append(css,
				fmt.Sprintf("%s:%s", n[0], ar[0]))
		case 2:
			css = append(css,
				fmt.Sprintf("min-%s:%s", n[0], ar[0]),
				fmt.Sprintf("max-%s:%s", n[0], ar[1]))
		case 3:
			css = append(css,
				fmt.Sprintf("min-%s:%s", n[0], ar[0]),
				fmt.Sprintf("%s:%s", n[0], ar[1]),
				fmt.Sprintf("max-%s:%s", n[0], ar[2]))
		}
	}
	if len(css) == 0 {
		return ""
	}
	return " style=\"" + strings.Join(css, ";") + "\""
}
