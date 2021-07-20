// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package records

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "record")
}

func (this *IndexAction) RunGet(params struct {
	DomainId int64
	Type     string
	Keyword  string
	RouteId  int64
}) {
	this.Data["type"] = params.Type
	this.Data["keyword"] = params.Keyword
	this.Data["routeId"] = params.RouteId

	// 域名信息
	domainResp, err := this.RPC().NSDomainRPC().FindEnabledNSDomain(this.AdminContext(), &pb.FindEnabledNSDomainRequest{NsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := domainResp.NsDomain
	if domain == nil {
		this.NotFound("nsDomain", params.DomainId)
		return
	}
	this.Data["domain"] = maps.Map{
		"id":   domain.Id,
		"name": domain.Name,
	}

	// 记录
	countResp, err := this.RPC().NSRecordRPC().CountAllEnabledNSRecords(this.AdminContext(), &pb.CountAllEnabledNSRecordsRequest{
		NsDomainId: params.DomainId,
		Type:       params.Type,
		NsRouteId:  params.RouteId,
		Keyword:    params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	recordsResp, err := this.RPC().NSRecordRPC().ListEnabledNSRecords(this.AdminContext(), &pb.ListEnabledNSRecordsRequest{
		NsDomainId: params.DomainId,
		Type:       params.Type,
		NsRouteId:  params.RouteId,
		Keyword:    params.Keyword,
		Offset:     page.Offset,
		Size:       page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var recordMaps = []maps.Map{}
	for _, record := range recordsResp.NsRecords {
		routeMaps := []maps.Map{}
		for _, route := range record.NsRoutes {
			routeMaps = append(routeMaps, maps.Map{
				"id":   route.Id,
				"name": route.Name,
			})
		}

		recordMaps = append(recordMaps, maps.Map{
			"id":          record.Id,
			"name":        record.Name,
			"type":        record.Type,
			"value":       record.Value,
			"ttl":         record.Ttl,
			"weight":      record.Weight,
			"description": record.Description,
			"isOn":        record.IsOn,
			"routes":      routeMaps,
		})
	}
	this.Data["records"] = recordMaps

	// 所有记录类型
	this.Data["types"] = dnsconfigs.FindAllRecordTypeDefinitions()

	this.Show()
}
