package de.kaffeeshare.server.plugins;

import java.net.URL;
import java.util.logging.Logger;

import org.jsoup.nodes.Document;

import de.kaffeeshare.server.datastore.Item;

/**
 * Abstract base class for all plugins
 */
public abstract class BasePlugin {

	protected static final Logger log = Logger.getLogger(BasePlugin.class.getName());
	
	public abstract boolean match(URL url);
	public abstract Item createItem(URL url);

	protected String getProperty(Document doc, String prop) {
		String caption;
		caption = doc.getElementsByAttributeValue("property", prop)
		             .first().attr("content");
		return caption;
	}	
}
