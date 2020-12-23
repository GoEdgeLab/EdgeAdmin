package http

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("http")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	httpConfig := &serverconfigs.HTTPProtocolConfig{}
	if len(server.HttpJSON) > 0 {
		err := json.Unmarshal(server.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpConfig.IsOn = true
	}

	this.Data["serverType"] = server.Type
	this.Data["httpConfig"] = maps.Map{
		"isOn":      httpConfig.IsOn,
		"addresses": httpConfig.Listen,
	}

	// 跳转相关设置
	webConfig, err := webutils.FindWebConfigWithServerId(this.Parent(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id
	this.Data["redirectToHTTPSConfig"] = webConfig.RedirectToHttps

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	IsOn      bool
	Addresses string

	WebId               int64
	RedirectToHTTPSJSON []byte

	Must *actions.Must
}) {
	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "修改服务 %d 的HTTP设置", params.ServerId)

	addresses := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	httpConfig := &serverconfigs.HTTPProtocolConfig{}
	if len(server.HttpJSON) > 0 {
		err = json.Unmarshal(server.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	httpConfig.IsOn = params.IsOn
	httpConfig.Listen = addresses
	configData, err := json.Marshal(httpConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerHTTP(this.AdminContext(), &pb.UpdateServerHTTPRequest{
		ServerId: params.ServerId,
		HttpJSON: configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 设置跳转到HTTPS
	// TODO 校验设置
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRedirectToHTTPS(this.AdminContext(), &pb.UpdateHTTPWebRedirectToHTTPSRequest{
		WebId:               params.WebId,
		RedirectToHTTPSJSON: params.RedirectToHTTPSJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
