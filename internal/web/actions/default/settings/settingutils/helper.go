//go:build !plus

package settingutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
	tab string
}

func NewHelper(tab string) *Helper {
	return &Helper{
		tab: tab,
	}
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	goNext = true

	var action = actionPtr.Object()

	// 左侧菜单
	action.Data["teaMenu"] = "settings"

	// 标签栏
	var tabbar = actionutils.NewTabbar()
	var session = action.Session()
	var adminId = session.GetInt64("adminId")
	if configloaders.AllowModule(adminId, configloaders.AdminModuleCodeSetting) {
		tabbar.Add("Web服务", "", "/settings/server", "", this.tab == "server")
		tabbar.Add("管理界面设置", "", "/settings/ui", "", this.tab == "ui")
		tabbar.Add("安全设置", "", "/settings/security", "", this.tab == "security")
		tabbar.Add("检查更新", "", "/settings/updates", "", this.tab == "updates")
	}
	tabbar.Add("个人资料", "", "/settings/profile", "", this.tab == "profile")
	tabbar.Add("登录设置", "", "/settings/login", "", this.tab == "login")
	actionutils.SetTabbar(actionPtr, tabbar)

	return
}
