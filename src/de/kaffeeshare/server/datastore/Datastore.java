package de.kaffeeshare.server.datastore;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

import com.google.appengine.api.datastore.DatastoreService;
import com.google.appengine.api.datastore.DatastoreServiceFactory;
import com.google.appengine.api.datastore.Entity;
import com.google.appengine.api.datastore.FetchOptions;
import com.google.appengine.api.datastore.Key;
import com.google.appengine.api.datastore.PreparedQuery;

/**
 * Datastore helper class.
 */
public class Datastore {

	private static DatastoreService datastore = DatastoreServiceFactory.getDatastoreService();

	/**
	 * Stores an item in the DB.
	 */
	public static Item storeItem(Item item) {
		Entity entity = item.toEntity();
		datastore.put(entity);
		return item;
	}

	/**
	 * Stores a list of items in the DB.
	 */
	public static void storeItems(List<Item> items) {
		List<Entity> entities = new ArrayList<Entity>();
		for (Item item : items) {
			entities.add(item.toEntity());
		}
		datastore.put(entities);
	}
	
	/**
	 * Check if current namespace is unused
	 */
	public static boolean isEmpty() {	
		PreparedQuery pq = datastore.prepare(Item.isEmpty());
		if (pq.asList(FetchOptions.Builder.withLimit(1)).size() > 0) return false;

		return true;
	}
	
	@SuppressWarnings("unused")
	/**
	 * Deletes an item in DB, currently not used and only kept for reference.
	 */
	private static void deleteItem(Item item) {
		Key key = item.getDBKey();
		datastore.delete(key);
	}

	/**
	 * Gets the newest @param maxNumber item from DB.
	 */
	public static List<Item> getItems(int maxNumber) {
		PreparedQuery pq = datastore.prepare(Item.getDBQuery());
		Collection<Entity> entities = pq.asList(FetchOptions.Builder.withLimit(maxNumber));
		return getItems(entities);
	}

	private static List<Item> getItems(Collection<Entity> entities) {
		List<Item> items = new ArrayList<Item>();
		for (Entity entity : entities) {
			items.add(new Item(entity));
		}
		return items;
	}
}