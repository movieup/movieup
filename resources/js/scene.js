new Vue({
  el: '#tabs',
  data: {
    activeTabNumber: 2
  },
  methods: {
    clickTab: function(tabNumber){
      this.activeTabNumber = tabNumber
    }
  }
})
