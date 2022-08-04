package actionutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	rpcerrors "github.com/TeaOSLab/EdgeCommon/pkg/rpc/errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// Fail 提示服务器错误信息
func Fail(action actions.ActionWrapper, err error) {
	if err != nil {
		logs.Println("[" + reflect.TypeOf(action).String() + "]" + findStack(err.Error()))
	}
	action.Object().Fail(teaconst.ErrServer + "（" + err.Error() + "）")
}

// FailPage 提示页面错误信息
func FailPage(action actions.ActionWrapper, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}

	logs.Println("[" + reflect.TypeOf(action).String() + "]" + findStack(err.Error()))

	// 当前API终端地址
	var apiEndpoints = []string{}
	apiConfig, apiConfigErr := configs.LoadAPIConfig()
	if apiConfigErr == nil && apiConfig != nil {
		apiEndpoints = append(apiEndpoints, apiConfig.RPC.Endpoints...)
	}

	var isRPCConnError bool
	err, isRPCConnError = rpcerrors.HumanError(err, apiEndpoints, Tea.ConfigFile("api.yaml"))

	action.Object().ResponseWriter.WriteHeader(http.StatusInternalServerError)
	if len(action.Object().Request.Header.Get("X-Requested-With")) > 0 {
		action.Object().WriteString(teaconst.ErrServer)
	} else {
		// 本地的一些错误提示
		var isLocalAPI = false
		if isRPCConnError {
			host, _, hostErr := net.SplitHostPort(action.Object().Request.Host)
			if hostErr == nil {
				for _, endpoint := range apiEndpoints {
					if strings.HasPrefix(endpoint, "http://"+host) || strings.HasPrefix(endpoint, "https://"+host) || strings.HasPrefix(endpoint, host) {
						isLocalAPI = true
						break
					}
				}
			}
		}

		var issuesHTML = ""
		if isLocalAPI {
			// 读取本地API节点的issues
			issuesData, issuesErr := os.ReadFile(Tea.Root + "/edge-api/logs/issues.log")
			if issuesErr == nil {
				var issueMaps = []maps.Map{}
				issuesErr = json.Unmarshal(issuesData, &issueMaps)
				if issuesErr == nil && len(issueMaps) > 0 {
					var issueMap = issueMaps[0]
					issuesHTML = "本地API节点启动错误：" + issueMap.GetString("message") + "，处理建议：" + issueMap.GetString("suggestion")
				}
			}
		}

		var html = `<!DOCTYPE html>
<html>
	<head>
		<title>有系统错误需要处理</title>
		<meta charset="UTF-8"/>
		<style type="text/css">
		hr { border-top: 1px #ccc solid; }
		.red { color: red; }
		</style>
	</head>
	<body>
	<div style="background: #eee; border: 1px #ccc solid; padding: 10px; font-size: 12px; line-height: 1.8">
	` + teaconst.ErrServer + `
		<div>可以通过查看 <strong><em>$安装目录/logs/run.log</em></strong> 日志文件查看具体的错误提示。</div>
		<hr/>
		<div class="red">Error: ` + err.Error() + `</div>`

		if len(issuesHTML) > 0 {
			html += `        <hr/>
        <div class="red">` + issuesHTML + `</div>`
		}

		action.Object().WriteString(html + `
		</div>
	</body>
</html>`)
	}
}

// MatchPath 判断动作的文件路径是否相当
func MatchPath(action *actions.ActionObject, path string) bool {
	return action.Request.URL.Path == path
}

// FindParentAction 查找父级Action
func FindParentAction(actionPtr actions.ActionWrapper) *ParentAction {
	parentActionValue := reflect.ValueOf(actionPtr).Elem().FieldByName("ParentAction")
	if parentActionValue.IsValid() {
		parentAction, isOk := parentActionValue.Interface().(ParentAction)
		if isOk {
			return &parentAction
		}
	}
	return nil
}

func findStack(err string) string {
	_, currentFilename, _, currentOk := runtime.Caller(1)
	if currentOk {
		for i := 1; i < 32; i++ {
			_, filename, lineNo, ok := runtime.Caller(i)
			if !ok {
				break
			}

			if filename == currentFilename || filepath.Base(filename) == "parent_action.go" {
				continue
			}

			goPath := os.Getenv("GOPATH")
			if len(goPath) > 0 {
				absGoPath, err := filepath.Abs(goPath)
				if err == nil {
					filename = strings.TrimPrefix(filename, absGoPath)[1:]
				}
			} else if strings.Contains(filename, "src") {
				filename = filename[strings.Index(filename, "src"):]
			}

			err += "\n\t\t" + string(filename) + ":" + fmt.Sprintf("%d", lineNo)

			break
		}
	}

	return err
}
