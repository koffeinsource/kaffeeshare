package de.kaffeeshare.server;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 * Servlet to generate the startpage.
 */
public class Page extends HttpServlet {

	private static final long serialVersionUID = -1666980819128719605L;

	public void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
		resp.setContentType("text/html; charset=UTF-8");
		resp.getWriter().append(printPage());
		resp.getWriter().flush();
	}

	private String printPage() {
		StringBuffer out = new StringBuffer();
		out.append("<html>");
		out.append("<head>");
		out.append("<title>" + Config.Name + "</title>");
		out.append("<meta http-equiv=\"content-type\" content=\"text/html; charset=UTF-8\" />");
		out.append("<meta name=\"keywords\" content=\"" + Config.Name + "\" />");
		out.append("<meta name=\"description\" content=\"" + Config.Name
				+ " | " + Config.Phrase + "\" />");
		out.append("<meta name=\"viewport\" content=\"width=400; initial-scale=1.0; maximum-scale=2.0; user-scalable=1;\">");
		out.append("<link type=\"text/css\" rel=\"stylesheet\" href=\"/style.css\" />");
		out.append("<link rel=\"alternate\" type=\"application/rss+xml\" title=\"RSS Feed\" href=\"/feed\">");
		out.append("</head>");
		out.append("<body>");
		out.append("<div id=\"allWrap\">");
		out.append("<div id=\"view\">");
		printHeader(out);
		out.append("</div>");
		out.append("</div>");
		out.append("</body>");
		out.append("</html>");
		return out.toString();
	}

	private void printHeader(StringBuffer out) {
		out.append("<div class=\"widget\">");
		out.append("<div class=\"siteTitle\">");
		out.append("<a class=\"siteTitleText\" href=\"/\">" + Config.Name
				+ "</a>");
		out.append("</div>");// title
		out.append("<div class=\"avatar\">");
		out.append("<a href=\"/\">");
		out.append("<img src=\"/comic_150_150.png\" />");
		out.append("</a>");
		out.append("</div>");// avatar
		out.append("<div class=\"metaText\">" + Config.Phrase + "</div>");
		out.append("<div class=\"links\">");
		out.append("<a href=\"/feed\" title=\"RSS\">»RSS</a><br />");
		out.append("<a href=\"/download/" + Config.Name
				+ ".crx\" title=\"Chromium Extension\">»Chromium Extension</a>");
		out.append("</div>");// links
		out.append("<div class=\"clear\"></div>");
		out.append("</div>");// widget
	}

}