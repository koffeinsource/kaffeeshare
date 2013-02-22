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
	 * @throws URISyntaxException, URISyntaxException 
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
		if (url == null) return false;
		
		try {
			validate(url);
		} catch (Exception e) {
			log.info("URL is not valid: " + url);
			return false;
		}
		
		return true;
	}
}
