package de.kaffeeshare.server.plugins;

import de.kaffeeshare.server.UrlImporter;
import de.kaffeeshare.server.datastore.Item;

public class PluginSandbox {

	/**
	 * @param args
	 */
	public static void main(String[] args) {
		testUrl("http://dilbert.com/strips/comic/2012-04-05/?utm_source=feedburner&utm_medium=feed&utm_campaign=Feed%3A+dilbert%2Fdaily_strip+%28Dilbert+Daily+Strip+-+UU%29");
		testUrl("http://www.youtube.com/watch?v=C6O2ZIkcs-s&feature=g-vrec&context=G208b504RVAAAAAAAAAQ");
		testUrl("http://vimeo.com/39249572");
	}
	
	private static void testUrl(String urlString) {
		Item item = UrlImporter.fetchUrl(urlString);
		System.out.println("---------------------------------------");
		System.out.println("url: " + urlString);
		System.out.println("caption: " + item.getCaption());
		System.out.println("description: " + item.getDescription());
		System.out.println("image: " + item.getImageUrl());
	}

}
