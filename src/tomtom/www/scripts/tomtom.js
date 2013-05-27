function TomTomCtrl ($scope, $http) {
    $scope.feeds = [];
    $scope.items = [];
    $scope.new_url = '';
    
    $scope.refreshFeeds = function() {
        $http.get ("/feeds").success (function (data) {
            $scope.feeds = data;
        });
    };
   
    $scope.loadFeed = function (id) {
        $http.get ("/feed/" + id).success (function (data) {
            $scope.items = data;
        });
    }
    
    $scope.addUrl = function() {
        $http.post ("/feeds/add", { "url": $scope.new_url }).success (function (data) {
            $scope.feeds = data;
        });
    };
    
    $scope.refreshFeeds();
}
