package providers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	ProviderId int64
}) {
	// TODO 检查权限

	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "删除DNS服务商 %d", params.ProviderId)

	// 检查是否被使用
	countClustersResp, err := this.RPC().NodeClusterRPC().CountAllEnabledNodeClustersWithDNSProviderId(this.AdminContext(), &pb.CountAllEnabledNodeClustersWithDNSProviderIdRequest{DnsProviderId: params.ProviderId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countClustersResp.Count > 0 {
		this.Fail("当前DNS服务商账号正在被" + numberutils.FormatInt64(countClustersResp.Count) + "个集群所使用，所以不能删除。请修改后再操作。")
	}

	// 执行删除
	_, err = this.RPC().DNSProviderRPC().DeleteDNSProvider(this.AdminContext(), &pb.DeleteDNSProviderRequest{DnsProviderId: params.ProviderId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
