Vue.component("node-log-row", {
	props: ["v-log", "v-keyword"],
	data: function () {
		return {
			log: this.vLog,
			keyword: this.vKeyword
		}
	},
	template: `<div>
	<pre class="log-box" style="margin: 0; padding: 0"><span :class="{red:log.level == 'error', orange:log.level == 'warning', green: log.level == 'success'}"><span v-if="!log.isToday">[{{log.createdTime}}]</span><strong v-if="log.isToday">[{{log.createdTime}}]</strong><keyword :v-word="keyword">[{{log.tag}}]{{log.description}}</keyword></span> &nbsp; <span v-if="log.count > 0" class="ui label tiny" :class="{red:log.level == 'error', orange:log.level == 'warning'}">共{{log.count}}条</span></pre>
</div>`
})