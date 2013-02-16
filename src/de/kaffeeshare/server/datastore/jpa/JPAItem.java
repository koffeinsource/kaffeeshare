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
