package de.kaffeeshare.server.output;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Logger;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.google.appengine.api.datastore.Cursor;
import com.google.appengine.labs.repackaged.org.json.JSONArray;
import com.google.appengine.labs.repackaged.org.json.JSONException;
import com.google.appengine.labs.repackaged.org.json.JSONObject;

import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.datastore.Item;

public class JSON extends HttpServlet {
	private static final long serialVersionUID = 5705854772226283363L;

	private static final String PARAM_NAMESPACE = "ns";
	private static final String PARAM_CURSOR    = "cursor";
	
	// TODO magic number
	private static final int RESULTS_PER_REQUEST = 20;
	
	private static final Logger log = Logger.getLogger(JSON.class.getName());

	/**
	 * Handles a post request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOException
	 */
	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		String namespace = req.getParameter(PARAM_NAMESPACE);
		String cursorStr = req.getParameter(PARAM_CURSOR);		
		
		log.info("JSON request for namespace " + namespace);
		
		DatastoreManager.setNamespace(namespace);
		
		Cursor cursor = null;
		
		if (cursorStr != null) {
			log.info("JSON request with cursor " + cursorStr);
			cursor = Cursor.fromWebSafeString(cursorStr);
		}
		
		List<Item> items = new ArrayList<Item>();
		cursor = DatastoreManager.getDatastore().getItems(RESULTS_PER_REQUEST, cursor, items);
		log.info("cur " + cursor.toWebSafeString());
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
			returnee.put(PARAM_CURSOR, cursor.toWebSafeString());
			
			log.info("cur " + cursor.toWebSafeString());
		} catch (JSONException e) {
			log.warning("JSON Exception @ JSON API");
			e.printStackTrace();
		}
		resp.setContentType("text; charset=UTF-8");
		
		resp.getWriter().append(returnee.toString());	
	}
	
	/**
	 * Handles a get request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, IOException
	 */
	public void doGet (HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		doPost(req, resp);
	}
}
