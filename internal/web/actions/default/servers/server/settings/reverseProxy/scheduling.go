package reverseProxy

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/schedulingconfigs"
)

type SchedulingAction struct {
	actionutils.ParentAction
}

func (this *SchedulingAction) Init() {
	this.FirstMenu("scheduling")
}

func (this *SchedulingAction) RunGet(params struct {
	ServerId int64
}) {
	reverseProxyResp, err := this.RPC().ServerRPC().FindAndInitServerReverseProxyConfig(this.AdminContext(), &pb.FindAndInitServerReverseProxyConfigRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxy := &serverconfigs.ReverseProxyConfig{}
	err = json.Unmarshal(reverseProxyResp.ReverseProxy, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["reverseProxyId"] = reverseProxy.Id

	schedulingCode := reverseProxy.FindSchedulingConfig().Code
	schedulingMap := schedulingconfigs.FindSchedulingType(schedulingCode)
	if schedulingMap == nil {
		this.ErrorPage(errors.New("invalid scheduling code '" + schedulingCode + "'"))
		return
	}
	this.Data["scheduling"] = schedulingMap

	this.Show()
}
