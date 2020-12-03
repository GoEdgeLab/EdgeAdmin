Tea.context(function () {
	this.step = "prepare"

	/**
	 * 准备工作
	 */
	this.authType = "http"

	this.doPrepare = function () {
		this.step = "user"
	}

	/**
	 * 选择用户
	 */
	this.userId = 0

	this.goPrepare = function () {
		this.step = "prepare"
	}

	this.createUser = function () {
		let that = this
		teaweb.popup("/servers/certs/acme/users/createPopup", {
			callback: function (resp) {
				teaweb.successToast("创建成功")

				let acmeUser = resp.data.acmeUser
				let description = acmeUser.description
				if (description.length > 0) {
					description = "（" + description + "）"
				}
				that.userId = acmeUser.id
				that.users.unshift({
					id: acmeUser.id,
					description: description,
					email: acmeUser.email
				})
			}
		})
	}

	this.doUser = function () {
		if (this.userId == 0) {
			teaweb.warn("请选择一个申请证书的用户")
			return
		}
		this.step = "dns"
	}

	/**
	 * 设置DNS解析
	 */
	this.dnsProviderId = 0
	this.dnsDomain = ""
	this.autoRenew = true
	this.domains = []
	this.taskId = 0
	this.isRequesting = false

	this.goUser = function () {
		this.step = "user"
	}

	this.changeDomains = function (v) {
		this.domains = v
	}

	this.doDNS = function () {
		this.isRequesting = true
		let that = this

		this.$post("$")
			.params({
				authType: this.authType,
				acmeUserId: this.userId,
				dnsProviderId: this.dnsProviderId,
				dnsDomain: this.dnsDomain,
				domains: this.domains,
				autoRenew: this.autoRenew ? 1 : 0,
				taskId: this.taskId
			})
			.success(function (resp) {
				this.taskId = resp.data.taskId

				this.isRequesting = true
				this.$post(".run")
					.timeout(300)
					.params({
						taskId: this.taskId
					})
					.success(function (resp) {
						that.certId = resp.data.certId
						that.step = "finish"
					})
					.done(function () {
						that.isRequesting = false
					})
			})
			.done(function () {
				this.isRequesting = false
			})
	}

	/**
	 * 完成
	 */
	this.certId = 0

	this.goDNS = function () {
		this.step = "dns"
	}

	this.doFinish = function () {
		window.location = "/servers/certs/acme"
	}

	this.viewCert = function () {
		teaweb.popup("/servers/certs/certPopup?certId=" + this.certId, {
			height: "28em",
			width: "48em"
		})
	}
})