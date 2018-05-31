var apiURL = 'http://localhost/v1/users/'

var demo = new Vue({

  el: '#app',

  data: {
    userId: null,
    user: null,
    repositories: null,
    statistics: null
  },

  computed: {
    isUserIdValid: function() {
      return this.userId != null && this.userId.length > 0;
    },
    anyUserData: function() {
      return this.user != null;
    },
    statisticsBySizeDesc: function () {
      return _.sortBy(this.statistics, 'Size').reverse();
    },
    repositoriesByName: function() {
      return _.sortBy(this.repositories, function (i) { return i.name.toLowerCase();} );
    },
    anyRepos: function() {
      return this.repositories != null && this.repositories.length > 0;
    }
  },

  methods: {
    fetchData: function () {
      if (!this.isUserIdValid) {
        return;
      }

      $.LoadingOverlay("show");
      var xhr = new XMLHttpRequest();
      var self = this;
      xhr.open('GET', apiURL + self.userId);
      xhr.onload = function () {
        self.user = JSON.parse(xhr.responseText);
        $.LoadingOverlay("hide");

        if (isUserDataEmpty(self)) {
          resetUserData(self)
          showNoDataAlert();
        } else {
          assignUserData(self);
        }
      }
      
      xhr.send()
    }
  }
})

function isUserDataEmpty(data) {
  return data.user.login === "";
}

function assignUserData(data) {
  data.repositories = data.user.Repositories;
  data.statistics = data.user.Statistics;
}

function resetUserData(data) {
  data.user = null;
  data.repositories = null;
  data.statistics = null;
}

function showNoDataAlert() {
  bootbox.alert({ 
    size: "small",
    message: "No data to display"
  });
}
