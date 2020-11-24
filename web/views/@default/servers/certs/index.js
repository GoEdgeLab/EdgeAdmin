Tea.context(function () {
	// 上传证书
	this.uploadCert = function () {
		teaweb.popup("/servers/certs/uploadPopup", {
			height: "28em",
			callback: function () {
				teaweb.success("上传成功", function () {
					window.location.reload()
				})
			}
		})
	}

	// 删除证书
	this.deleteCert = function (certId) {
		let that = this
		teaweb.confirm("确定要删除此证书吗？", function () {
			that.$post("/servers/certs/delete")
				.params({
					certId: certId
				})
				.refresh()
		})
	}

	// 查看证书
	this.viewCert = function (certId) {
		teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
			height: "28em",
			width: "48em"
		})
	}

	// 修改证书
	this.updateCert = function (certId) {
		teaweb.popup("/servers/certs/updatePopup?certId=" + certId, {
			height: "28em",
			callback: function () {
				teaweb.success("上传成功", function () {
					window.location.reload()
				})
			}
		})
	}
})