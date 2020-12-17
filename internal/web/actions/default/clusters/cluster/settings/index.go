package settings

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/grants/grantutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cluster := clusterResp.Cluster
	if cluster == nil {
		this.WriteString("not found cluster")
		return
	}

	// 认证
	var grantMap interface{} = nil

	if cluster.GrantId > 0 {
		grantResp, err := this.RPC().NodeGrantRPC().FindEnabledGrant(this.AdminContext(), &pb.FindEnabledGrantRequest{GrantId: cluster.GrantId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		grant := grantResp.Grant
		if grant != nil {
			grantMap = maps.Map{
				"id":         grant.Id,
				"name":       grant.Name,
				"method":     grant.Method,
				"methodName": grantutils.FindGrantMethodName(grant.Method),
			}
		}
	}
	this.Data["grant"] = grantMap

	this.Data["cluster"] = maps.Map{
		"id":         cluster.Id,
		"name":       cluster.Name,
		"installDir": cluster.InstallDir,
	}

	this.Show()
}

// 保存设置
func (this *IndexAction) RunPost(params struct {
	ClusterId  int64
	Name       string
	GrantId    int64
	InstallDir string

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改集群基础设置 %d", params.ClusterId)

	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	_, err := this.RPC().NodeClusterRPC().UpdateNodeCluster(this.AdminContext(), &pb.UpdateNodeClusterRequest{
		NodeClusterId: params.ClusterId,
		Name:          params.Name,
		GrantId:       params.GrantId,
		InstallDir:    params.InstallDir,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
