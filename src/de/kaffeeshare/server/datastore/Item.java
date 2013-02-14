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