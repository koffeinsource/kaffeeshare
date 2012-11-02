package de.kaffeeshare.server.output;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.sun.syndication.feed.synd.SyndContent;
import com.sun.syndication.feed.synd.SyndContentImpl;
import com.sun.syndication.feed.synd.SyndEntry;
import com.sun.syndication.feed.synd.SyndEntryImpl;
import com.sun.syndication.feed.synd.SyndFeed;
import com.sun.syndication.feed.synd.SyndFeedImpl;
import com.sun.syndication.io.SyndFeedOutput;

import de.kaffeeshare.server.Config;
import de.kaffeeshare.server.datastore.Datastore;
import de.kaffeeshare.server.datastore.Item;
import de.kaffeeshare.server.datastore.Namespace;

/**
 * This servlet generates a rss feed with the latest news.
 */
public class Feed extends HttpServlet {

	private static final long serialVersionUID = -5819674729148390595L;

	private static String PARAM_NAMESPACE = "ns";
	
	private Logger log = Logger.getLogger(Feed.class.getName());
	
	/**
	 * Called when the feed url is requested.
	 */
	public void doGet(HttpServletRequest req, HttpServletResponse resp)
			throws ServletException, IOException {
		try {
			resp.setContentType("text; charset=UTF-8");
			String feedType = "rss_2.0";
			String ns = null;

			try {
				ns = req.getParameter(PARAM_NAMESPACE);
				Namespace.setNamespace(ns);
			} catch (Exception e) {
				// return an empty page if something isn't ok
				return;
			}

			SyndFeed feed = new SyndFeedImpl();
			feed.setFeedType(feedType);

			feed.setTitle(Config.Name + " - " + ns);
			feed.setLink("http://"+req.getServerName());
			feed.setDescription(Config.Phrase);

			List<SyndEntry> feedEntries = new ArrayList<SyndEntry>();
			List<Item> items = Datastore.getItems(20);
			for (Item item : items) {
				SyndEntry feedEntry;
				SyndContent feedContent;

				feedEntry = new SyndEntryImpl();
				feedEntry.setTitle(item.getCaption());
				feedEntry.setLink(item.getUrl());
				feedEntry.setPublishedDate(item.getCreatedAt());
				feedContent = new SyndContentImpl();
				feedContent.setType("html");
				String content = "";
				String imageUrl = item.getImageUrl();
				if (Item.isUrl(imageUrl)) {
					content += "<div style=\"float:left; margin-right:16px; margin-bottom:16px;\"><img width=\"200\" src=\""
							+ imageUrl + "\" alt=\"\"/></div>";
				}
				// escaped in validator - why?
				content += "<p>" + item.getDescription() + "</p>";
				if (item.hasUrl()) {
					content = content + " <a href=\"" + item.getUrl()
							+ "\">&raquo; " + item.getUrl() + "</a>";
				}
				feedContent.setValue(content);
				feedEntry.setDescription(feedContent);
				feedEntries.add(feedEntry);
			}

			feed.setEntries(feedEntries);

			SyndFeedOutput output = new SyndFeedOutput();
			output.output(feed, resp.getWriter());
			resp.getWriter().close();
		} catch (Exception ex) {
			log.log(Level.SEVERE, ex.getMessage());
		}
	}

}