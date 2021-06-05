package origins

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"net/url"
	"regexp"
	"strings"
)

// AddPopupAction 添加源站
type AddPopupAction struct {
	actionutils.ParentAction
}

func (this *AddPopupAction) RunGet(params struct {
	ServerId       int64
	ReverseProxyId int64
	OriginType     string
}) {
	this.Data["reverseProxyId"] = params.ReverseProxyId
	this.Data["originType"] = params.OriginType

	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverType := serverTypeResp.Type
	this.Data["serverType"] = serverType

	// 是否为HTTP
	this.Data["isHTTP"] = serverType == "httpProxy" || serverType == "httpWeb"

	this.Show()
}

func (this *AddPopupAction) RunPost(params struct {
	OriginType string

	ReverseProxyId int64
	Weight         int32
	Protocol       string
	Addr           string
	Name           string

	ConnTimeout  int
	ReadTimeout  int
	MaxConns     int32
	MaxIdleConns int32
	IdleTimeout  int

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("addr", params.Addr).
		Require("请输入源站地址")

	addr := params.Addr

	// 是否是完整的地址
	if params.Protocol == "http" || params.Protocol == "https" {
		u, err := url.Parse(addr)
		if err == nil {
			addr = u.Host
		}
	}

	addr = regexp.MustCompile(`\s+`).ReplaceAllString(addr, "")
	portIndex := strings.LastIndex(addr, ":")
	if portIndex < 0 {
		if params.Protocol == "http" {
			addr += ":80"
		} else if params.Protocol == "https" {
			addr += ":443"
		} else {
			this.Fail("地址中需要带有端口")
		}
		portIndex = strings.LastIndex(addr, ":")
	}
	host := addr[:portIndex]
	port := addr[portIndex+1:]

	connTimeoutJSON, err := (&shared.TimeDuration{
		Count: int64(params.ConnTimeout),
		Unit:  shared.TimeDurationUnitSecond,
	}).AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	readTimeoutJSON, err := (&shared.TimeDuration{
		Count: int64(params.ReadTimeout),
		Unit:  shared.TimeDurationUnitSecond,
	}).AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	idleTimeoutJSON, err := (&shared.TimeDuration{
		Count: int64(params.IdleTimeout),
		Unit:  shared.TimeDurationUnitSecond,
	}).AsJSON()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	createResp, err := this.RPC().OriginRPC().CreateOrigin(this.AdminContext(), &pb.CreateOriginRequest{
		Name: params.Name,
		Addr: &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		},
		Description:     params.Description,
		Weight:          params.Weight,
		IsOn:            params.IsOn,
		ConnTimeoutJSON: connTimeoutJSON,
		ReadTimeoutJSON: readTimeoutJSON,
		IdleTimeoutJSON: idleTimeoutJSON,
		MaxConns:        params.MaxConns,
		MaxIdleConns:    params.MaxIdleConns,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	originId := createResp.OriginId
	originRef := &serverconfigs.OriginRef{
		IsOn:     true,
		OriginId: originId,
	}

	reverseProxyResp, err := this.RPC().ReverseProxyRPC().FindEnabledReverseProxy(this.AdminContext(), &pb.FindEnabledReverseProxyRequest{ReverseProxyId: params.ReverseProxyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxy := reverseProxyResp.ReverseProxy
	if reverseProxy == nil {
		this.ErrorPage(errors.New("reverse proxy should not be nil"))
		return
	}

	origins := []*serverconfigs.OriginRef{}
	switch params.OriginType {
	case "primary":
		if len(reverseProxy.PrimaryOriginsJSON) > 0 {
			err = json.Unmarshal(reverseProxy.PrimaryOriginsJSON, &origins)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	case "backup":
		if len(reverseProxy.BackupOriginsJSON) > 0 {
			err = json.Unmarshal(reverseProxy.BackupOriginsJSON, &origins)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
	}
	origins = append(origins, originRef)
	originsData, err := json.Marshal(origins)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	switch params.OriginType {
	case "primary":
		_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxyPrimaryOrigins(this.AdminContext(), &pb.UpdateReverseProxyPrimaryOriginsRequest{
			ReverseProxyId: params.ReverseProxyId,
			OriginsJSON:    originsData,
		})
	case "backup":
		_, err = this.RPC().ReverseProxyRPC().UpdateReverseProxyBackupOrigins(this.AdminContext(), &pb.UpdateReverseProxyBackupOriginsRequest{
			ReverseProxyId: params.ReverseProxyId,
			OriginsJSON:    originsData,
		})
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "为反向代理服务 %d 添加源站 %d", params.ReverseProxyId, originId)

	this.Success()
}
