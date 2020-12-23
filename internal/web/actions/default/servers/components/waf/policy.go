package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type PolicyAction struct {
	actionutils.ParentAction
}

func (this *PolicyAction) Init() {
	this.Nav("", "", "index")
}

func (this *PolicyAction) RunGet(params struct {
	FirewallPolicyId int64
}) {
	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyConfig(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy == nil {
		this.NotFound("firewallPolicy", params.FirewallPolicyId)
		return
	}

	internalGroups := []maps.Map{}
	if firewallPolicy.Inbound != nil {
		for _, group := range firewallPolicy.Inbound.Groups {
			internalGroups = append(internalGroups, maps.Map{
				"name": group.Name,
				"isOn": group.IsOn,
			})
		}
	}
	if firewallPolicy.Outbound != nil {
		for _, group := range firewallPolicy.Outbound.Groups {
			internalGroups = append(internalGroups, maps.Map{
				"name": group.Name,
				"isOn": group.IsOn,
			})
		}
	}

	this.Data["firewallPolicy"] = maps.Map{
		"id":           firewallPolicy.Id,
		"name":         firewallPolicy.Name,
		"isOn":         firewallPolicy.IsOn,
		"description":  firewallPolicy.Description,
		"groups":       internalGroups,
		"blockOptions": firewallPolicy.BlockOptions,
	}

	// 正在使用此策略的集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithHTTPFirewallPolicyId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithHTTPFirewallPolicyIdRequest{HttpFirewallPolicyId: params.FirewallPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
