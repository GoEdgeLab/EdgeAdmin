package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	action.Data["teaMenu"] = "servers"

	selectedTabbar, _ := action.Data["mainTab"]

	tabbar := actionutils.NewTabbar()
	tabbar.Add("服务管理", "", "/servers", "", selectedTabbar == "server")
	actionutils.SetTabbar(action, tabbar)
}
