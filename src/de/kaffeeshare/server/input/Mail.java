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
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;
import java.util.logging.Logger;

import javax.mail.Address;
import javax.mail.Message;
import javax.mail.MessagingException;
import javax.mail.Multipart;
import javax.mail.Part;
import javax.mail.Session;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import de.kaffeeshare.server.UrlImporter;
import de.kaffeeshare.server.datastore.DatastoreManager;
import de.kaffeeshare.server.exception.InputErrorException;
import de.kaffeeshare.server.exception.SystemErrorException;

/**
 * Servlet to handle incomming emails.
 */
public class Mail extends HttpServlet {

	private static final long serialVersionUID = 294584452111372279L;
	private static final Logger log = Logger.getLogger(Mail.class.getName());

	/**
	 * Handle a post request.
	 * @param req Request
	 * @param resp Response
	 * @throws ServletException, SystemErrorException
	 */
	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, SystemErrorException {
        Properties props = new Properties(); 
        Session session = Session.getDefaultInstance(props, null); 
        try {
            MimeMessage message = new MimeMessage(session, req.getInputStream());			
            
			log.info("Got email from " + ((InternetAddress)message.getFrom()[0]).getAddress());
			
			List<String> toAddresses = new ArrayList<String>();
			Address[] recipients = message.getRecipients(Message.RecipientType.TO);
			for (Address address : recipients) {
			    toAddresses.add(address.toString());
			}
			
			for (String to : toAddresses) {
			
				// check if to is name <address@domain.tld>
				if (to.contains("<")) {
					int start = to.indexOf("<");
					int end = to.indexOf(">");
					if (start == -1 || end == -1) throw new InputErrorException();
					++start;
					--end;
					to = to.substring(start, end);
				}
				
				to = to.split("@")[0];
				DatastoreManager.setNamespace(to);
				
				// first lets see if there is plain text with url
				if (UrlImporter.importFromText(getText(message)) != null) continue;
				
				log.info("No URLs found in plain text, will look in HTML");
				
				UrlImporter.importFromHTML(getHTML(message));
			}
			
		} catch (MessagingException e) {
			throw new SystemErrorException();
		} catch (IOException e) {
			throw new SystemErrorException();
		}

	}
	
	static private String getText(Part p) {
		return getMime(p, "text/plain");
	}
	
	static private String getHTML(Part p) {
		return getMime(p, "text/html");
	}
	
	/**
	 * Return the text content of the message with the matching mime type.
	 * @param p Part
	 * @param mime Mime type
	 * @return Text content
	 * @throws InputErrorException, SystemErrorException
	 */
	static private String getMime(Part p, String mime) throws InputErrorException, SystemErrorException {
		try {
			// ok is p text? yes => return body
			if (p.isMimeType(mime)) {
				log.info("Found text in email ("+mime+"):");
				log.info((String)p.getContent());
				return (String)p.getContent();
			}

			StringBuffer sb = new StringBuffer();
			
			// multipart/alternative? so we do have the text multiple times
			if (p.isMimeType("multipart/alternative")) {
	
				log.info("Found multipart/alternative in email:");
				
				Multipart mp = (Multipart)p.getContent();
				
				for (int i = 0; i < mp.getCount(); i++) {
					Part bp = mp.getBodyPart(i);
					
					String s = getMime(bp, mime);
					if (s != null) sb.append(" " + s);
				}
				
			} else if (p.isMimeType("multipart/*")) {
				log.info("Found multipart/* in email:");
				
				Multipart mp = (Multipart)p.getContent();
				for (int i = 0; i < mp.getCount(); i++) {
					String s = getMime(mp.getBodyPart(i), mime);
					if (s != null) sb.append(" " + s);
				}
			}
	
			return sb.toString();
		} catch (MessagingException e) {
			throw new InputErrorException();
		} catch (IOException e) {
			throw new SystemErrorException();
		}
	}
}
