/**
 * 一级菜单
 */
Vue.component("first-menu", {
	template: ' \
		<div class="first-menu"> \
			<div class="ui menu text blue small">\
				<slot></slot>\
			</div> \
			<div class="ui divider"></div> \
			<div class="margin"></div> \
		</div>'
});