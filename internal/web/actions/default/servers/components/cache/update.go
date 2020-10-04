package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	CachePolicyId int64
}) {
	configResp, err := this.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPCachePolicyConfigRequest{CachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	configJSON := configResp.CachePolicyJSON
	if len(configJSON) == 0 {
		this.NotFound("cachePolicy", params.CachePolicyId)
		return
	}

	cachePolicy := &serverconfigs.HTTPCachePolicy{}
	err = json.Unmarshal(configJSON, cachePolicy)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["cachePolicy"] = cachePolicy

	// 其他选项
	this.Data["types"] = serverconfigs.AllCachePolicyTypes

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	CachePolicyId int64

	Name string
	Type string

	// file
	FileDir string

	CapacityJSON []byte
	MaxSizeJSON  []byte
	MaxKeys      int64

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	// 校验选项
	var options interface{}
	switch params.Type {
	case serverconfigs.CachePolicyTypeFile:
		params.Must.
			Field("fileDir", params.FileDir).
			Require("请输入缓存目录")
		options = &serverconfigs.HTTPFileCacheConfig{
			Dir: params.FileDir,
		}
	case serverconfigs.CachePolicyTypeMemory:
		options = &serverconfigs.HTTPMemoryCacheConfig{
		}
	default:
		this.Fail("请选择正确的缓存类型")
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPCachePolicyRPC().UpdateHTTPCachePolicy(this.AdminContext(), &pb.UpdateHTTPCachePolicyRequest{
		CachePolicyId: params.CachePolicyId,
		IsOn:          params.IsOn,
		Name:          params.Name,
		Description:   params.Description,
		CapacityJSON:  params.CapacityJSON,
		MaxKeys:       params.MaxKeys,
		MaxSizeJSON:   params.MaxSizeJSON,
		Type:          params.Type,
		OptionsJSON:   optionsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
