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
import java.util.concurrent.atomic.AtomicBoolean;

import org.jsoup.nodes.Document;

import com.fasterxml.jackson.databind.ObjectMapper;

import de.kaffeeshare.server.exception.PluginErrorException;
import de.kaffeeshare.server.plugins.soundcloud.JSON_OEmbedInfo;

/**
 * Plugin to handle Soundcloud pages.
 */
public class Soundcloud extends BasePlugin {

	private String description = "";
	private String caption = "";
	private AtomicBoolean guard = new AtomicBoolean(false);
	private URL url;

	@Override
	public String getCaption(Document doc) {
		return caption;
	}

	@Override
	public String getDescription(Document doc) {
		String d = description;
		guard.set(false);
		return d;
	}

	@Override
	public boolean match(URL url) {
		
		if (match(url, "m.soundcloud.com") || match(url, "soundcloud.com")) {
			// wait until guard is false
			while (guard.getAndSet(true) == true) {}
			
			// remember url
			this.url = url;
			// reset vars
			description = "";
			caption = "";
			fetchOEmbedInfo();
			
			return true;
		}

		return false;
	}

	private void fetchOEmbedInfo() {
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
		} catch (Exception e) {
			log.warning("Could not receive OEmbed info: " + e + " "
					+ e.getMessage());
			guard.set(false);
			throw new PluginErrorException(this);
		}
	}

	@Override
	public String getImageUrl(Document doc) {
		return null;
	}

}