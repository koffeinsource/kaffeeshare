package de.kaffeeshare.server.datastore;

import java.util.MissingResourceException;
import java.util.ResourceBundle;

/**
 * Message class to handle strings.
 */
public class Messages {
	private static final String BUNDLE_NAME = "de.kaffeeshare.server.datastore.messages";

	private static final ResourceBundle RESOURCE_BUNDLE = ResourceBundle
			.getBundle(BUNDLE_NAME);

	private Messages() {
	}

	public static String getString(String key) {
		try {
			return RESOURCE_BUNDLE.getString(key);
		} catch (MissingResourceException e) {
			return '!' + key + '!';
		}
	}
}
