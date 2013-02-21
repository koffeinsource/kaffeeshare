package de.kaffeeshare.server.plugins;

import java.util.MissingResourceException;
import java.util.ResourceBundle;

/**
 * Message class to handle strings.
 */
public class Messages {
	private static final String BUNDLE_NAME = "de.kaffeeshare.server.plugins.messages"; //$NON-NLS-1$

	private static final ResourceBundle RESOURCE_BUNDLE = ResourceBundle
			.getBundle(BUNDLE_NAME);

	/**
	 * Void-constructor.
	 */
	private Messages() {
	}

	/**
	 * Get the string by a key.
	 * @param key Key
	 * @return String
	 */
	public static String getString(String key) {
		try {
			return RESOURCE_BUNDLE.getString(key);
		} catch (MissingResourceException e) {
			return '!' + key + '!';
		}
	}
}
