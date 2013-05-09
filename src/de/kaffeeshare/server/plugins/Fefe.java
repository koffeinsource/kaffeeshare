package de.kaffeeshare.server.plugins;

import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;

import de.kaffeeshare.server.exception.PluginErrorException;

public class Fefe extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (str.startsWith("http://blog.fefe.de/?ts") || str.startsWith("https://blog.fefe.de/?ts"));
	}

	@Override
	public String getDescription(Document doc) {

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
	public String getCaption(Document doc) {
		return "Fefes Blog - " + Jsoup.parse(getDescription(doc)).text().replaceAll("(?<=.{20})\\b.*", "...");
	}
	
	@Override
	public String getImageUrl(Document doc) {
		return "";
	}
}
