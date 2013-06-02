package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;

import de.kaffeeshare.server.exception.PluginErrorException;

public class Fefe extends BasePlugin {

	@Override
	public boolean match(URL url) {
		return match(url, "blog.fefe.de/?ts");
	}

	@Override
	protected String getDescription(Document doc) {

		String description = "";

		try {
			description = doc.getElementsByTag("li").html();
		} catch (Exception e) {
			throw new PluginErrorException(this);
		}
		
		description = description.substring(description.indexOf("</a>")+4, description.length());
				
		return description;
	}
	
	@Override
	protected String getCaption(Document doc) {
		return "Fefes Blog - " + Jsoup.parse(getDescription(doc)).text().replaceAll("(?<=.{20})\\b.*", "...");
	}
	
	@Override
	protected String getImageUrl(Document doc) {
		return null;
	}
}
