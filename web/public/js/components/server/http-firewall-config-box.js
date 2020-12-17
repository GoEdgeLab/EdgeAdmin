Vue.component("http-firewall-config-box", {
	props: ["v-firewall-config", "v-is-location", "v-firewall-policy"],
	data: function () {
		let firewall = this.vFirewallConfig
		if (firewall == null) {
			firewall = {
				isPrior: false,
				isOn: false,
				firewallPolicyId: 0
			}
		}

		return {
			firewall: firewall
		}
	},
	template: `<div>
	<input type="hidden" name="firewallJSON" :value="JSON.stringify(firewall)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="firewall" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || firewall.isPrior">
			<tr>
				<td>WAF策略</td>
				<td>
					<div v-if="vFirewallPolicy != null">{{vFirewallPolicy.name}}
						<p class="comment">使用当前服务所在集群的设置。</p>
					</div>
					<span v-else class="red">当前集群没有设置WAF策略，当前配置无法生效。</span>
				</td>
			</tr>
			<tr>
				<td class="title">是否启用Web防火墙</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="firewall.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})