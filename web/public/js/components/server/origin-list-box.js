Vue.component("origin-list-box", {
	props: ["v-primary-origins", "v-backup-origins", "v-server-type", "v-params"],
	data: function () {
		return {
			primaryOrigins: this.vPrimaryOrigins,
			backupOrigins: this.vBackupOrigins
		}
	},
	methods: {
		createPrimaryOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=primary&" + this.vParams, {
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		createBackupOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=backup&" + this.vParams, {
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		updateOrigin: function (originId, originType) {
			teaweb.popup("/servers/server/settings/origins/updatePopup?originType=" + originType + "&" + this.vParams + "&originId=" + originId, {
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		deleteOrigin: function (originId, originType) {
			let that = this
			teaweb.confirm("确定要删除此源站吗？", function () {
				Tea.action("/servers/server/settings/origins/delete?" + that.vParams + "&originId=" + originId + "&originType=" + originType)
					.post()
					.success(function () {
						teaweb.success("删除成功", function () {
							window.location.reload()
						})
					})
			})
		}
	},
	template: `<div>
	<h3>主要源站 <a href="" @click.prevent="createPrimaryOrigin()">[添加主要源站]</a> </h3>
	<p class="comment" v-if="primaryOrigins.length == 0">暂时还没有主要源站。</p>
	<origin-list-table v-if="primaryOrigins.length > 0" :v-origins="vPrimaryOrigins" :v-origin-type="'primary'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>

	<h3>备用源站 <a href="" @click.prevent="createBackupOrigin()">[添加备用源站]</a></h3>
	<p class="comment" v-if="backupOrigins.length == 0" :v-origins="primaryOrigins">暂时还没有备用源站。</p>
	<origin-list-table v-if="backupOrigins.length > 0" :v-origins="backupOrigins" :v-origin-type="'backup'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>
</div>`
})

Vue.component("origin-list-table", {
	props: ["v-origins", "v-origin-type"],
	data: function () {
		let hasMatchedDomains = false
		let origins = this.vOrigins
		if (origins != null && origins.length > 0) {
			origins.forEach(function (origin) {
				if (origin.domains != null && origin.domains.length > 0) {
					hasMatchedDomains = true
				}
			})
		}

		return {
			hasMatchedDomains: hasMatchedDomains
		}
	},
	methods: {
		deleteOrigin: function (originId) {
			this.$emit("deleteOrigin", originId, this.vOriginType)
		},
		updateOrigin: function (originId) {
			this.$emit("updateOrigin", originId, this.vOriginType)
		}
	},
	template: `
<table class="ui table selectable">
	<thead>
		<tr>
			<th>源站地址</th>
			<th>权重</th>
			<th class="width10">状态</th>
			<th class="two op">操作</th>
		</tr>	
	</thead>
	<tbody>
		<tr v-for="origin in vOrigins">
			<td :class="{disabled:!origin.isOn}">
				<a href="" @click.prevent="updateOrigin(origin.id)" :class="{disabled:!origin.isOn}">{{origin.addr}} &nbsp;<i class="icon expand small"></i></a>
				<div style="margin-top: 0.3em">
					<tiny-basic-label v-if="origin.isOSS"><i class="icon hdd outline"></i>对象存储</tiny-basic-label>
					<tiny-basic-label v-if="origin.name.length > 0">{{origin.name}}</tiny-basic-label>
					<tiny-basic-label v-if="origin.hasCert">证书</tiny-basic-label>
					<tiny-basic-label v-if="origin.host != null && origin.host.length > 0">主机名: {{origin.host}}</tiny-basic-label>
					<tiny-basic-label v-if="origin.followPort">端口跟随</tiny-basic-label>
					<tiny-basic-label v-if="origin.addr != null && origin.addr.startsWith('https://') && origin.http2Enabled">HTTP/2</tiny-basic-label>
	
					<span v-if="origin.domains != null && origin.domains.length > 0"><tiny-basic-label v-for="domain in origin.domains">匹配: {{domain}}</tiny-basic-label></span>
					<span v-else-if="hasMatchedDomains"><tiny-basic-label>匹配: 所有域名</tiny-basic-label></span>
				</div>
			</td>
			<td :class="{disabled:!origin.isOn}">{{origin.weight}}</td>
			<td>
				<label-on :v-is-on="origin.isOn"></label-on>
			</td>
			<td>
				<a href="" @click.prevent="updateOrigin(origin.id)">修改</a> &nbsp;
				<a href="" @click.prevent="deleteOrigin(origin.id)">删除</a>
			</td>
		</tr>
	</tbody>
</table>`
})