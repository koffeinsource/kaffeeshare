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
package de.kaffeeshare.server.input;

import java.io.IOException;
import java.util.logging.Logger;

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.google.appengine.api.xmpp.JID;
import com.google.appengine.api.xmpp.Message;
import com.google.appengine.api.xmpp.MessageBuilder;
import com.google.appengine.api.xmpp.SendResponse;
import com.google.appengine.api.xmpp.XMPPService;
import com.google.appengine.api.xmpp.XMPPServiceFactory;

import de.kaffeeshare.server.UrlImporter;
import de.kaffeeshare.server.datastore.DatastoreManager;

/**
 * Servlet to handle incomming jabber / xmpp messages.
 */
public class Jabber extends HttpServlet {

	private static final long serialVersionUID = -6000464462591650319L;
	private static final Logger log = Logger.getLogger(Jabber.class.getName());

	private static final XMPPService xmppService = XMPPServiceFactory.getXMPPService();

	/**
	 * Handle a post request.
	 * This is called whenever a jabber message is sent to the app.
	 * @param req Request
	 * @param resp Response
	 * @throws IOException
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

	/**
	 * Import a URL.
	 * @param message Message
	 */
	private void importUrl(Message message) {
		String replyMessageBody = null;

		// TODO will break if there are multiple recipients of that message
		JID recieverID = message.getRecipientJids()[0];
		String recieverIDStr =recieverID.getId();
		// is message sent to <AppID>@appspot.com              <- namespace: default
		// or is it   sent to anything@<AppID>.appspotchat.com <- namespace: anything
		if (!recieverIDStr.contains("@appspot.com") ) {
			recieverIDStr = recieverIDStr.split("@")[0];
			DatastoreManager.setNamespace(recieverIDStr);
		} else {
			DatastoreManager.setNamespace(null);
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
		                  .withFromJid(recieverID)
		                  .withRecipientJids(message.getFromJid())
		                  .withBody(replyMessageBody).build();

		SendResponse status = xmppService.sendMessage(msg);
		boolean messageSent = (status.getStatusMap().get(message.getFromJid()) == SendResponse.Status.SUCCESS);

		if (!messageSent)
			log.warning("Reply message could not be sent, this should not happen(?)");
	}
}
