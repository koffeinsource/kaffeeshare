package de.kaffeeshare.server.plugins;

import java.io.IOException;
import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;

import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.InputErrorException;
import de.kaffeeshare.server.exception.SystemErrorException;

public class XKCD extends BasePlugin {

	@Override
	public boolean match(URL url) {
		return url.toString().contains("xkcd.com");
	}

	@Override
	public Item createItem(URL url) {
		log.info("Running xkcd plugin!");

		Document doc;
		try {
			doc = Jsoup.parse(url, 10000);
		} catch (IOException e1) {
			throw new SystemErrorException();
		}
		String caption = null;
		
		try {
			caption = getProperty(doc, "og:title");
		} catch (Exception e) {}
		
		if (caption == null) {
			try {
				caption = doc.select("title").first().text();
			} catch (Exception e) {
				caption = "";
			}
		}
		
		log.info("caption: " + caption);

		String description = "";

		String imageUrl = null;
		try {
			Elements elements = doc.select("#comic");
			log.info(elements.toString());
			description  = "<img src=\"";
			description += elements.get(0).child(0).attr("src");//elements.get(0).attr("src");
			log.info(description);
			description += "\" />";
		} catch (Exception e) {
			throw new InputErrorException();
		}
		log.info("description: " + description);
		
		String urlString = null;
		try {
			urlString = getProperty(doc, "og:url");
		} catch (Exception e) {
			urlString = url.toString();
		}
		
		return null;//new Item(caption, urlString, description, imageUrl);
	}

}
