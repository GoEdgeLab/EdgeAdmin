package teaconst

const (
	Version = "0.6.4.1"

	APINodeVersion = "0.6.4.1"

	ProductName   = "Edge Admin"
	ProcessName   = "edge-admin"
	ProductNameZH = "Edge"

	Role = "admin"

	EncryptKey    = "8f983f4d69b83aaa0d74b21a212f6967"
	EncryptMethod = "aes-256-cfb"

	ErrServer = "服务器出了点小问题，请联系技术人员处理。"
	CookieSID = "edgesid"

	SystemdServiceName = "edge-admin"
	UpdatesURL         = "https://goedge.cn/api/boot/versions?os=${os}&arch=${arch}&version=${version}"
)
