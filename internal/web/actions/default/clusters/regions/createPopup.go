package regions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入区域名称")

	createResp, err := this.RPC().NodeRegionRPC().CreateNodeRegion(this.AdminContext(), &pb.CreateNodeRegionRequest{Name: params.Name})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["region"] = maps.Map{
		"id":   createResp.NodeRegionId,
		"name": params.Name,
	}

	// 日志
	defer this.CreateLogInfo("创建节点区域 %d", createResp.NodeRegionId)

	this.Success()
}
