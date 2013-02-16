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
		NamespaceValidator.validate(ns);
		
		log.info("Namespace set to " + ns);
		getDatastore().setNamespace(ns);
	}

}
