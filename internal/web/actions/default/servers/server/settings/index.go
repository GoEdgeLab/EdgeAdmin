package settings

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// IndexAction 服务基本信息设置
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 所有集群
	resp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range resp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 当前服务信息
	serverResp, err := this.RPC().ServerRPC().FindEnabledServer(this.AdminContext(), &pb.FindEnabledServerRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	server := serverResp.Server
	if server == nil {
		this.NotFound("server", params.ServerId)
		return
	}

	// 用户
	if server.User != nil {
		this.Data["user"] = maps.Map{
			"id":       server.User.Id,
			"fullname": server.User.Fullname,
			"username": server.User.Username,
		}
	} else {
		this.Data["user"] = nil
	}

	// 集群
	clusterId := int64(0)
	this.Data["clusterName"] = ""
	if server.NodeCluster != nil {
		clusterId = server.NodeCluster.Id
		this.Data["clusterName"] = server.NodeCluster.Name
	}

	// 分组
	groupMaps := []maps.Map{}
	if len(server.Groups) > 0 {
		for _, group := range server.Groups {
			groupMaps = append(groupMaps, maps.Map{
				"id":   group.Id,
				"name": group.Name,
			})
		}
	}

	this.Data["server"] = maps.Map{
		"id":          server.Id,
		"clusterId":   clusterId,
		"type":        server.Type,
		"name":        server.Name,
		"description": server.Description,
		"isOn":        server.IsOn,
		"groups":      groupMaps,
	}

	serverType := serverconfigs.FindServerType(server.Type)
	if serverType == nil {
		this.ErrorPage(errors.New("invalid server type '" + server.Type + "'"))
		return
	}

	typeName := serverType.GetString("name")
	this.Data["typeName"] = typeName

	// 记录最近使用
	_, err = this.RPC().LatestItemRPC().IncreaseLatestItem(this.AdminContext(), &pb.IncreaseLatestItemRequest{
		ItemType: "server",
		ItemId:   params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

// RunPost 保存
func (this *IndexAction) RunPost(params struct {
	ServerId    int64
	Name        string
	Description string
	ClusterId   int64
	GroupIds    []int64
	IsOn        bool

	Must *actions.Must
}) {
	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "修改代理服务 %d 基本信息", params.ServerId)

	params.Must.
		Field("name", params.Name).
		Require("请输入服务名称")

	if params.ClusterId <= 0 {
		this.Fail("请选择部署的集群")
	}

	_, err := this.RPC().ServerRPC().UpdateServerBasic(this.AdminContext(), &pb.UpdateServerBasicRequest{
		ServerId:      params.ServerId,
		Name:          params.Name,
		Description:   params.Description,
		NodeClusterId: params.ClusterId,
		IsOn:          params.IsOn,
		GroupIds:      params.GroupIds,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
