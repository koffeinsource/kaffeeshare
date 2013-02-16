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
package de.kaffeeshare.server.plugins;

import java.io.IOException;
import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.SystemErrorException;

public class Vimeo extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (str.startsWith("http://vimeo.com/") || str.startsWith("https://vimeo.com/"));
	}

	@Override
	public Item createItem(URL url) {
		log.info("Running Vimeo plugin!");

		Document doc;
		try {
			doc = Jsoup.parse(url, 10000);
		} catch (IOException e1) {
			throw new SystemErrorException();
		}

		String videoId = null;
		try {
			videoId = getProperty(doc, "og:url").replace("http://vimeo.com/",
					"");
		} catch (Exception e) {
		}

		String caption = null;

		try {
			caption = getProperty(doc, "og:title");
		} catch (Exception e) {
		}

		if (caption == null) {
			try {
				caption = doc.select("title").first().text();
			} catch (Exception e) {
				caption = "";
			}
		}

		log.info("caption: " + caption);

		String description = null;

		try {
			description = getProperty(doc, "og:description");
		} catch (Exception e) {
		}

		if (description == null) {
			try {
				description = doc
						.getElementsByAttributeValue("name", "description")
						.first().attr("content");
			} catch (Exception e) {
				description = "";
			}
		}

		if (videoId != null) {
			description += "<br /><br /><br /><iframe src=\"http://player.vimeo.com/video/"
					+ videoId
					+ "?title=0&amp;byline=0&amp;portrait=0\" width=\"400\" height=\"225\" frameborder=\"0\" webkitAllowFullScreen mozallowfullscreen allowFullScreen></iframe>";
		}

		log.info("desc: " + description);

		String imageUrl = "";

		String urlString = null;
		try {
			urlString = getProperty(doc, "og:url");
		} catch (Exception e) {
			urlString = url.toString();
		}

		return DatastoreManager.getDatastore().createItem(caption,urlString, description, imageUrl);
	}

}
