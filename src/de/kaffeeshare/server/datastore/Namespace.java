package de.kaffeeshare.server.datastore;

import java.util.Arrays;
import java.util.List;
import java.util.logging.Logger;

import com.google.appengine.api.NamespaceManager;

import de.kaffeeshare.server.exception.ReservedNamespaceException;

public class Namespace {
	
	private static final Logger log = Logger.getLogger(Namespace.class.getName());

	private static List<String> reservedNamespaces = Arrays.asList("public", "private", "kaffee");
	
	static public void setNamespace(String ns) {
		
		if (reservedNamespaces.contains(ns) || ns == null) {
			log.info("Trying to use reserve namespace!");
			
			throw new ReservedNamespaceException();
		}
		
		log.info("Namespace set to " + ns);
		NamespaceManager.set(ns);
	}
}
