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
package de.kaffeeshare.server.input;

import java.io.IOException;
import java.util.logging.Logger;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import de.kaffeeshare.server.UrlImporter;
import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.exception.ReservedNamespaceException;

/**
 * Called by the browser extension-
 */
public class OneClickShare extends HttpServlet {

	private static final long serialVersionUID = 294584452111372279L;

	private static final Logger log = Logger.getLogger(OneClickShare.class.getName());

	private static String PARAM_URL = "url";
	private static String PARAM_NAMESPACE = "ns";
	
	/**
	 * imports url given by url parameter
	 * reponse is "0" on success and "1" on failure 
	 */
	public void doGet(HttpServletRequest req, HttpServletResponse resp)
			throws ServletException, IOException {
		String url = req.getParameter(PARAM_URL);
		resp.setContentType("text; charset=UTF-8");
		if (url != null) {
			try {
				DatastoreManager.setNamespace(req.getParameter(PARAM_NAMESPACE));
			} catch (ReservedNamespaceException e) {
				resp.getWriter().append("1"); // TODO return a different error code
				return;
			}
			try {
				log.info("Try to add url " + url + " to DB");
				DatastoreManager.getDatastore().storeItem(UrlImporter.fetchUrl(url));
				resp.getWriter().append("0");
			} catch (Exception e) {
				log.warning("exeption during import of url: " + e);
				resp.getWriter().append("1");
			}
		} else {
			log.warning("no url provided!");
			resp.getWriter().append("1");
		}
	}

}
