Tea.context(function () {
	this.teaweb = teaweb

	this.deleteNode = function (nodeId) {
		teaweb.confirm("确定要从当前集群中删除这个节点吗？", function () {
			this.$post("/nodes/delete")
				.params({
					clusterId: this.clusterId,
					nodeId: nodeId
				})
				.refresh();
		})
	}

	this.upNode = function (nodeId) {
		teaweb.confirm("确定要手动上线此节点吗？", function () {
			this.$post("/clusters/cluster/node/up")
				.params({
					nodeId: nodeId
				})
				.refresh()
		})
	}
})