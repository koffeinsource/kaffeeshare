package de.kaffeeshare.server.input;

import java.io.IOException;
import java.util.logging.Logger;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.google.appengine.api.xmpp.Message;
import com.google.appengine.api.xmpp.MessageBuilder;
import com.google.appengine.api.xmpp.SendResponse;
import com.google.appengine.api.xmpp.XMPPService;
import com.google.appengine.api.xmpp.XMPPServiceFactory;

import de.kaffeeshare.server.UrlImporter;

/**
 * Servlet to handle incomming jabber / xmpp messages.
 */
public class Jabber extends HttpServlet {

	private static final long serialVersionUID = -6000464462591650319L;
	private static final Logger log = Logger.getLogger(Jabber.class.getName());

	private static final XMPPService xmppService = XMPPServiceFactory.getXMPPService();

	/**
	 * This is called whenever a jabber message is sent to the app.
	 */
	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws IOException {
		Message message;
		try {
			message = xmppService.parseMessage(req);
		} catch (Exception e) {
			// http://code.google.com/p/googleappengine/issues/detail?id=2082&can=5&colspec=ID%20Type%20Component%20Status%20Stars%20Summary%20Language%20Priority%20Owner%20Log
			// there may be messages that are no messages, we ignore them
			// for now
			return;
		}

		log.info("got XMPP message from: " + message.getFromJid());
		log.info("Message body: " + message.getBody());

		importUrl(message);
	}

	private void importUrl(Message message) {
		String replyMessageBody = null;

		String receiverID = message.getRecipientJids()[0].getId();
		// is message sent to <AppID>@appspot.com              <- namespace: default
		// or is it   sent to anything@<AppID>.appspotchat.com <- namespace: anything
		if (!receiverID.contains("@appspot.com") ) {
			receiverID = receiverID.split("@")[0];
			//UrlImporter.setNamespace(receiverID);
		} else {
			//UrlImporter.setNamespace(null);
		}

		try {
			String url = UrlImporter.importFromText(message.getBody());

			if (url != null)
				replyMessageBody = "URL " + url + " added to DB.";
			else
				replyMessageBody = "No URL found!";

		} catch (Exception e) {
			// server error
			replyMessageBody = "Server error, should not happen!";
		}

		Message msg = new MessageBuilder()
		                  .withRecipientJids(message.getFromJid())
		                  .withBody(replyMessageBody).build();

		SendResponse status = xmppService.sendMessage(msg);
		boolean messageSent = (status.getStatusMap().get(message.getFromJid()) == SendResponse.Status.SUCCESS);

		if (!messageSent)
			log.warning("Reply message could not be sent, this should not happen(?)");
	}
}
