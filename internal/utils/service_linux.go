// +build linux

package utils

import (
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var systemdServiceFile = "/etc/systemd/system/edge-admin.service"
var initServiceFile = "/etc/init.d/" + teaconst.SystemdServiceName

// 安装服务
func (this *ServiceManager) Install(exePath string, args []string) error {
	if os.Getgid() != 0 {
		return errors.New("only root users can install the service")
	}

	systemd, err := exec.LookPath("systemctl")
	if err != nil {
		return this.installInitService(exePath, args)
	}

	return this.installSystemdService(systemd, exePath, args)
}

// 启动服务
func (this *ServiceManager) Start() error {
	if os.Getgid() != 0 {
		return errors.New("only root users can start the service")
	}

	if files.NewFile(systemdServiceFile).Exists() {
		systemd, err := exec.LookPath("systemctl")
		if err != nil {
			return err
		}

		return exec.Command(systemd, "start", teaconst.SystemdServiceName+".service").Start()
	}
	return exec.Command("service", teaconst.ProcessName, "start").Start()
}

// 删除服务
func (this *ServiceManager) Uninstall() error {
	if os.Getgid() != 0 {
		return errors.New("only root users can uninstall the service")
	}

	if files.NewFile(systemdServiceFile).Exists() {
		systemd, err := exec.LookPath("systemctl")
		if err != nil {
			return err
		}

		// disable service
		exec.Command(systemd, "disable", teaconst.SystemdServiceName+".service").Start()

		// reload
		exec.Command(systemd, "daemon-reload")

		return files.NewFile(systemdServiceFile).Delete()
	}

	f := files.NewFile(initServiceFile)
	if f.Exists() {
		return f.Delete()
	}
	return nil
}

// install init service
func (this *ServiceManager) installInitService(exePath string, args []string) error {
	shortName := teaconst.SystemdServiceName
	scriptFile := Tea.Root + "/scripts/" + shortName
	if !files.NewFile(scriptFile).Exists() {
		return errors.New("'scripts/" + shortName + "' file not exists")
	}

	data, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		return err
	}

	data = regexp.MustCompile("INSTALL_DIR=.+").ReplaceAll(data, []byte("INSTALL_DIR="+Tea.Root))
	err = ioutil.WriteFile(initServiceFile, data, 0777)
	if err != nil {
		return err
	}

	chkCmd, err := exec.LookPath("chkconfig")
	if err != nil {
		return err
	}

	err = exec.Command(chkCmd, "--add", teaconst.ProcessName).Start()
	if err != nil {
		return err
	}

	return nil
}

// install systemd service
func (this *ServiceManager) installSystemdService(systemd, exePath string, args []string) error {
	shortName := teaconst.SystemdServiceName
	longName := "GoEdge API" // TODO 将来可以修改

	desc := `# Provides:          ` + shortName + `
# Required-Start:    $all
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:
# Short-Description: ` + longName + ` Service
### END INIT INFO

[Unit]
Description=` + longName + ` Service
Before=shutdown.target
After=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=1s
ExecStart=` + exePath + ` daemon
ExecStop=` + exePath + ` stop
ExecReload=` + exePath + ` reload

[Install]
WantedBy=multi-user.target`

	// write file
	err := ioutil.WriteFile(systemdServiceFile, []byte(desc), 0777)
	if err != nil {
		return err
	}

	// stop current systemd service if running
	exec.Command(systemd, "stop", shortName+".service")

	// reload
	exec.Command(systemd, "daemon-reload")

	// enable
	cmd := exec.Command(systemd, "enable", shortName+".service")
	return cmd.Run()
}
