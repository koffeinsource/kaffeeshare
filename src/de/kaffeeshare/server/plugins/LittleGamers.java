package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;

import de.kaffeeshare.server.exception.PluginErrorException;

public class LittleGamers extends BasePlugin {

	@Override
	public boolean match(URL url) {
		return match(url, "www.little-gamers.com");
	}

	@Override
	protected String getDescription(Document doc) {

		String description = null;

		try {
			Elements elements = doc.getElementById("comic").getAllElements();
			description  = "<img src=\"";
			description += elements.get(0).attr("src");
			description += "\" />";
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
