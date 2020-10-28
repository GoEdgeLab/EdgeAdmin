package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	GroupId int64
}) {
	groupResp, err := this.RPC().NodeGroupRPC().FindEnabledNodeGroup(this.AdminContext(), &pb.FindEnabledNodeGroupRequest{GroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	group := groupResp.Group
	if group == nil {
		this.NotFound("nodeGroup", params.GroupId)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	GroupId int64
	Name    string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")
	_, err := this.RPC().NodeGroupRPC().UpdateNodeGroup(this.AdminContext(), &pb.UpdateNodeGroupRequest{
		GroupId: params.GroupId,
		Name:    params.Name,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
