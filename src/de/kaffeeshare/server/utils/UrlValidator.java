package de.kaffeeshare.server.utils;

import java.net.MalformedURLException;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.logging.Logger;

/**
 * Namespace class to handle namespaces.
 */
public class UrlValidator {
	
	private static final Logger log = Logger.getLogger(UrlValidator.class.getName());
	
	/**
	 * Validates a url. An exception if thrown in case it is not valid
	 * @param url the url to be validated
	 * @throws URISyntaxException 
	 * @throws MalformedURLException 
	 */
	private static void validate(String url) throws URISyntaxException, MalformedURLException {
		URL temp = new URL(url);
		temp.toURI();
	}
	
	/**
	 * Checks if a url is valid
	 * @param url the url to be validated
	 * @return true if url is valid, false otherwise
	 */
	public static boolean isValide(String url)  {
		try {
			validate(url);
		} catch (Exception e) {
			log.info("URL is not valid: " + url);
			return false;
		}
		
		return true;
	}
}
