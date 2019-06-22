new Vue({
  el: '#tabs',
  data: {
    activeTabNumber: 0
  },
  methods: {
    clickTab: function(tabNumber){
      this.activeTabNumber = tabNumber
    }
  }
})
