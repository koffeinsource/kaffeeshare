package de.kaffeeshare.server.datastore;

import java.net.MalformedURLException;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.Arrays;
import java.util.Iterator;
import java.util.List;
import java.util.logging.Logger;

/**
 * Namespace class to handle namespaces.
 */
public class UrlValidator {
	
	private static final Logger log = Logger.getLogger(UrlValidator.class.getName());

	/**
	 * A list of namespaces we want to reserver for special use
	 */
	private static List<String> validProtocols = Arrays.asList("http://", "https://");
	
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
	 * Encode a valid url string.
	 * @param url Url
	 * @return Valid url or null
	 */
	public static String encodeUrl(String url) {
		if(isValide(url)) {
			Iterator<String> iter = validProtocols.iterator();
			while(iter.hasNext()) {
				if(url.startsWith(iter.next())) {
					return url;
				}
			}
			
			return validProtocols.get(0) + url;
		}
		
		return null;
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
			log.info("Trying to use an illegal url: " + url);
			return false;
		}
		
		return true;
	}
}
