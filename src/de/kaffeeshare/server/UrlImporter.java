/*******************************************************************************
 * Copyright 2013 Jens Breitbart, Daniel Klemm, Dennis Obermann
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 ******************************************************************************/
package de.kaffeeshare.server;

import java.io.File;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.MissingResourceException;
import java.util.ResourceBundle;
import java.util.Vector;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.InputErrorException;
import de.kaffeeshare.server.exception.PluginErrorException;
import de.kaffeeshare.server.plugins.BasePlugin;
import de.kaffeeshare.server.plugins.DefaultPlugin;

/**
 * Imports a URL into our DB. The code calls plugins to extract the data from
 * the homepage.
 */
public class UrlImporter {
	
	private static final Logger log = Logger.getLogger(UrlImporter.class.getName());
	private static Vector<BasePlugin> plugins;
	
	/**
	 * Adds the first URL of a String to the DB.
	 * @param text Plain text to be parsed
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
	 * Adds the first URL of a String to the DB.
	 * @param text Html to be parsed
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
	 * @param urlString URL
	 * @return Item
	 * @throws InputErrorException
	 */
	public static Item fetchUrl(String urlString) throws InputErrorException {
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
	 * Lazy init of plugins. use this method to get the plugins. If you create a
	 * plugin - make sure to make it available here.
	 * @return Plugins
	 */
	private synchronized static Vector<BasePlugin> getPlugins() {

		if(plugins == null) {
			
			plugins = new Vector<BasePlugin>();
			
			// Get a file object for the plugins package
			File directory = new File(Thread.currentThread()
										.getContextClassLoader()
										.getResource("de/kaffeeshare/server/plugins").getFile());
			
			if(directory.exists()) {
				
				ResourceBundle bundle = ResourceBundle.getBundle("de.kaffeeshare.server.config");
				
				// Get the list of the files contained in the package
				String[] files = directory.list();
				for(int i = 0; i < files.length; i++) {
				
					// we are only interested in .class files
					if(files[i].endsWith(".class")) {
						
						// removes the .class extension
						String className = files[i].substring(0, files[i].length() - 6);
						
						// Check if the plugin is disabled in the config.properties
						boolean use = true;
						try {
							use = Boolean.valueOf(bundle.getString(className));
						} catch(MissingResourceException e) {
							// Plugin not defined in config.properties -> Use it
						}
						
						if(use && !(className.equals("BasePlugin") || className.equals("DefaultPlugin"))) {
							try {
								Object plugin = Class.forName("de.kaffeeshare.server.plugins." + className)
								                     .newInstance();
								
								if(plugin instanceof BasePlugin) {
									log.info("Start plugin: " +  className);
									plugins.add((BasePlugin)plugin);
								}

							} catch (Exception e) {
								log.warning("Can't start plugin " + className);
							}
						}
					}
				}
				
				
				
			}
		}
		
		return plugins;
	}

	/**
	 * Naive approach on finding the right plugin ;-)
	 * @param url URL
	 * @return Plugin 
	 */
	public static Item callMatchingPlugin(URL url) {
		
		try {
			for (BasePlugin plugin: getPlugins()) {
				if (plugin.match(url)) {
					return plugin.createItem(url);
				}
			}
		} catch(PluginErrorException e) {
			log.info(e.getMessage());
			// Plugin error, use the default plugin
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
	 * Extracts the first URL from a given string.
	 * @param text String
	 * @return URL
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

	/**
	 * URL pattern, html.
	 */
	static private Pattern urlPatternHTML = Pattern
			.compile("\\s*(?i)href\\s*=\\s*(\"([^\"]*\")|'[^']*'|([^'\">\\s]+))");

	/**
	 * Extracts the first URL from given HTML code.
	 * @param html Html string
	 * @return URL
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
