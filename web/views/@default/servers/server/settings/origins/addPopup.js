Tea.context(function () {
	this.addr = ""
	this.protocol = ""
	this.isOSS = false

	this.addrError = ""

	if (this.isHTTP) {
		this.protocol = "http"
	} else if (this.serverType == "tcpProxy") {
		this.protocol = "tcp"
	} else if (this.serverType == "udpProxy") {
		this.protocol = "udp"
	}

	this.changeProtocol = function () {
		this.isOSS = this.protocol.startsWith("oss:")

		this.checkPort()
	}

	this.changeAddr = function () {
		if (this.serverType == "httpProxy") {
			if (this.addr.startsWith("http://")) {
				this.protocol = "http"
			} else if (this.addr.startsWith("https://")) {
				this.protocol = "https"
			}
		}

		this.checkPort()
	}

	this.checkPort = function () {
		this.addrError = ""

		// HTTP
		if (this.protocol == "http") {
			if (this.addr.endsWith(":443")) {
				this.addrError = "443通常是HTTPS协议端口，请确认源站协议选择是否正确。"
			} else if (this.addr.endsWith(":8443")) {
				this.addrError = "8443通常是HTTPS协议端口，请确认源站协议选择是否正确。"
			}
		}

		// HTTPS
		if (this.protocol == "https") {
			if (this.addr.endsWith(":80")) {
				this.addrError = "80通常是HTTP协议端口，请确认源站协议选择是否正确。"
			} else if (this.addr.endsWith(":8080")) {
				this.addrError = "8080通常是HTTP协议端口，请确认源站协议选择是否正确。"
			}
		}
	}
})