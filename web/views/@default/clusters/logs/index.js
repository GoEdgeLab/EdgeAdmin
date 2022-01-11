Tea.context(function () {
	this.$delay(function () {
		teaweb.datepicker("day-from-picker")
		teaweb.datepicker("day-to-picker")
	})

	this.updateRead = function (logId) {
		this.$post(".readLogs")
			.params({
				logIds: [logId]
			})
			.success(function () {
				teaweb.reload()
			})
	}

	this.updatePageRead = function () {
		let logIds = this.logs.map(function (v) {
			return v.id
		})
		teaweb.confirm("确定要设置本页日志为已读吗？", function () {
			this.$post(".readLogs")
				.params({
					logIds: logIds
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.updateAllRead = function () {
		teaweb.confirm("确定要设置所有日志为已读吗？", function () {
			this.$post(".readAllLogs")
				.params({})
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.changeCluster = function (clusterId) {
		this.clusterId = clusterId
	}
})