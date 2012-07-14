// avoid conficts with other js frameworks
jQuery.noConflict();

window.addEventListener('keydown', function(event) {
	// Bind to both command (for Mac) and control (for Win/Linux)
	var modifier = event.ctrlKey || event.metaKey;
	// ctrl + .
	if (modifier && event.keyCode == 190) {
		// Send message to background page to toggle tab
		chrome.extension.sendRequest({
			share_page : true
		}, function(response) {
			// Do stuff on successful response
		});
	}
}, false);

chrome.extension.onRequest.addListener(function(request, sender, sendResponse) {
	var url = null;
	try {
		url = jQuery('#current-entry .entry-title-link').attr('href');
	} catch (e) {
	}
	sendResponse({
		urltoshare : url
	});
});