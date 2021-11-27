package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type CheckAction struct {
	actionutils.ParentAction
}

func (this *CheckAction) RunPost(params struct{}) {
	resp, err := this.RPC().NodeTaskRPC().ExistsNodeTasks(this.AdminContext(), &pb.ExistsNodeTasksRequest{
		ExcludeTypes: []string{"ipItemChanged"},
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["isDoing"] = resp.ExistTasks
	this.Data["hasError"] = resp.ExistError

	this.Success()
}
