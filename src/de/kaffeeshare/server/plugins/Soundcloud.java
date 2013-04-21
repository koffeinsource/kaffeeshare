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

import com.fasterxml.jackson.databind.ObjectMapper;

import de.kaffeeshare.server.plugins.soundcloud.JSON_OEmbedInfo;

/**
 * Plugin to handle Soundcloud pages. XXX NOT THREADSAFE Has to use variables to
 * store data until getCaption() and getDescription() are called
 */
public class Soundcloud extends BasePlugin {

	private String description = "";
	private String caption = "";
	private boolean oEmbedInfoFetched = false;
	private URL url;

	@Override
	public String getCaption(Document doc) {
		maybeFetchOEmbedInfo();
		return caption;
	}

	@Override
	public String getDescription(Document doc) {
		maybeFetchOEmbedInfo();
		return description;
	}

	@Override
	public boolean match(URL url) {
		// remember url
		this.url = url;
		// reset vars
		oEmbedInfoFetched = false;
		description = "";
		caption = "";

		String str = url.toString();
		return (str.startsWith("http://m.soundcloud.com/")
				|| str.startsWith("https://m.soundcloud.com/")
				|| str.startsWith("https://soundcloud.com/")
				|| str.startsWith("https://www.soundcloud.com/")
				|| str.startsWith("http://soundcloud.com/") || str
					.startsWith("http://www.soundcloud.com/"));
	}

	private void maybeFetchOEmbedInfo() {
		// only once
		if (!oEmbedInfoFetched) {
			try {
				String oEmbedUrl = "http://soundcloud.com/oembed?format=json&url="
						+ url.toExternalForm();
				log.info("fetching oebemd info: " + oEmbedUrl);
				ObjectMapper mapper = new ObjectMapper();
				JSON_OEmbedInfo json_OEmbedInfo = mapper.readValue(new URL(
						oEmbedUrl), JSON_OEmbedInfo.class);
				log.info("title: " + json_OEmbedInfo.getTitle());
				log.info("description: " + json_OEmbedInfo.getDescription());
				log.info("html: " + json_OEmbedInfo.getHtml());
				caption = json_OEmbedInfo.getTitle();
				description = json_OEmbedInfo.getHtml()
						+ "<br /><br /><br />"
						+ json_OEmbedInfo.getDescription();
				oEmbedInfoFetched = true;
			} catch (Exception e) {
				log.warning("Could not receive OEmbed info: " + e + " "
						+ e.getMessage());
			}
		}
	}

	@Override
	public String getImageUrl(Document doc) {
		return null;
	}

}
