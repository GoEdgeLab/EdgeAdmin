package setup

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(Helper)).
			Prefix("/setup").
			Get("", new(IndexAction)).
			Post("/validateApi", new(ValidateApiAction)).
			Post("/validateDb", new(ValidateDbAction)).
			Post("/validateAdmin", new(ValidateAdminAction)).
			Post("/install", new(InstallAction)).
			EndAll()
	})
}
