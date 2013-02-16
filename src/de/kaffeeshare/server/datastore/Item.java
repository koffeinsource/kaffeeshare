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
package de.kaffeeshare.server.datastore;

import java.util.Date;

import javax.persistence.MappedSuperclass;

/**
 * A news item.
 */
@MappedSuperclass
public class Item {
	
	protected String caption = null;
	protected String url = null;
	protected String imageUrl = null;
	protected String description = null;
	protected Date createdAt = null;

	/**
	 * Creates a new item.
	 * This constructor is important for JPA
	 */
	public Item() {
	}
	
	/**
	 * Creates a new item.
	 */
	public Item(String caption, String url, String description, String imageUrl) {
		this.caption = caption;
		this.url = url;
		this.description = description;
		this.createdAt = new Date();
		this.imageUrl = imageUrl;
	}
	
	/**
	 * Creates a new item with defined creation date.
	 */
	public Item(String caption, String url, String description, String imageUrl, Date createdAt) {
		this(caption, url, description, imageUrl);
		this.createdAt = createdAt;
	}

	
	public String getCaption() {
		return caption;
	}

	public void setCaption(String caption) {
		this.caption = caption;
	}

	public String getUrl() {
		return url;
	}
	
	public boolean hasUrl() {
		return url != null;
	}

	public void setUrl(String url) {
		this.url = url;
	}

	public String getDescription() {
		return description;
	}

	public void setDescription(String description) {
		this.description = description;
	}

	public Date getCreatedAt() {
		return createdAt;
	}

	public String getImageUrl() {
		return imageUrl;
	}

	public void setImageUrl(String imageUrl) {
		this.imageUrl = imageUrl;
	}

}
