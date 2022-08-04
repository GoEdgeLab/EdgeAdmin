package certs

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	CertId int64
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除SSL证书 %d", params.CertId)

	// 是否正在被服务使用
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledServersWithSSLCertIdRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此证书正在被某些服务引用，请先修改服务后再删除。")
	}

	// 是否正在被API节点使用
	countResp, err = this.RPC().APINodeRPC().CountAllEnabledAPINodesWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledAPINodesWithSSLCertIdRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("此证书正在被某些API节点引用，请先修改API节点后再删除")
	}

	// 是否正在被用户节点使用
	if teaconst.IsPlus {
		countResp, err = this.RPC().UserNodeRPC().CountAllEnabledUserNodesWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledUserNodesWithSSLCertIdRequest{SslCertId: params.CertId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if countResp.Count > 0 {
			this.Fail("此证书正在被某些用户节点引用，请先修改相关用户节点后再删除")
		}
	}

	// 是否正在被NS集群使用
	if teaconst.IsPlus {
		countResp, err = this.RPC().NSClusterRPC().CountAllNSClustersWithSSLCertId(this.AdminContext(), &pb.CountAllNSClustersWithSSLCertIdRequest{SslCertId: params.CertId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if countResp.Count > 0 {
			this.Fail("此证书正在被某些DNS集群节点引用，请先修改相关DNS集群设置后再删除")
		}
	}

	_, err = this.RPC().SSLCertRPC().DeleteSSLCert(this.AdminContext(), &pb.DeleteSSLCertRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
