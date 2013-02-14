package de.kaffeeshare.server.plugins;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.InputErrorException;
import de.kaffeeshare.server.exception.SystemErrorException;

public class Dilbert extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (str.startsWith("https://feed.dilbert.com/")) || (str.startsWith("http://feed.dilbert.com/"))
				|| (str.startsWith("http://dilbert.com/strips/") || str.startsWith("https://dilbert.com/strips/"));
	}

	@Override
	public Item createItem(URL url) {
		log.info("Running Dilbert plugin!");

		Document doc;
		try {
			BufferedReader in = new BufferedReader(new InputStreamReader(
					url.openStream()));
			String html = "";
			String inputLine;
			while ((inputLine = in.readLine()) != null)
				html += (inputLine);
			in.close();
			doc = Jsoup.parse(html);
		} catch (IOException e1) {
			throw new SystemErrorException();
		}
		String caption = null;

		try {
			caption = getProperty(doc, "og:title");
		} catch (Exception e) {
		}

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
			imageUrl = doc.getElementsByClass("STR_Image").first().children()
					.first().attr("src");
			if (imageUrl != null) {
				if (!imageUrl.contains("http://dilbert.com")) {
					imageUrl = "http://dilbert.com" + imageUrl;
				}
				description = "<img src=\"";
				description += imageUrl;
				description += "\" />";
				imageUrl = "";
			}
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

		return DatastoreManager.getDatastore().createItem(caption,urlString, description, imageUrl);
	}

}
