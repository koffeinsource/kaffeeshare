package de.kaffeeshare.server;

import java.io.IOException;
import java.util.logging.Logger;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import de.kaffeeshare.server.datastore.Datastore;
import de.kaffeeshare.server.datastore.Namespace;

public class NamespaceCheck extends HttpServlet {

	private static final long serialVersionUID = -4417846428551748172L;
	private static String PARAM_NAMESPACE = "ns";
	
	private static final Logger log = Logger.getLogger(NamespaceCheck.class.getName());

	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		String namespace = req.getParameter(PARAM_NAMESPACE);
		resp.setContentType("text; charset=UTF-8");
				
		if (namespace != null) {
			log.info("Check status of namespace: " + namespace);
			// check if namespace is valid 
			try {
				Namespace.setNamespace(namespace);
			} catch (Exception e) {
				resp.getWriter().append("{\"status\": \"error\"}");
				return;
			}
			
			// TODO check if email and jabber address are valid!
			
			// check if namespace is empty
			if (Datastore.getItems(1).size() > 0) {
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
