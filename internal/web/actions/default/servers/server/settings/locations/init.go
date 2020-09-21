package locations

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(serverutils.NewServerHelper()).
			Prefix("/servers/server/settings/locations").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			Post("/delete", new(DeleteAction)).
			EndAll()
	})
}
