package websrv

import (
	"github.com/hooto/httpsrv"
	// "github.com/valuedig/valueui/pkg/uiapi"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("valueui_index")

	module.RouteSet(httpsrv.Route{
		Type:       httpsrv.RouteTypeStatic,
		Path:       "~",
		StaticPath: "./",
	})

	module.ControllerRegister(new(Index))

	return module
}

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	c.RenderString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>dce</title>
  <script src="/demo/~/main.js"></script>
  <script type="text/javascript">
    valueui.basepath = "/demo";
    valueui.uipath = "~";
    window.onload = valueui.main();
  </script>
</head>
<body id="body-content">loading</body>
</html>
`)
}
