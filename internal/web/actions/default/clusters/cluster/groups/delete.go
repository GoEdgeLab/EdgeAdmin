package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GroupId int64
}) {
	// 检查是否正在使用
	countResp, err := this.RPC().NodeRPC().CountAllEnabledNodesWithGroupId(this.AdminContext(), &pb.CountAllEnabledNodesWithGroupIdRequest{GroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if countResp.Count > 0 {
		this.Fail("此分组正在被使用不能删除，请修改节点后再删除")
	}

	_, err = this.RPC().NodeGroupRPC().DeleteNodeGroup(this.AdminContext(), &pb.DeleteNodeGroupRequest{GroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除集群分组 %d", params.GroupId)

	this.Success()
}
