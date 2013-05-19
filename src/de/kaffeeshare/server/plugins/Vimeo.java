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

/**
 * Plugin to handle vimeo pages.
 */
public class Vimeo extends BasePlugin {

	@Override
	public boolean match(URL url) {
		String str = url.toString();
		return match(str, "vimeo.com/");
	}

	@Override
	protected String getDescription(Document doc) {

		String videoId = null;
		try {
			videoId = getProperty(doc, "og:url").replace("http://vimeo.com/",
					"");
		} catch (Exception e) {}

		String description = super.getDescription(doc);
		if (videoId != null) {
			description += "<br /><br /><br /><iframe src=\"http://player.vimeo.com/video/"
					+ videoId
					+ "?title=0&amp;byline=0&amp;portrait=0\" width=\"400\" height=\"225\" frameborder=\"0\" webkitAllowFullScreen mozallowfullscreen allowFullScreen></iframe>";
		}

		return description;
	}
	
	@Override
	protected String getImageUrl(Document doc) {
		return null;
	}

}
