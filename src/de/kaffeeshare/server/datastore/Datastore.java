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
