/*******************************************************************************
 * Copyright 2013 Jens Breitbart
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
 * Datastore interface.
 */
public interface Datastore {
	
	public Item createItem(String catption, String url, String description, String imageUrl);
	
	public Item storeItem(Item item);
	
	public void storeItems(List<Item> items);
	
	public List<Item> getItems(int maxNumber);
	
	public void setNamespace(String ns);
	
	public boolean isEmpty();
}
