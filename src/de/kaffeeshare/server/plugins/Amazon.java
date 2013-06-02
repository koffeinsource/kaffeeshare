package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.nodes.Document;

import de.kaffeeshare.server.Config;

/**
 * Plugin to handle amazon pages.
 */
public class Amazon extends BasePlugin {

	// Add here the url extension string to earn money ;-)
	private final static String urlExtension =  "tag=" + Config.getString("amazon_url_extension");

	@Override
	public boolean match(URL url) {
		return match(url, "amazon");
	}

	@Override
	protected String getURL(URL url) {

		String urlStr = url.toString();
		int index = urlStr.indexOf("tag=");
		if(index != -1) {
			int endIndex = urlStr.indexOf("&", index + 1);
			String replace = "";
			if(endIndex != -1) {
				replace = urlStr.substring(index, endIndex);
			} else {
				replace = urlStr.substring(index);
			}
			return urlStr.replace(replace, urlExtension);
		}

		if(!urlStr.contains("?")) {
			urlStr += "?";
		} else {
			urlStr += "&";
		}
		return urlStr + urlExtension;
	}

	@Override
	protected String getImageUrl(Document doc) {

		String imageUrl = "";
		try {
			imageUrl = doc.getElementById("main-image").attr("src");
		} catch(Exception e) {
			// No image available
		}
		return imageUrl;

	}
}
