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
import de.kaffeeshare.server.exception.SystemErrorException;

public class JSON extends HttpServlet {
	private static final long serialVersionUID = 5705854772226283363L;

	private static final String PARAM_NAMESPACE  = "ns";
	private static final String PARAM_CURSOR     = "cursor";
	private static final String PARAM_OP         = "op";
	private static final String PARAM_OP_GET     = "get";
	private static final String PARAM_OP_UPDATED = "updated";
	
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
		JSONObject returnee = new JSONObject();
		
		String op = req.getParameter(PARAM_OP);
		if (op.equals(PARAM_OP_GET)) {
			getItems(req, returnee);
		}
		if (op.equals(PARAM_OP_UPDATED)) {
			String namespace = req.getParameter(PARAM_NAMESPACE);
			log.info("JSON " + PARAM_OP_UPDATED + " request for namespace " + namespace);
			
			DatastoreManager.setNamespace(namespace);
			List<Item> items = new ArrayList<Item>();
			DatastoreManager.getDatastore().getItems(1, items);
			try {
				if (items.size()==0) returnee.put("last_update", 0);
				else returnee.put("last_update", items.get(0).getCreatedAt().getTime());
			} catch (JSONException e) {
				log.warning("JSON Exception @ JSON API");
				e.printStackTrace();
				throw new SystemErrorException();
			}
	
		}
		resp.setContentType("text; charset=UTF-8");
		
		resp.getWriter().append(returnee.toString());	
	}

	private void getItems(HttpServletRequest req, JSONObject returnee) {
		String namespace = req.getParameter(PARAM_NAMESPACE);
		String cursorStr = req.getParameter(PARAM_CURSOR);		
		
		log.info("JSON " + PARAM_OP_GET + " request for namespace " + namespace);
		
		DatastoreManager.setNamespace(namespace);
		
		Cursor cursor = null;
		
		if (cursorStr != null) {
			log.info("JSON request with cursor " + cursorStr);
			cursor = Cursor.fromWebSafeString(cursorStr);
		}
		
		List<Item> items = new ArrayList<Item>();
		cursor = DatastoreManager.getDatastore().getItems(RESULTS_PER_REQUEST, cursor, items);
		log.info("cur " + cursor.toWebSafeString());
		JSONArray jarr = new JSONArray();
				
		for (Item i : items) {
			try {
				jarr.put(i.toJSON());
			} catch (JSONException e) {
				log.warning("JSON Exception @ JSON API");
				e.printStackTrace();
				throw new SystemErrorException();
			}
		}

		try {
			returnee.put("items", jarr);
			returnee.put(PARAM_CURSOR, cursor.toWebSafeString());
			
			log.info("cur " + cursor.toWebSafeString());
		} catch (JSONException e) {
			log.warning("JSON Exception @ JSON API");
			e.printStackTrace();
			throw new SystemErrorException();
		}
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
