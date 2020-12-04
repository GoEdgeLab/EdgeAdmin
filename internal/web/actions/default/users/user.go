package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/users/userutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type UserAction struct {
	actionutils.ParentAction
}

func (this *UserAction) Init() {
	this.Nav("", "", "index")
}

func (this *UserAction) RunGet(params struct {
	UserId int64
}) {
	err := userutils.InitUser(this.Parent(), params.UserId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	userResp, err := this.RPC().UserRPC().FindEnabledUser(this.AdminContext(), &pb.FindEnabledUserRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	user := userResp.User
	if user == nil {
		this.NotFound("user", params.UserId)
		return
	}

	this.Data["user"] = maps.Map{
		"id":       user.Id,
		"username": user.Username,
		"fullname": user.Fullname,
		"email":    user.Email,
		"tel":      user.Tel,
		"remark":   user.Remark,
		"mobile":   user.Mobile,
		"isOn":     user.IsOn,
	}

	this.Show()
}
