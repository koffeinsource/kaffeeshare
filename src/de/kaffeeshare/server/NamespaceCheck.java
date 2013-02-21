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
			log.info(Messages.getString("NamespaceCheck.check_ns") + namespace);
			
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
			log.warning(Messages.getString("NamespaceCheck.no_ns"));
			return;
		}
	}
	
	public void doGet (HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		doPost(req, resp);
	}
	
}
