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
import org.jsoup.select.Elements;

import de.kaffeeshare.server.exception.PluginErrorException;

/**
 * Plugin to handle garfield pages.
 */
public class Garfield extends BasePlugin {

	@Override
	public boolean match(URL url) {
		return match(url, "www.gocomics.com/garfield");
	}

	@Override
	protected String getDescription(Document doc) {

		String description = null;

		try {
			Elements elements = doc.getElementsByClass("strip");
			description  = "<img src=\"";
			description += elements.get(0).attr("src");
			description += "\" />";
		} catch (Exception e) {
			throw new PluginErrorException(this);
		}
		
		return description;
	}

	@Override
	protected String getImageUrl(Document doc) {
		return null;
	}

}
