package de.kaffeeshare.server.plugins;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.ProtocolException;
import java.net.URL;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;

public class Image extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String contentType = null;
		try {
			try {
				// lets try HTTP HEAD
				HttpURLConnection connection = (HttpURLConnection) url.openConnection();
				connection.setRequestMethod("HEAD");
				connection.connect();
				contentType = connection.getContentType();
				log.info(Messages.getString("Image.http_get_successfull") + url.toString());
			} catch (ProtocolException e) {
				// ok, HTTP HEAD was not working, lets try a GET
				log.info(Messages.getString("Image.http_get_fallback") + url.toString());
				contentType = url.openConnection().getContentType();
			}
		} catch (IOException e) {
			log.warning(Messages.getString("Image.io_exception") + url.toString());
			return false;
		}
		
		if (contentType == null) return false;
		
		log.info(Messages.getString("Image.content_type_found") + contentType + ": " + url.toString());
		
		return contentType.startsWith("image/");
	}

	@Override
	public Item createItem(URL url) {
		log.info(Messages.getString("Image.running"));

		String caption = url.getFile().substring(1); // .getFile() returns "/<name>"
		String urlString = url.toString();
		String description = "<img src='" +url.toString()+ "'>";

		String imageUrl = "";
		return DatastoreManager.getDatastore().createItem(caption,urlString, description, imageUrl);
	}

}
