package http

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	// 跳转相关设置
	webConfig, err := webutils.FindWebConfigWithLocationId(this.Parent(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id
	this.Data["redirectToHTTPSConfig"] = webConfig.RedirectToHttps

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId               int64
	RedirectToHTTPSJSON []byte

	Must *actions.Must
}) {
	// 设置跳转到HTTPS
	// TODO 校验设置
	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebRedirectToHTTPS(this.AdminContext(), &pb.UpdateHTTPWebRedirectToHTTPSRequest{
		WebId:               params.WebId,
		RedirectToHTTPSJSON: params.RedirectToHTTPSJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
