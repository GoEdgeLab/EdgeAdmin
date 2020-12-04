package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/users/userutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
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

func (this *UpdateAction) RunPost(params struct {
	UserId   int64
	Username string
	Pass1    string
	Pass2    string
	Fullname string
	Mobile   string
	Tel      string
	Email    string
	Remark   string
	IsOn     bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改用户 %d", params.UserId)

	params.Must.
		Field("username", params.Username).
		Require("请输入用户名").
		Match(`^[a-zA-Z0-9_]+$`, "用户名中只能含有英文、数字和下划线")

	checkUsernameResp, err := this.RPC().UserRPC().CheckUsername(this.AdminContext(), &pb.CheckUsernameRequest{
		UserId:   params.UserId,
		Username: params.Username,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if checkUsernameResp.Exists {
		this.FailField("username", "此用户名已经被占用，请换一个")
	}

	if len(params.Pass1) > 0 {
		params.Must.
			Field("pass1", params.Pass1).
			Require("请输入密码").
			Field("pass2", params.Pass2).
			Require("请再次输入确认密码").
			Equal(params.Pass1, "两次输入的密码不一致")
	}

	params.Must.
		Field("fullname", params.Fullname).
		Require("请输入全名")

	if len(params.Mobile) > 0 {
		params.Must.
			Field("mobile", params.Mobile).
			Mobile("请输入正确的手机号")
	}
	if len(params.Email) > 0 {
		params.Must.
			Field("email", params.Email).
			Email("请输入正确的电子邮箱")
	}

	_, err = this.RPC().UserRPC().UpdateUser(this.AdminContext(), &pb.UpdateUserRequest{
		UserId:   params.UserId,
		Username: params.Username,
		Password: params.Pass1,
		Fullname: params.Fullname,
		Mobile:   params.Mobile,
		Tel:      params.Tel,
		Email:    params.Email,
		Remark:   params.Remark,
		IsOn:     params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
