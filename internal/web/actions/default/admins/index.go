package admins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().AdminRPC().CountAllEnabledAdmins(this.AdminContext(), &pb.CountAllEnabledAdminsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	adminsResp, err := this.RPC().AdminRPC().ListEnabledAdmins(this.AdminContext(), &pb.ListEnabledAdminsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	adminMaps := []maps.Map{}
	for _, admin := range adminsResp.Admins {
		adminMaps = append(adminMaps, maps.Map{
			"id":          admin.Id,
			"isOn":        admin.IsOn,
			"isSuper":     admin.IsSuper,
			"username":    admin.Username,
			"fullname":    admin.Fullname,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", admin.CreatedAt),
		})
	}
	this.Data["admins"] = adminMaps

	this.Show()
}
