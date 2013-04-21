package de.kaffeeshare.server.output;

import java.io.IOException;
import java.util.List;
import java.util.logging.Logger;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.google.appengine.labs.repackaged.org.json.JSONArray;
import com.google.appengine.labs.repackaged.org.json.JSONException;
import com.google.appengine.labs.repackaged.org.json.JSONObject;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;

public class JSON extends HttpServlet {
	private static final long serialVersionUID = 5705854772226283363L;

	private static final String PARAM_NAMESPACE = "ns";
	private static final String PARAM_OFFSET    = "start";
	
	// TODO magic number
	private static final int RESULTS_PER_REQUEST = 20;
	
	private static final Logger log = Logger.getLogger(JSON.class.getName());

	/**
	 * Handles a get request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOException
	 */
	public void doGet (HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		doPost(req, resp);
	}
	
	/**
	 * Handles a post request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOException
	 */
	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		String namespace = req.getParameter(PARAM_NAMESPACE);
		Integer offset = new Integer(req.getParameter(PARAM_OFFSET));
		offset *= RESULTS_PER_REQUEST;
		
		log.info("JSON request for namespace " + namespace + " offset " + offset);
		
		DatastoreManager.setNamespace(namespace);
		
		List<Item> items = DatastoreManager.getDatastore().getItems(RESULTS_PER_REQUEST, offset);
		JSONObject returnee = new JSONObject();
		JSONArray jarr = new JSONArray();
		
		for (Item i : items) {
			try {
				jarr.put(i.toJSON());
			} catch (JSONException e) {
				log.warning("JSON Exception @ JSON API");
				e.printStackTrace();
			}
		}

		try {
			returnee.put("items", jarr);
		} catch (JSONException e) {
			log.warning("JSON Exception @ JSON API");
			e.printStackTrace();
		}
		resp.setContentType("text; charset=UTF-8");
		
		resp.getWriter().append(returnee.toString());	
	}
}
