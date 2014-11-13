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
package de.kaffeeshare.server.cron;

import java.io.IOException;
import java.text.DateFormat;
import java.util.Calendar;
import java.util.ResourceBundle;
import java.util.logging.Logger;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import de.kaffeeshare.server.datastore.DatastoreManager;

/**
 * Service to start garbage collection on the db.
 */
public class GarbageCollection extends HttpServlet {

	private static final long serialVersionUID = 8864732716177174197L;

	private static final Logger log = Logger.getLogger(GarbageCollection.class.getName());

	private static final int maxKeepNumber = Integer.valueOf(ResourceBundle.getBundle("de.kaffeeshare.server.config").getString("max_ns_items"));
	private static final int lifeTime = Integer.valueOf(ResourceBundle.getBundle("de.kaffeeshare.server.config").getString("item_life_time"));

	/**
	 * Handle a get request.
	 * Remove old data from the db.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOExcption
	 */
	public void doGet(HttpServletRequest req, HttpServletResponse resp)
			throws ServletException, IOException {
		
		// Calculate date: Current date - liftTime
		Calendar deleteDate = Calendar.getInstance();
		deleteDate.add(Calendar.DAY_OF_YEAR, -1*lifeTime);
		
		log.info("Start garbage collection. Max ns items: " + String.valueOf(maxKeepNumber) + 
				 ", Life time: " + String.valueOf(lifeTime) + " -> Deletion date: " + DateFormat.getDateInstance().format(deleteDate.getTime()));
		DatastoreManager.getDatastore().garbageCollection(maxKeepNumber, deleteDate.getTime());

	}

}
