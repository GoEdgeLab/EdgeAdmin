Tea.context(function () {
	this.params = [
		{
			"name": "DOCUMENT_ROOT",
			"value": "",
			"nameZh": "文档目录"
		},
		{
			"name": "SCRIPT_FILENAME",
			"value": "",
			"nameZh": "脚本文件"
		}
	]

	this.addParam = function () {
		this.params.push({
			"name": "",
			"value": "",
			"nameZh": ""
		})
		this.$delay(function () {
			this.$find("form input[name='paramNames']").last().focus()
		})
	}

	this.removeParam = function (index) {
		this.params.$remove(index)
	}
})