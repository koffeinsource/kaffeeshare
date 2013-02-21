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

	private static final String PARAM_URL = "url";
	private static final String PARAM_NAMESPACE = "ns";
	
	private static final String OK = "0";
	private static final String ERROR = "1";
	
	/**
	 * imports url given by url parameter
	 * reponse is OK on success and ERROR on failure 
	 */
	public void doGet(HttpServletRequest req, HttpServletResponse resp)
			throws ServletException, IOException {
		String url = req.getParameter(PARAM_URL);
		resp.setContentType("text; charset=UTF-8");
		if (url != null) {
			try {
				DatastoreManager.setNamespace(req.getParameter(PARAM_NAMESPACE));
			} catch (ReservedNamespaceException e) {
				resp.getWriter().append(ERROR); // TODO return a different error code
				return;
			}
			try {
				log.info(Messages.getString("OneClickShare.url_add") + url);
				DatastoreManager.getDatastore().storeItem(UrlImporter.fetchUrl(url));
				resp.getWriter().append(OK);
			} catch (Exception e) {
				log.warning(Messages.getString("OneClickShare.url_exception") + e);
				resp.getWriter().append(ERROR);
			}
		} else {
			log.warning(Messages.getString("OneClickShare.no_url"));
			resp.getWriter().append(ERROR); 
		}
	}

}