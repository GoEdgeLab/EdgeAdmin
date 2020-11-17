package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
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
	// 预置分组
	groups := []maps.Map{}
	templatePolicy := firewallconfigs.HTTPFirewallTemplate()
	for _, group := range templatePolicy.AllRuleGroups() {
		groups = append(groups, maps.Map{
			"code": group.Code,
			"name": group.Name,
			"isOn": group.IsOn,
		})
	}
	this.Data["groups"] = groups

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name        string
	GroupCodes  []string
	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	createResp, err := this.RPC().HTTPFirewallPolicyRPC().CreateHTTPFirewallPolicy(this.AdminContext(), &pb.CreateHTTPFirewallPolicyRequest{
		IsOn:               params.IsOn,
		Name:               params.Name,
		Description:        params.Description,
		FirewallGroupCodes: params.GroupCodes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 日志
	this.CreateLog(oplogs.LevelInfo, "创建WAF策略 %d", createResp.FirewallPolicyId)

	this.Success()
}
