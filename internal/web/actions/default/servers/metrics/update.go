// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package metrics

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	ItemId int64
}) {
	item, err := InitItem(this.Parent(), params.ItemId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["valueDefinitions"] = serverconfigs.FindAllMetricValueDefinitions(item.Category)

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	ItemId     int64
	Name       string
	KeysJSON   []byte
	PeriodJSON []byte
	Value      string
	IsOn       bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入指标名称")

	// 统计对象
	if len(params.KeysJSON) == 0 {
		this.FailField("keys", "请选择指标统计的对象")
	}
	var keys = []string{}
	err := json.Unmarshal(params.KeysJSON, &keys)
	if err != nil {
		this.FailField("keys", "解析指标对象失败")
	}
	if len(keys) == 0 {
		this.FailField("keys", "请选择指标统计的对象")
	}

	var periodMap = maps.Map{}
	err = json.Unmarshal(params.PeriodJSON, &periodMap)
	if err != nil {
		this.FailField("period", "解析统计周期失败")
	}
	var period = periodMap.GetInt32("period")
	var periodUnit = periodMap.GetString("unit")

	_, err = this.RPC().MetricItemRPC().UpdateMetricItem(this.AdminContext(), &pb.UpdateMetricItemRequest{
		MetricItemId: params.ItemId,
		Name:         params.Name,
		Keys:         keys,
		Period:       period,
		PeriodUnit:   periodUnit,
		Value:        params.Value,
		IsOn:         params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer this.CreateLogInfo("修改统计指标 %d", params.ItemId)
	this.Success()
}
