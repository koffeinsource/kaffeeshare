package de.kaffeeshare.server.plugins;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.ProtocolException;
import java.net.URL;

import de.kaffeeshare.server.datastore.Item;

public class Image extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String contentType = null;
		try {
			try {
				// lets try HTTP HEAD
				HttpURLConnection connection = (HttpURLConnection)  url.openConnection();
				connection.setRequestMethod("HEAD");
				connection.connect();
				contentType = connection.getContentType();
				log.info("HTTP GET for " + url.toString() + " successfull");
			} catch (ProtocolException e) {
				// ok, HTTP HEAD was not working, lets try a GET
				log.info("Fallback to HTTP GET for " + url.toString());
				contentType = url.openConnection().getContentType();
			}
		} catch (IOException e) {
			log.warning("Got IOException when calling getContent on " + url.toString());
			return false;
		}
		
		log.info("Found content type: " + contentType + " for " + url.toString());
		
		return contentType.startsWith("image/");
	}

	@Override
	public Item createItem(URL url) {
		log.info("Running Image plugin!");

		String caption = url.getFile().substring(1); // .getFile() returns "/<name>"
		String urlString = url.toString();
		String description = "<img src='" +url.toString()+ "'>";

		String imageUrl = "";
		return new Item(caption, urlString, description, imageUrl);
		}

}
