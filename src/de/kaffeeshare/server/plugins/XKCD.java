package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.nodes.Document;

import de.kaffeeshare.server.exception.PluginErrorException;

public class XKCD extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (match(str, "xkcd.com") || match(str, "www.xkcd.com"));
	}
	
	@Override
	protected String getDescription(Document doc) {

		String description = null;

		try {
			description  = doc.getElementById("comic").html();
		} catch (Exception e) {
			throw new PluginErrorException(this);
		}
		
		return description;
	}

	@Override
	protected String getCaption(Document doc) {
		return "XKCD - " + doc.getElementById("ctitle").html();
	}
	
	@Override
	protected String getImageUrl(Document doc) {
		return null;
	}
}
