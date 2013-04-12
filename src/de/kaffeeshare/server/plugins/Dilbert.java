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

import java.net.URL;

import org.jsoup.nodes.Document;

import de.kaffeeshare.server.exception.PluginErrorException;

/**
 * Plugin to handle dilbert pages.
 */
public class Dilbert extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return (str.startsWith("https://feed.dilbert.com/")) || (str.startsWith("http://feed.dilbert.com/"))
				|| (str.startsWith("http://dilbert.com/strips/") || str.startsWith("https://dilbert.com/strips/"));
	}

	@Override
	public String getDescription(Document doc) {

		String description = "";
		try {
			String imageUrl = doc.getElementsByClass("STR_Image").first().children()
					.first().attr("src");
			if (imageUrl != null) {
				if (!imageUrl.contains("http://dilbert.com")) {
					imageUrl = "http://dilbert.com" + imageUrl;
				}
				description = "<img src=\"";
				description += imageUrl;
				description += "\" />";
			}
		} catch (Exception e) {
			throw new PluginErrorException(this);
		}

		return description;
	}

}
