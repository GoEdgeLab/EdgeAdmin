package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/waf/ipadmin"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth(configloaders.AdminModuleCodeServer)).
			Helper(serverutils.NewServerHelper()).
			Prefix("/servers/server/settings/waf").
			GetPost("", new(IndexAction)).
			Get("/ipadmin/allowList", new(ipadmin.AllowListAction)).
			Get("/ipadmin/denyList", new(ipadmin.DenyListAction)).
			GetPost("/ipadmin/countries", new(ipadmin.CountriesAction)).
			GetPost("/ipadmin/provinces", new(ipadmin.ProvincesAction)).
			GetPost("/ipadmin/createIPPopup", new(ipadmin.CreateIPPopupAction)).
			GetPost("/ipadmin/updateIPPopup", new(ipadmin.UpdateIPPopupAction)).
			Post("/ipadmin/deleteIP", new(ipadmin.DeleteIPAction)).

			// 规则相关
			Get("/groups", new(GroupsAction)).
			Get("/group", new(GroupAction)).
			EndAll()
	})
}
