package de.kaffeeshare.server.plugins.soundcloud;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class JSON_OEmbedInfo {

	@JsonProperty("title")
	private String title;
	@JsonProperty("description")
	private String description;
	@JsonProperty("html")
	private String html;

	/**
	 * Required by Jackson
	 */
	@JsonCreator
	public JSON_OEmbedInfo() {
	}

	public String getDescription() {
		return description;
	}

	public String getHtml() {
		return html;
	}

	public String getTitle() {
		return title;
	}

	public void setDescription(String description) {
		this.description = description;
	}

	public void setHtml(String html) {
		this.html = html;
	}

	public void setTitle(String title) {
		this.title = title;
	}

}
