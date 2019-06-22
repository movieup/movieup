var vueHeader = new Vue({
	el: '#header',
	data: {
	 	isShownMenu: false
	},
    methods: {
		clickHamburgerOverlay: function(event) {
			hamburgerBtn.click();
		}
    }
})

var hamburgerBtn;
// onload
document.addEventListener('DOMContentLoaded', function() {
	hamburgerBtn = document.getElementById("hamburger_btn");
	// ハンバーガーにクリックイベントを付与
	hamburgerBtn.addEventListener('click', function (event) {
		vueHeader.isShownMenu = !vueHeader.isShownMenu;
	});
});
