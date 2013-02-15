package de.kaffeeshare.server.datastore;

import java.util.ResourceBundle;
import java.util.logging.Logger;

import de.kaffeeshare.server.datastore.appengine.AppEngineDatastore;
import de.kaffeeshare.server.datastore.jpa.JPADatastore;
import de.kaffeeshare.server.exception.DatastoreConfigException;
import de.kaffeeshare.server.utils.NamespaceValidator;

/**
 * Manager class to handle different datastore interfaces.
 */
public class DatastoreManager {
	
	private static final String datastoreConfig = ResourceBundle.getBundle("de.kaffeeshare.server.datastore.config").getString("datastore");
	
	private static final String APPENGINE = "AppEngine";
	private static final String JPA = "JPA";
	
	private static final Logger log = Logger.getLogger(DatastoreManager.class.getName());
	
	private static ThreadLocal<Datastore> datastore = new ThreadLocal<Datastore>();
	
	/**
	 * Get the used datastore.
	 * @return Datastore
	 */
	public static Datastore getDatastore() {
		
		if(datastoreConfig.equals(APPENGINE)) {
			if (datastore.get() == null) {
				log.info("Use AppEngine datastore interface.");
				datastore.set(new AppEngineDatastore());
			}
		} else if(datastoreConfig.equals(JPA)) {
			if (datastore.get() == null) {
				log.info("Use JPA datastore interface.");
				datastore.set(new JPADatastore());
			}
		} else {
			log.severe("No datastore interface defined. (See config.properties)");
			throw new DatastoreConfigException();
		}
		
		return datastore.get();
	}
	
	/**
	 * Sets namespace
	 * @param ns the namespace to be set
	 */
	public static void setNamespace(String ns) {
		
		if(NamespaceValidator.isValide(ns)) {
			log.info("Namespace set to " + ns);
			getDatastore().setNamespace(ns);
		}
	}

}
