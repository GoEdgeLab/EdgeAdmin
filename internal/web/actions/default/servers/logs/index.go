// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package logs

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"regexp"
	"strings"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
	NodeId    int64
	Day       string
	Hour      string
	Keyword   string
	Ip        string
	Domain    string
	HasError  int
	HasWAF    int

	RequestId string
	ServerId  int64

	PageSize int64
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Y-m-d")
	}

	this.Data["clusterId"] = params.ClusterId
	this.Data["nodeId"] = params.NodeId
	this.Data["serverId"] = 0
	this.Data["path"] = this.Request.URL.Path
	this.Data["day"] = params.Day
	this.Data["hour"] = params.Hour
	this.Data["keyword"] = params.Keyword
	this.Data["ip"] = params.Ip
	this.Data["domain"] = params.Domain
	this.Data["accessLogs"] = []interface{}{}
	this.Data["hasError"] = params.HasError
	this.Data["hasWAF"] = params.HasWAF
	this.Data["pageSize"] = params.PageSize
	this.Data["isSlowQuery"] = false
	this.Data["slowQueryCost"] = ""

	day := params.Day
	ipList := []string{}

	if len(day) > 0 && regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).MatchString(day) {
		day = strings.ReplaceAll(day, "-", "")
		size := params.PageSize
		if size < 1 {
			size = 20
		}

		this.Data["hasError"] = params.HasError

		var before = time.Now()
		resp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
			RequestId:         params.RequestId,
			NodeClusterId:     params.ClusterId,
			NodeId:            params.NodeId,
			ServerId:          params.ServerId,
			HasError:          params.HasError > 0,
			HasFirewallPolicy: params.HasWAF > 0,
			Day:               day,
			HourFrom:          params.Hour,
			HourTo:            params.Hour,
			Keyword:           params.Keyword,
			Ip:                params.Ip,
			Domain:            params.Domain,
			Size:              size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		var cost = time.Since(before).Seconds()
		if cost > 5 {
			this.Data["slowQueryCost"] = fmt.Sprintf("%.2f", cost)
			this.Data["isSlowQuery"] = true
		}

		if len(resp.HttpAccessLogs) == 0 {
			this.Data["accessLogs"] = []interface{}{}
		} else {
			this.Data["accessLogs"] = resp.HttpAccessLogs
			for _, accessLog := range resp.HttpAccessLogs {
				if len(accessLog.RemoteAddr) > 0 {
					if !lists.ContainsString(ipList, accessLog.RemoteAddr) {
						ipList = append(ipList, accessLog.RemoteAddr)
					}
				}
			}
		}
		this.Data["hasMore"] = resp.HasMore
		this.Data["nextRequestId"] = resp.RequestId

		// 上一个requestId
		this.Data["hasPrev"] = false
		this.Data["lastRequestId"] = ""
		if len(params.RequestId) > 0 {
			this.Data["hasPrev"] = true
			prevResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
				RequestId:         params.RequestId,
				NodeClusterId:     params.ClusterId,
				NodeId:            params.NodeId,
				ServerId:          params.ServerId,
				HasError:          params.HasError > 0,
				HasFirewallPolicy: params.HasWAF > 0,
				Day:               day,
				Keyword:           params.Keyword,
				Ip:                params.Ip,
				Domain:            params.Domain,
				Size:              size,
				Reverse:           true,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if int64(len(prevResp.HttpAccessLogs)) == size {
				this.Data["lastRequestId"] = prevResp.RequestId
			}
		}
	}

	// 根据IP查询区域
	regionMap := map[string]string{} // ip => region
	if len(ipList) > 0 {
		resp, err := this.RPC().IPLibraryRPC().LookupIPRegions(this.AdminContext(), &pb.LookupIPRegionsRequest{IpList: ipList})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if resp.IpRegionMap != nil {
			for ip, region := range resp.IpRegionMap {
				regionMap[ip] = region.Summary
			}
		}
	}
	this.Data["regions"] = regionMap

	this.Show()
}
