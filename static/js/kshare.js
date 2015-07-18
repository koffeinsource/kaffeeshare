function getNamespaceRSS(namespace) {
    var RSSFeedURL = "https://kaffeeshare.appspot.com/k/show/rss/"+ encodeURIComponent(namespace);
    return RSSFeedURL;
}

var language = window.navigator.userLanguage || window.navigator.language;
moment.locale(language);

// global variables to store the request
// this is done so we can cancel the request
// TODO move it into jquery data attribute?
var request;
var db_cursor;
