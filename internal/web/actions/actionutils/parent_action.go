package actionutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"strconv"
)

type ParentAction struct {
	actions.ActionObject

	rpcClient *rpc.RPCClient
}

// Parent 可以调用自身的一个简便方法
func (this *ParentAction) Parent() *ParentAction {
	return this
}

func (this *ParentAction) ErrorPage(err error) {
	if err == nil {
		return
	}

	// 日志
	this.CreateLog(oplogs.LevelError, "系统发生错误：%s", err.Error())

	if this.Request.Method == http.MethodGet {
		FailPage(this, err)
	} else {
		Fail(this, err)
	}
}

func (this *ParentAction) ErrorText(err string) {
	this.ErrorPage(errors.New(err))
}

func (this *ParentAction) NotFound(name string, itemId int64) {
	this.ErrorPage(errors.New(name + " id: '" + strconv.FormatInt(itemId, 10) + "' is not found"))
}

func (this *ParentAction) NewPage(total int64, size ...int64) *Page {
	if len(size) > 0 {
		return NewActionPage(this, total, size[0])
	}
	return NewActionPage(this, total, 10)
}

func (this *ParentAction) Nav(mainMenu string, tab string, firstMenu string) {
	this.Data["mainMenu"] = mainMenu
	this.Data["mainTab"] = tab
	this.Data["firstMenuItem"] = firstMenu
}

func (this *ParentAction) FirstMenu(menuItem string) {
	this.Data["firstMenuItem"] = menuItem
}

func (this *ParentAction) SecondMenu(menuItem string) {
	this.Data["secondMenuItem"] = menuItem
}

func (this *ParentAction) TinyMenu(menuItem string) {
	this.Data["tinyMenuItem"] = menuItem
}

func (this *ParentAction) AdminId() int64 {
	return this.Context.GetInt64("adminId")
}

func (this *ParentAction) CreateLog(level string, description string, args ...interface{}) {
	desc := fmt.Sprintf(description, args...)
	if level == oplogs.LevelInfo {
		if this.Code != 200 {
			level = oplogs.LevelWarn
			if len(this.Message) > 0 {
				desc += " 失败：" + this.Message
			}
		}
	}
	err := dao.SharedLogDAO.CreateAdminLog(this.AdminContext(), level, this.Request.URL.Path, desc, this.RequestRemoteIP())
	if err != nil {
		utils.PrintError(err)
	}
}

func (this *ParentAction) CreateLogInfo(description string, args ...interface{}) {
	this.CreateLog(oplogs.LevelInfo, description, args...)
}

// RPC 获取RPC
func (this *ParentAction) RPC() *rpc.RPCClient {
	if this.rpcClient != nil {
		return this.rpcClient
	}

	// 所有集群
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		logs.Fatal(err)
		return nil
	}
	this.rpcClient = rpcClient

	return rpcClient
}

// AdminContext 获取Context
func (this *ParentAction) AdminContext() context.Context {
	if this.rpcClient == nil {
		rpcClient, err := rpc.SharedRPC()
		if err != nil {
			logs.Fatal(err)
			return nil
		}
		this.rpcClient = rpcClient
	}
	return this.rpcClient.Context(this.AdminId())
}
