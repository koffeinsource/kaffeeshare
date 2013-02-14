package de.kaffeeshare.server;

import java.net.MalformedURLException;
import java.net.URL;
import java.util.Vector;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.InputErrorException;
import de.kaffeeshare.server.plugins.BasePlugin;
import de.kaffeeshare.server.plugins.DefaultPlugin;
import de.kaffeeshare.server.plugins.Dilbert;
import de.kaffeeshare.server.plugins.Garfield;
import de.kaffeeshare.server.plugins.Image;
import de.kaffeeshare.server.plugins.Vimeo;
import de.kaffeeshare.server.plugins.Youtube;

/**
 * Imports a URL into our DB. The code calls plugins to extract the data from
 * the homepage.
 */
public class UrlImporter {

	private static final Logger log = Logger.getLogger(UrlImporter.class.getName());
	private static Vector<BasePlugin> plugins;
	
	/**
	 * Adds the first URL of a String to the DB
	 * 
	 * @param text
	 *            plain text to be parsed
	 * @return null if no URL is in the text, otherwise the url
	 */
	public static String importFromText(String text) {
		String url = getURLPlain(text);

		if (url != null) {
			log.info("Try to add url " + url + " to DB");
			DatastoreManager.getDatastore().storeItem(fetchUrl(url));
		}

		return url;
	}

	/**
	 * Adds the first URL of a String to the DB
	 * 
	 * @param text
	 *            html to be parsed
	 * @return null if no URL is in the text, otherwise the url
	 */
	public static String importFromHTML(String text) {
		String url = getURLHTML(text);

		if (url != null) {
			log.info("Try to add url " + url + " to DB");
			DatastoreManager.getDatastore().storeItem(fetchUrl(url));
		}

		return url;
	}

	/**
	 * Generates an item from a url.
	 */
	public static Item fetchUrl(String urlString) {
		try {
			if (urlString.startsWith("http://") || urlString.startsWith("https://")) {
			} else {
				urlString = "http://" + urlString;
			}
			URL url = new URL(urlString);
			return callMatchingPlugin(url);
		} catch (MalformedURLException e) {
			throw new InputErrorException();
		}
	}

	/**
	 * lazy init of plugins. use this method to get the plugins. if you create a
	 * plugin - make sure to make it available here
	 * 
	 * @return plugins
	 */
	private synchronized static Vector<BasePlugin> getPlugins() {
		if (plugins == null) {
			plugins = new Vector<BasePlugin>();
			plugins.add(new Image());
			plugins.add(new Garfield());
			plugins.add(new Youtube());
			plugins.add(new Vimeo());
			plugins.add(new Dilbert());
		}
		return plugins;
	}

	/**
	 * naive approach on finding the right plugin ;-)
	 * @param url
	 * @return plugin 
	 */
	public static Item callMatchingPlugin(URL url) {
		for (BasePlugin plugin: getPlugins()) {
			if (plugin.match(url)) {
				return plugin.createItem(url);
			}
		}
		return new DefaultPlugin().createItem(url);
	}

	/**
	 * URL pattern, public domain.
	 */
	static private Pattern urlPatternPlain = Pattern
			.compile("\\b((http(s?)\\:\\/\\/|~\\/|\\/)|www.)"
					+ "(\\w+:\\w+@)?(([-\\w]+\\.)+(com|org|net|gov"
					+ "|mil|biz|info|mobi|name|aero|jobs|museum"
					+ "|travel|[a-z]{2}))(:[\\d]{1,5})?"
					+ "(((\\/([-\\w~!$+|.,=]|%[a-f\\d]{2})+)+|\\/)+|\\?|#)?"
					+ "((\\?([-\\w~!$+|.,*:]|%[a-f\\d{2}])+=?"
					+ "([-\\w~!$+|.,*:=]|%[a-f\\d]{2})*)"
					+ "(&(?:[-\\w~!$+|.,*:]|%[a-f\\d{2}])+=?"
					+ "([-\\w~!$+|.,*:=]|%[a-f\\d]{2})*)*)*"
					+ "(#([-\\w~!$+|.,*:=]|%[a-f\\d]{2})*)?\\b");

	/**
	 * Extracts the first URL from a given Strings
	 */
	static private String getURLPlain(String text) {
		Matcher m = urlPatternPlain.matcher(text);

		while (m.find()) {
			String url = m.group();
			if (url.startsWith("(") && url.endsWith(")")) {
				url = url.substring(1, url.length() - 1);
			}

			log.info("found url: " + url);
			return url;
		}
		return null;

	}

	static private Pattern urlPatternHTML = Pattern
			.compile("\\s*(?i)href\\s*=\\s*(\"([^\"]*\")|'[^']*'|([^'\">\\s]+))");

	/**
	 * Extracts the first URL from given HTML code
	 */
	static private String getURLHTML(String html) {
		Matcher m = urlPatternHTML.matcher(html);

		while (m.find()) {
			String url = m.group();

			log.info("found url: " + url);
			return url;
		}

		return null;
	}

}
