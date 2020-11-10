package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(clusterutils.NewClustersHelper()).
			Data("teaSubMenu", "grant").
			Prefix("/clusters/grants").

			// 授权管理
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/delete", new(DeleteAction)).
			Get("/grant", new(GrantAction)).
			GetPost("/selectPopup", new(SelectPopupAction)).
			GetPost("/createPopup", new(CreatePopupAction)).
			GetPost("/updatePopup", new(UpdatePopupAction)).
			EndAll()
	})
}
