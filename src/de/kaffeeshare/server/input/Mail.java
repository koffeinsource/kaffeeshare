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

	public void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException {
        Properties props = new Properties(); 
        Session session = Session.getDefaultInstance(props, null); 
        try {
            MimeMessage message = new MimeMessage(session, req.getInputStream());			
            
			log.info(Messages.getString("Mail.email_from") + ((InternetAddress)message.getFrom()[0]).getAddress());
			
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
				
				log.info(Messages.getString("Mail.no_url_found"));
				
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
	 */
	static private String getMime(Part p, String mime) {
		try {
			// ok is p text? yes => return body
			if (p.isMimeType(mime)) {
				log.info(Messages.getString("Mail.text_found") + "(" + mime + "):" );
				log.info((String)p.getContent());
				return (String)p.getContent();
			}

			StringBuffer sb = new StringBuffer();
			
			// multipart/alternative? so we do have the text multiple times
			if (p.isMimeType("multipart/alternative")) {
	
				log.info(Messages.getString("Mail.multipart_found_alternative"));
				
				Multipart mp = (Multipart)p.getContent();
				
				for (int i = 0; i < mp.getCount(); i++) {
					Part bp = mp.getBodyPart(i);
					
					String s = getMime(bp, mime);
					if (s != null) sb.append(" " + s);
				}
				
			} else if (p.isMimeType("multipart/*")) {
				log.info(Messages.getString("Mail.multipart_found"));
				
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