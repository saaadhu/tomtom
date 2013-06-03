/*
angular.module('scroll', []).directive('whenScrolled', function() {
    return function(scope, elm, attr) {
        var raw = elm[0];
        
        elm.bind('$scroll', function() {
            if (raw.scrollTop + raw.offsetHeight >= raw.scrollHeight - 100) {
                scope.$apply(attr.whenScrolled);
            }
        });
    };
});
*/

function TomTomCtrl ($scope, $http, $location, $anchorScroll) {
    $scope.feeds = [];
    $scope.items = [];
    $scope.pairs = [];
    $scope.new_url = '';

    var offset = 0;
    $scope.current_feed_id = '';
    $scope.last_loaded_feed_id = '';
    $scope.more_to_fetch = false;
    
    function fetchAndLoadFeedContents (id) {
        $scope.loading_feed_items = true;
        $scope.more_to_fetch = false;
        
        $http.get ("/feed/" + id + "?o=" + offset).success (function (data) {
            $scope.refreshFeeds();
            $scope.loading_feed_items = false;

            var i = 0;
            for (; i<data.length; ++i)
                $scope.items.push (data[i]);
        
            $scope.last_loaded_feed_id = id;
            
            offset += i;
            if (i != 0)
            {
                $location.hash(data[0].Id);
                $anchorScroll();
            }

            $scope.more_to_fetch = i == 5;
        });
    }
    
    $scope.showRecentFeeds = function() {
        $scope.loading_feed_items = true;
        $http.get ("/recent").success (function (data) {
            $scope.loading_feed_items = false;
            $scope.pairs = data;
        });
    }
    
    $scope.moreOfFeed = function () {
        fetchAndLoadFeedContents ($scope.last_loaded_feed_id);
    };
    
    $scope.refreshFeeds = function() {
        $scope.loading_feeds = true;
        $http.get ("/feeds").success (function (data) {
            $scope.loading_feeds = false; 
            $scope.feeds = data;
        });
    };
   
    $scope.loadFeed = function (id) {
        offset = 0;
        $scope.current_feed_id = id;
        fetchAndLoadFeedContents (id);
    };
    
    $scope.addUrl = function() {
        $scope.adding_feed = true;
        $http.post ("/feeds/add", { "url": $scope.new_url }).success (function (data) {
            $scope.new_url = '';
            $scope.adding_feed = false;
            $scope.feeds = data;
        });
    };
    
    $scope.refreshFeeds();
    $scope.showRecentFeeds();
}
