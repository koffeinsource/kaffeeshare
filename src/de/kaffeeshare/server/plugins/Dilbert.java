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

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URL;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.exception.InputErrorException;
import de.kaffeeshare.server.exception.SystemErrorException;

public class Dilbert extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (str.startsWith("https://feed.dilbert.com/")) || (str.startsWith("http://feed.dilbert.com/"))
				|| (str.startsWith("http://dilbert.com/strips/") || str.startsWith("https://dilbert.com/strips/"));
	}

	@Override
	public Item createItem(URL url) {
		log.info("Running Dilbert plugin!");

		Document doc;
		try {
			BufferedReader in = new BufferedReader(new InputStreamReader(
					url.openStream()));
			String html = "";
			String inputLine;
			while ((inputLine = in.readLine()) != null)
				html += (inputLine);
			in.close();
			doc = Jsoup.parse(html);
		} catch (IOException e1) {
			throw new SystemErrorException();
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

		String description = "";

		String imageUrl = null;
		try {
			imageUrl = doc.getElementsByClass("STR_Image").first().children()
					.first().attr("src");
			if (imageUrl != null) {
				if (!imageUrl.contains("http://dilbert.com")) {
					imageUrl = "http://dilbert.com" + imageUrl;
				}
				description = "<img src=\"";
				description += imageUrl;
				description += "\" />";
				imageUrl = "";
			}
		} catch (Exception e) {
			throw new InputErrorException();
		}
		log.info("description: " + description);

		String urlString = null;
		try {
			urlString = getProperty(doc, "og:url");
		} catch (Exception e) {
			urlString = url.toString();
		}

		return DatastoreManager.getDatastore().createItem(caption,urlString, description, imageUrl);
	}

}
