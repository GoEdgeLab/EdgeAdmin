package security

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/securitymanager"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	config, err := securitymanager.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 国家和地区
	countryMaps := []maps.Map{}
	for _, countryId := range config.AllowCountryIds {
		countryResp, err := this.RPC().RegionCountryRPC().FindEnabledRegionCountry(this.AdminContext(), &pb.FindEnabledRegionCountryRequest{CountryId: countryId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		country := countryResp.Country
		if country != nil {
			countryMaps = append(countryMaps, maps.Map{
				"id":   country.Id,
				"name": country.Name,
			})
		}
	}
	this.Data["countries"] = countryMaps

	// 省份
	provinceMaps := []maps.Map{}
	for _, provinceId := range config.AllowProvinceIds {
		provinceResp, err := this.RPC().RegionProvinceRPC().FindEnabledRegionProvince(this.AdminContext(), &pb.FindEnabledRegionProvinceRequest{ProvinceId: provinceId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		province := provinceResp.Province
		if province != nil {
			provinceMaps = append(provinceMaps, maps.Map{
				"id":   province.Id,
				"name": province.Name,
			})
		}
	}
	this.Data["provinces"] = provinceMaps

	this.Data["config"] = config
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	Frame           string
	CountryIdsJSON  []byte
	ProvinceIdsJSON []byte
	AllowLocal      bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改管理界面安全设置")

	config, err := securitymanager.LoadSecurityConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 框架
	config.Frame = params.Frame

	// 国家和地区
	countryIds := []int64{}
	if len(params.CountryIdsJSON) > 0 {
		err = json.Unmarshal(params.CountryIdsJSON, &countryIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.AllowCountryIds = countryIds

	// 省份
	provinceIds := []int64{}
	if len(params.ProvinceIdsJSON) > 0 {
		err = json.Unmarshal(params.ProvinceIdsJSON, &provinceIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.AllowProvinceIds = provinceIds

	// 允许本地
	config.AllowLocal = params.AllowLocal

	err = securitymanager.UpdateSecurityConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
