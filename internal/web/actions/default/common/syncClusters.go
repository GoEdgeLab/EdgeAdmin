package common

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type SyncClustersAction struct {
	actionutils.ParentAction
}

func (this *SyncClustersAction) RunPost(params struct{}) {
	// TODO 将来可以单独选择某一个集群进行单独的同步

	// 所有有变化的集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllChangedClusters(this.AdminContext(), &pb.FindAllChangedClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusters := clustersResp.Clusters

	for _, cluster := range clusters {
		_, err := this.RPC().NodeRPC().SyncNodesVersionWithCluster(this.AdminContext(), &pb.SyncNodesVersionWithClusterRequest{
			ClusterId: cluster.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Success()
}
