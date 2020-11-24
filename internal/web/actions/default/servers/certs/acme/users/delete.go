package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	UserId int64
}) {
	defer this.CreateLogInfo("删除ACME用户 %d", params.UserId)

	countResp, err := this.RPC().SSLCertRPC().CountSSLCertsWithACMEUserId(this.AdminContext(), &pb.CountSSLCertsWithACMEUserIdRequest{AcmeUserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if countResp.Count > 0 {
		this.Fail("有证书正在和这个用户关联，所以不能删除")
	}

	_, err = this.RPC().ACMEUserRPC().DeleteACMEUser(this.AdminContext(), &pb.DeleteACMEUserRequest{AcmeUserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
