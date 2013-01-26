package de.kaffeeshare.server.plugins;

import java.io.IOException;
import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;

import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.SystemErrorException;

/**
 * The default plugin. Should be used to parse generic
 * HTML webpages.
 */
public class DefaultPlugin extends BasePlugin {

	@Override
	public boolean match(URL url) {
		return true;
	}

	@Override
	public Item createItem(URL url) {
		Document doc;
		try {
			doc = Jsoup.parse(url, 10000);
		} catch (IOException e1) {
			log.warning(e1.getLocalizedMessage());
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

		String description = null;
		
		try {
			description = getProperty(doc, "og:description");
		} catch (Exception e) {
		}
		
		if (description == null) {
			try {
				description = doc.getElementsByAttributeValue("name", "description")
				                 .first().attr("content");
			} catch (Exception e) {
				description = "";
			}
		}
		
		log.info("desc: " + description);

		String imageUrl = "";
		try {
			imageUrl = getProperty(doc, "og:image");
			imageUrl.replace(" ", "");
			
			// check if URL is a legit URL as og:url is broken on some sites 
			// throws an exception if URL is not ok
			URL u = new URL(imageUrl);
			u.toURI(); 			
		} catch (Exception e) {
			imageUrl = "";
		}
		log.info("imageUrl: " + imageUrl);
		
		String urlString = null;
		try {
			urlString = getProperty(doc, "og:url");
			
			// check if URL is a legit URL as og:url is broken on some sites 
			// throws an exception if URL is not ok
			URL u = new URL(urlString);
			u.toURI(); 
		} catch (Exception e) {
			urlString = url.toString();
		}
		
		return new Item(caption, urlString, description, imageUrl);
	}

}
