package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;

import de.kaffeeshare.server.exception.PluginErrorException;

public class Gfycat extends BasePlugin {

	@Override
	public boolean match(URL url) {
		return match(url, "gfycat.com/");
	}

	@Override
	protected String getDescription(Document doc) {

		String description = null;

		try {
			Elements elements = doc.getElementsByClass("gfyVid");
			description = elements.toString();
		} catch (Exception e) {
			throw new PluginErrorException(this);
		}
		
		return description;
	}

	@Override
	protected String getImageUrl(Document doc) {
		return null;
	}
	
}
