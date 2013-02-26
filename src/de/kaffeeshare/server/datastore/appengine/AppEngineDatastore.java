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
package de.kaffeeshare.server.datastore.appengine;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Date;
import java.util.List;
import java.util.logging.Logger;

import com.google.appengine.api.NamespaceManager;
import com.google.appengine.api.datastore.DatastoreService;
import com.google.appengine.api.datastore.DatastoreServiceFactory;
import com.google.appengine.api.datastore.Entities;
import com.google.appengine.api.datastore.Entity;
import com.google.appengine.api.datastore.FetchOptions;
import com.google.appengine.api.datastore.Key;
import com.google.appengine.api.datastore.KeyFactory;
import com.google.appengine.api.datastore.PreparedQuery;
import com.google.appengine.api.datastore.Query;
import com.google.appengine.api.datastore.Query.Filter;
import com.google.appengine.api.datastore.Query.FilterOperator;
import com.google.appengine.api.datastore.Query.FilterPredicate;
import com.google.appengine.api.datastore.Query.SortDirection;
import com.google.appengine.api.datastore.Text;

import de.kaffeeshare.server.datastore.Datastore;
import de.kaffeeshare.server.datastore.Item;

/**
 * Datastore helper class for googles app engine.
 */
public class AppEngineDatastore implements Datastore {

	private static final String DB_KIND_ITEM = "Item";
	private static final String DB_ITEM_CAPTION = "Caption";
	private static final String DB_ITEM_DESCRIPTION = "Description";
	private static final String DB_ITEM_CREATEDAT = "CreatedAt";
	private static final String DB_ITEM_IMAGEURL = "imageUrl";
	
	private static final Logger log = Logger.getLogger(AppEngineDatastore.class.getName());
	
	private static DatastoreService datastore = DatastoreServiceFactory.getDatastoreService();
	
	@Override
	public Item createItem(String caption, String url, String description, String imageUrl) {
		return new Item(caption, url, description, imageUrl);
	}
	
	@Override
	public Item storeItem(Item item) {
		Entity entity = toEntity(item);
		datastore.put(entity);
		return item;
	}

	@Override
	public void storeItems(List<Item> items) {
		List<Entity> entities = new ArrayList<Entity>();
		for (Item item : items) {
			entities.add(toEntity(item));
		}
		datastore.put(entities);
	}
	
	@Override
	public List<Item> getItems(int maxNumber) {
		Query query = new Query(DB_KIND_ITEM, null);
		query.addSort(DB_ITEM_CREATEDAT, SortDirection.DESCENDING);
		PreparedQuery pq = datastore.prepare(query);
		Collection<Entity> entities = pq.asList(FetchOptions.Builder.withLimit(maxNumber));
		return getItems(entities);
	}
	
	@Override
	public void setNamespace(String ns) {
		NamespaceManager.set(ns);
	}
	
	@Override
	public boolean isEmpty() {
		
		Query query = new Query(DB_KIND_ITEM, null);
		query.setKeysOnly();
		
		PreparedQuery pq = datastore.prepare(query);
		if (pq.asList(FetchOptions.Builder.withLimit(1)).size() > 0) return false;

		return true;
	}
	
	@Override
	public void garbageCollection(int maxKeepNumber, Date eldestDate) {
		
		// Metadata query to find all namespaces
		Query nsQuery = new Query(Entities.NAMESPACE_METADATA_KIND);
		nsQuery.setKeysOnly();
		
		for (Entity nsEntity : datastore.prepare(nsQuery).asIterable()) {

			String ns = nsEntity.getKey().getName();
			log.info("Garbage collection for ns: " + ns);
			setNamespace(ns);
						
			// First deletes everything but the first "maxKeepNumber" Items
			{
				Query query = new Query(DB_KIND_ITEM, null);
				query.addSort(DB_ITEM_CREATEDAT, SortDirection.DESCENDING);
				query.setKeysOnly();
	
				PreparedQuery pq = datastore.prepare(query);
				
				Iterable<Entity> entities = pq.asIterable(FetchOptions.Builder.withOffset(maxKeepNumber));
				
				deleteEntities(entities);
			}
			
			// Second delete all items elder than date
			{
				Filter filter = new FilterPredicate(DB_ITEM_CREATEDAT,
													FilterOperator.GREATER_THAN,
													eldestDate);
	
				Query query = new Query(DB_KIND_ITEM, null);
				query.setFilter(filter);
				query.setKeysOnly();
				
				PreparedQuery pq = datastore.prepare(query);
				
				Iterable<Entity> entities = pq.asIterable(FetchOptions.Builder.withDefaults());
				deleteEntities(entities);
			}
		}
	}
	
	/**
	 * Deletes a set of entities from the datastore
	 * @param entities the entities to be deleted
	 */
	private void deleteEntities(Iterable<Entity> entities) {
		List<Key> keys = new ArrayList<Key>();
		
		for( Entity e : entities) {
			keys.add(e.getKey());
		}
		
		log.info("Going to delete " + keys.size() + " enteties");
		
		datastore.delete(keys);
	}

	/**
	 * Gets a list with items from entities.
	 * @param entities Collection with entities
	 * @return List with items
	 */
	private List<Item> getItems(Collection<Entity> entities) {
		List<Item> items = new ArrayList<Item>();
		for (Entity entity : entities) {
			items.add(fromEntity(entity));
		}
		return items;
	}
	
	/**
	 * Creates an item from a google DB entity.
	 * @param Entity
	 * @return Item
	 */
	private Item fromEntity(Entity e) {
		return new Item((String) e.getProperty(DB_ITEM_CAPTION),
						(String) e.getKey().getName(),
						((Text) e.getProperty(DB_ITEM_DESCRIPTION)).getValue(),
						(String) e.getProperty(DB_ITEM_IMAGEURL),
						new Date((Long)e.getProperty(DB_ITEM_CREATEDAT))
						);
	}
	
	/**
	 * Returns a DB entity for the current item.
	 * @param item Item
	 * @return Entity
	 */
	private Entity toEntity(Item item) {
		Entity entity = new Entity(getItemDBKey(item));
		entity.setUnindexedProperty(DB_ITEM_CAPTION, item.getCaption());
		entity.setUnindexedProperty(DB_ITEM_DESCRIPTION, new Text(item.getDescription()));
		entity.setUnindexedProperty(DB_ITEM_IMAGEURL, item.getImageUrl());
		entity.setProperty(DB_ITEM_CREATEDAT, item.getCreatedAt().getTime());
		return entity;
	}
	
	/**
	 * The DB key of the current item.
	 * @param item Item
	 * @return Key
	 */
	private Key getItemDBKey(Item item) {
		return KeyFactory.createKey(DB_KIND_ITEM, item.getUrl());
	}
	
}
