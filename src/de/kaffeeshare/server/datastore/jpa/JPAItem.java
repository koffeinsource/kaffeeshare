package de.kaffeeshare.server.datastore.jpa;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.NamedQueries;
import javax.persistence.NamedQuery;

import de.kaffeeshare.server.datastore.Item;

/**
 * A news item.
 */
@Entity
@NamedQueries ({
	@NamedQuery(name="findNsItems",query="SELECT i FROM JPAItem i WHERE i.namespace = :ns ORDER BY i.createdAt")
})
public class JPAItem extends Item {
	
	@Id
	@GeneratedValue(strategy = GenerationType.IDENTITY)
	private Long id = null;
	
	private String namespace = null;
	
	/**
	 * Creates a new item.
	 * This constructor is important for JPA
	 */
	public JPAItem() {
	}
	
	/**
	 * Creates a new item.
	 */
	public JPAItem(String caption, String url, String ns, String description, String imageUrl) {
		super(caption,url,description,imageUrl);
		this.namespace = ns;
	}

	public String getNamespace() {
		return namespace;
	}

	public void setNamespace(String ns) {
		this.namespace = ns;
	}

}