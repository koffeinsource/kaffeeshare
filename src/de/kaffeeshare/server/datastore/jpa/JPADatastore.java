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

import java.util.List;

import javax.persistence.EntityManager;
import javax.persistence.EntityManagerFactory;
import javax.persistence.Persistence;
import javax.persistence.Query;

import de.kaffeeshare.server.datastore.Datastore;
import de.kaffeeshare.server.datastore.Item;

/**
 * Datastore helper class for googles app engine.
 */
public class JPADatastore implements Datastore {
	
	private EntityManagerFactory emfInstance =Persistence.createEntityManagerFactory("transactions-optional");
	
	private String namespace = null;
	
	/**
	 * Constructor.
	 */
	public JPADatastore() {
	}
	
	/**
	 * Creates an item.
	 * @param caption Caption
	 * @param url URL
	 * @param description Description
	 * @param imageUrl Image URL
	 * @return Item
	 */
	public Item createItem(String caption, String url, String description, String imageUrl) {
		return new JPAItem(caption, url, namespace, description, imageUrl);
	}
	
	/**
	 * Stores an item in the DB.
	 * @param item Item to store
	 * @return Stored item
	 */
	public Item storeItem(Item item) {
		
		Item persistentItem = null;
		
		EntityManager entityManager = emfInstance.createEntityManager();
		try {
			entityManager.getTransaction().begin();
			persistentItem = entityManager.merge(item);
			entityManager.getTransaction().commit();
		} catch(Exception e) {
			e.printStackTrace();
			entityManager.getTransaction().rollback();
			entityManager.close();
		} finally {
			entityManager.close();
		}

		return persistentItem;
	}

	/**
	 * Stores a list of items in the DB.
	 * @param items List with items
	 */
	public void storeItems(List<Item> items) {

		EntityManager entityManager = emfInstance.createEntityManager();
		try {
			entityManager.getTransaction().begin();
			for (Item item : items) {
				entityManager.merge(item);
			}
			entityManager.getTransaction().commit();
		} catch(Exception e) {
			e.printStackTrace();
			entityManager.getTransaction().rollback();
			entityManager.close();
		} finally {
			entityManager.close();
		}

	}
	
	@SuppressWarnings("unused")
	/**
	 * Deletes an item in DB, currently not used and only kept for reference.
	 * @param Item Item to delete
	 */
	private void deleteItem(Item item) {
		
		EntityManager entityManager = emfInstance.createEntityManager();
		try {
			entityManager.getTransaction().begin();
			entityManager.remove(item);
			entityManager.getTransaction().commit();
		} catch(Exception e) {
			e.printStackTrace();
			entityManager.getTransaction().rollback();
			entityManager.close();
		} finally {
			entityManager.close();
		}
	}

	/**
	 * Gets the newest items.
	 * @param ns Namespace
	 * @param maxNumber Number of items from DB
	 * @return List with items
	 */
	@SuppressWarnings("unchecked")
	public List<Item> getItems(int maxNumber) {

		List<Item> items = null;
		
		EntityManager entityManager = emfInstance.createEntityManager();
		try {

			Query q = entityManager.createNamedQuery("findNsItems");
			q.setParameter("ns", namespace);
			
			q.setMaxResults(maxNumber);
			
			items = q.getResultList();
		
		} catch(Exception e) {
			e.printStackTrace();
			entityManager.getTransaction().rollback();
			entityManager.close();
		} finally {
			entityManager.close();
		}
		
		return items;

	}

	/**
	 * Set the namespace.
	 * @param ns Namespace
	 */
	public void setNamespace(String ns) {
		this.namespace = ns;
	}

	/**
	 * Check if current namespace is unused.
	 * @return true, if namespace is unused
	 */
	public boolean isEmpty() {

		List<Item> items = getItems(1);
		if(items == null || items.isEmpty()) {
			return true;
		}

		return false;
	}

}
