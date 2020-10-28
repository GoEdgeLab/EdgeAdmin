package cluster

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type InstallNodesAction struct {
	actionutils.ParentAction
}

func (this *InstallNodesAction) Init() {
	this.Nav("", "node", "install")
	this.SecondMenu("nodes")
}

func (this *InstallNodesAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["leftMenuItems"] = LeftMenuItemsForInstall(params.ClusterId, "register")

	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{ClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if clusterResp.Cluster == nil {
		this.NotFound("nodeCluster", params.ClusterId)
		return
	}

	cluster := clusterResp.Cluster

	clusterAPINodesResp, err := this.RPC().NodeClusterRPC().FindAPINodesWithNodeCluster(this.AdminContext(), &pb.FindAPINodesWithNodeClusterRequest{ClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	apiNodeAddrs := []string{}
	if clusterAPINodesResp.UseAllAPINodes {
		apiNodesResp, err := this.RPC().APINodeRPC().FindAllEnabledAPINodes(this.AdminContext(), &pb.FindAllEnabledAPINodesRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, apiNode := range apiNodesResp.Nodes {
			if !apiNode.IsOn {
				continue
			}
			apiNodeAddrs = append(apiNodeAddrs, apiNode.AccessAddrs...)
		}
	} else {
		for _, apiNode := range clusterAPINodesResp.ApiNodes {
			if !apiNode.IsOn {
				continue
			}
			apiNodeAddrs = append(apiNodeAddrs, apiNode.AccessAddrs...)
		}
	}

	this.Data["cluster"] = maps.Map{
		"uniqueId":  cluster.UniqueId,
		"secret":    cluster.Secret,
		"endpoints": "\"" + strings.Join(apiNodeAddrs, "\", \"") + "\"",
	}

	this.Show()
}
