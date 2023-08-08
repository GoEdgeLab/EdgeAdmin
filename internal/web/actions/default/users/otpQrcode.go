package users

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type OtpQrcodeAction struct {
	actionutils.ParentAction
}

func (this *OtpQrcodeAction) Init() {
	this.Nav("", "", "")
}

func (this *OtpQrcodeAction) RunGet(params struct {
	UserId int64
}) {
	loginResp, err := this.RPC().LoginRPC().FindEnabledLogin(this.AdminContext(), &pb.FindEnabledLoginRequest{
		UserId: params.UserId,
		Type:   "otp",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	login := loginResp.Login
	if login == nil || !login.IsOn {
		this.NotFound("userLogin", params.UserId)
		return
	}

	loginParams := maps.Map{}
	err = json.Unmarshal(login.ParamsJSON, &loginParams)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	secret := loginParams.GetString("secret")

	// 当前用户信息
	userResp, err := this.RPC().UserRPC().FindEnabledUser(this.AdminContext(), &pb.FindEnabledUserRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var user = userResp.User
	if user == nil {
		this.NotFound("user", params.UserId)
		return
	}

	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var productName = uiConfig.ProductName
	if len(productName) == 0 {
		productName = "GoEdge用户"
	}
	var url = gotp.NewDefaultTOTP(secret).ProvisioningUri(user.Username, productName)
	data, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.AddHeader("Content-Type", "image/png")
	_, _ = this.Write(data)
}
