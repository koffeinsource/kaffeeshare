/*******************************************************************************
 * Copyright 2013 Jens Breitbart
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
