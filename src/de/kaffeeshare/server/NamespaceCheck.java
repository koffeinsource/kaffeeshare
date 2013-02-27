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

import com.google.appengine.labs.repackaged.org.json.JSONException;
import com.google.appengine.labs.repackaged.org.json.JSONObject;

import de.kaffeeshare.server.datastore.DatastoreManager;

/**
 * Namespace check service.
 */
public class NamespaceCheck extends HttpServlet {

	private static final long serialVersionUID = -4417846428551748172L;
	private static final String PARAM_NAMESPACE = "ns";
	
	private static final Logger log = Logger.getLogger(NamespaceCheck.class.getName());

	/**
	 * Handle a post request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOException
	 */
	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		String namespace = req.getParameter(PARAM_NAMESPACE);
		resp.setContentType("text; charset=UTF-8");
		
		try {
			if (namespace != null) {
				log.info("Check status of namespace: " + namespace);
				
				JSONObject json = new JSONObject();
				
				// By convention, all namespaces starting with "_" (underscore) are reserved for system use.
				if (namespace.charAt(0) == '_') {
					json.put("status", "error");
					resp.getWriter().append(json.toString());
					return;
				}
				
				// check if namespace is valid 
				try {
					DatastoreManager.setNamespace(namespace);
				} catch (Exception e) {
					json.put("status", "error");
					resp.getWriter().append(json.toString());
					return;
				}
				
				// check if email and jabber address are valid!
				// I believe a valid email address also a valid jabber address
				// so I only check for valid email address
				try {
					InternetAddress emailAddr = new InternetAddress(namespace+"@abc.com");
					emailAddr.validate();
				} catch (AddressException ex) {
					json.put("status", "error");
					resp.getWriter().append(json.toString());
					return;
				}
				
				// check if namespace is empty
				if (!DatastoreManager.getDatastore().isEmpty()) {
					json.put("status", "use");
					resp.getWriter().append(json.toString());
					return;
				}
				
				json.put("status", "success");
				resp.getWriter().append(json.toString());
			} else {
				log.warning("no namespace provided!");
				return;
			}
		} catch (JSONException e) {
			log.warning("JSON error!");
			e.printStackTrace();
		}
	}
	
	/**
	 * Handle a get request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOException
	 */
	public void doGet (HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		doPost(req, resp);
	}
	
}
