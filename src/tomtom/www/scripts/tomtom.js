function TomTomCtrl ($scope, $http) {
    $scope.feeds = [];
    $scope.items = [];
    $scope.new_url = '';
    
    $scope.refreshFeeds = function() {
        $scope.loading_feeds = true;
        $http.get ("/feeds").success (function (data) {
            $scope.loading_feeds = false; 
            $scope.feeds = data;
        });
    };
   
    $scope.loadFeed = function (id) {
        $scope.loading_feed_items = true;
        $http.get ("/feed/" + id).success (function (data) {
            $scope.loading_feed_items = false;
            $scope.items = data;
            $scope.refreshFeeds();
        });
    }
    
    $scope.addUrl = function() {
        $scope.adding_feed = true;
        $http.post ("/feeds/add", { "url": $scope.new_url }).success (function (data) {
            $scope.new_url = '';
            $scope.adding_feed = true;
            $scope.feeds = data;
        });
    };
    
    $scope.refreshFeeds();
}
