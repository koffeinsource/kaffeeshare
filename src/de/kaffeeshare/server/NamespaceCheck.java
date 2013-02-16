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
package de.kaffeeshare.server;

import java.io.IOException;
import java.util.logging.Logger;

import javax.mail.internet.AddressException;
import javax.mail.internet.InternetAddress;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import de.kaffeeshare.server.datastore.DatastoreManager;

public class NamespaceCheck extends HttpServlet {

	private static final long serialVersionUID = -4417846428551748172L;
	private static String PARAM_NAMESPACE = "ns";
	
	private static final Logger log = Logger.getLogger(NamespaceCheck.class.getName());

	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		String namespace = req.getParameter(PARAM_NAMESPACE);
		resp.setContentType("text; charset=UTF-8");
		
		if (namespace != null) {
			log.info("Check status of namespace: " + namespace);
			
			// By convention, all namespaces starting with "_" (underscore) are reserved for system use.
			if (namespace.charAt(0) == '_') {
				resp.getWriter().append("{\"status\": \"error\"}");
				return;
			}
			
			// check if namespace is valid 
			try {
				DatastoreManager.setNamespace(namespace);
			} catch (Exception e) {
				resp.getWriter().append("{\"status\": \"error\"}");
				return;
			}
			
			// check if email and jabber address are valid!
			// I think a valid email address also a valid jabber address
			// so I only check for valid email address
			try {
				InternetAddress emailAddr = new InternetAddress(namespace+"@abc.com");
				emailAddr.validate();
			} catch (AddressException ex) {
				resp.getWriter().append("{\"status\": \"error\"}");
				return;
			}
			
			// check if namespace is empty
			if (!DatastoreManager.getDatastore().isEmpty()) {
				resp.getWriter().append("{\"status\": \"use\"}");
				return;
			}
			
			resp.getWriter().append("{\"status\": \"success\"}");
		} else {
			log.warning("no namespace provided!");
			return;
		}
	}
	
	public void doGet (HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		doPost(req, resp);
	}
	
}
