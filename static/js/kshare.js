function getNamespaceRSS(namespace) {
    var RSSFeedURL = "https://kaffeeshare.appspot.com/k/show/rss/"+ encodeURIComponent(namespace);
    return RSSFeedURL;
}

var language = window.navigator.userLanguage || window.navigator.language;
moment.locale(language);

var spinnerOpts = {
		lines: 13, // The number of lines to draw
		length: 4, // The length of each line
		width: 2, // The line thickness
		radius: 4, // The radius of the inner circle
		corners: 1, // Corner roundness (0..1)
		rotate: 0, // The rotation offset
		color: '#000', // #rgb or #rrggbb
		speed: 1, // Rounds per second
		trail: 60, // Afterglow percentage
		shadow: false, // Whether to render a shadow
		hwaccel: true, // Whether to use hardware acceleration
		className: 'spinner', // The CSS class to assign to the spinner
		zIndex: 2e9, // The z-index (defaults to 2000000000)
		top: 'auto', // Top position relative to parent in px
		left: 'auto' // Left position relative to parent in px
};

// global variables to store the request and the spinner
// this is done so we can cancel the request and delete the spinner
// TODO move it into jquery data attribute?
var request;
var spinner;
var db_cursor;
