package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"strings"
)

type ValidateDbAction struct {
	actionutils.ParentAction
}

func (this *ValidateDbAction) RunPost(params struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string

	Must *actions.Must
}) {
	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Match(`^[\w\.-]+$`, "主机地址中不能包含特殊字符").
		Field("port", params.Port).
		Require("请输入端口").
		Match(`^\d+$`, "端口中只能包含数字").
		Field("database", params.Database).
		Require("请输入数据库名称").
		Match(`^[\w\.-]+$`, "数据库名称中不能包含特殊字符").
		Field("username", params.Username).
		Require("请输入连接数据库的用户名").
		Match(`^[\w\.-]+$`, "用户名中不能包含特殊字符")

	// 测试连接
	db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
		Driver: "mysql",
		Dsn:    params.Username + ":" + params.Password + "@tcp(" + params.Host + ":" + params.Port + ")/" + params.Database,
		Prefix: "",
	})
	if err != nil {
		this.Fail("数据库信息错误：" + err.Error())
	}

	err = db.Raw().Ping()
	if err != nil {
		// 是否是数据库不存在
		if strings.Contains(err.Error(), "Error 1049") {
			db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
				Driver: "mysql",
				Dsn:    params.Username + ":" + params.Password + "@tcp(" + params.Host + ":" + params.Port + ")/",
				Prefix: "",
			})

			_, err = db.Exec("CREATE DATABASE `" + params.Database + "`")
			if err != nil {
				this.Fail("尝试创建数据库失败：" + err.Error())
			}
		} else {
			this.Fail("无法连接到数据库，请检查配置：" + err.Error())
		}
	}

	// 检查数据库版本
	one, err := db.FindOne("SELECT VERSION() AS v")
	if err != nil {
		this.Fail("检查数据库版本时出错：" + err.Error())
	}
	if one == nil {
		this.Fail("检查数据库版本时出错：无法获取数据库版本")
	}
	version := one.GetString("v")
	if stringutil.VersionCompare(version, "5.7.8") < 0 {
		this.Fail("数据库版本至少在v5.7.8以上，你现在使用的是v" + version)
	}

	this.Data["db"] = maps.Map{
		"host":         params.Host,
		"port":         params.Port,
		"database":     params.Database,
		"username":     params.Username,
		"password":     params.Password,
		"passwordMask": strings.Repeat("*", len(params.Password)),
	}

	this.Success()
}
