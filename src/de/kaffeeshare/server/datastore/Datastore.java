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

import java.util.List;

/**
 * Datastore interface. A datatore backend have to implement this interface.
 */
public interface Datastore {
	
	/**
	 * Create a new item.
	 * @param catption Captionj
	 * @param url URL
	 * @param description Description
	 * @param imageUrl Image URL
	 * @return Item
	 */
	public Item createItem(String catption, String url, String description, String imageUrl);
	
	/**
	 * Stores a item and returns the persistent entity.
	 * @param item Item to store
	 * @return Persistent item
	 */
	public Item storeItem(Item item);
	
	/**
	 * Stores a list of items.
	 * @param items List with items to store
	 */
	public void storeItems(List<Item> items);
	
	/**
	 * Gets a list with items.
	 * @param maxNumber Number of list entries
	 * @return List with items
	 */
	public List<Item> getItems(int maxNumber);
	
	/**
	 * Sets the namespace.
	 * @param ns Namespace
	 */
	public void setNamespace(String ns);
	
	/**
	 * Check whether the current datastore with currently set namespace is empty.
	 * @return true if empty, otherwise false
	 */
	public boolean isEmpty();
}
