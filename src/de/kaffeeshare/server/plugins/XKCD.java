package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.nodes.Document;

import de.kaffeeshare.server.exception.PluginErrorException;

public class XKCD extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (str.startsWith("http://xkcd.com") || str.startsWith("https://xkcd.com") ||
				str.startsWith("http://www.xkcd.com") || str.startsWith("https://www.xkcd.com")
				);
	}
	
	@Override
	public String getDescription(Document doc) {

		String description = null;

		try {
			description  = doc.getElementById("comic").html();
		} catch (Exception e) {
			throw new PluginErrorException(this);
		}
		
		return description;
	}

	@Override
	public String getCaption(Document doc) {
		return "XKCD - " + doc.getElementById("ctitle").html();
	}
	
	@Override
	public String getImageUrl(Document doc) {
		return "";
	}
}
