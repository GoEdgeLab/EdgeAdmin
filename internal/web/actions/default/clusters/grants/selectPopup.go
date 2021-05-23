package grants

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct{}) {
	// 所有的认证
	grantsResp, err := this.RPC().NodeGrantRPC().FindAllEnabledNodeGrants(this.AdminContext(), &pb.FindAllEnabledNodeGrantsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	grants := grantsResp.NodeGrants
	grantMaps := []maps.Map{}
	for _, grant := range grants {
		grantMaps = append(grantMaps, maps.Map{
			"id":          grant.Id,
			"name":        grant.Name,
			"method":      grant.Method,
			"methodName":  grantutils.FindGrantMethodName(grant.Method),
			"username":    grant.Username,
			"description": grant.Description,
		})
	}
	this.Data["grants"] = grantMaps

	this.Show()
}

func (this *SelectPopupAction) RunPost(params struct {
	GrantId int64
	Must    *actions.Must
}) {
	if params.GrantId <= 0 {
		this.Data["grant"] = maps.Map{
			"id":         params.GrantId,
			"name":       "",
			"method":     "",
			"methodName": "",
		}
		this.Success()
	}

	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledNodeGrant(this.AdminContext(), &pb.FindEnabledNodeGrantRequest{NodeGrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	grant := grantResp.NodeGrant
	if grant == nil {
		this.Fail("找不到要使用的认证")
	}
	this.Data["grant"] = maps.Map{
		"id":         grant.Id,
		"name":       grant.Name,
		"method":     grant.Method,
		"methodName": grantutils.FindGrantMethodName(grant.Method),
	}

	this.Success()
}
