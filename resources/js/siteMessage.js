new Vue({
	el: '#site_message',
	data: {
		show: false
	},
	methods: {
		handleScroll() {
			if (!this.show) {
				var top = this.$el.getBoundingClientRect().top;
				this.show = top < window.innerHeight + 100;
			}
		}
	},
	created() {
		window.addEventListener('scroll', this.handleScroll);
	},
	destroyed() {
		window.removeEventListener('scroll', this.handleScroll);
	},
	mounted () {
		this.handleScroll();
	}
});
