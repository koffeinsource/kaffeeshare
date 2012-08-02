package de.kaffeeshare.server.datastore;

import java.util.Arrays;
import java.util.List;
import java.util.logging.Logger;

import com.google.appengine.api.NamespaceManager;

import de.kaffeeshare.server.exception.IllegalNamespaceException;
import de.kaffeeshare.server.exception.ReservedNamespaceException;

public class Namespace {
	
	private static final Logger log = Logger.getLogger(Namespace.class.getName());

	private static List<String> reservedNamespaces = Arrays.asList("public", "private", "kaffee");
	
	static public void setNamespace(String ns) {
		validateNamespace(ns);
		
		log.info("Namespace set to " + ns);
		NamespaceManager.set(ns);
	}
	
	static public void validateNamespace(String ns) {
		try {
			NamespaceManager.validateNamespace(ns);
		} catch (IllegalArgumentException e) {
			log.info("Trying to use an illegal namespace: " + ns);
			throw new IllegalNamespaceException();
		}
		
		if (reservedNamespaces.contains(ns) || ns == null) {
			log.info("Trying to use a reserved namespace: " + ns);
			throw new ReservedNamespaceException();
		}
	}
	
	static public boolean isValide(String ns)  {
		try {
			validateNamespace(ns);
		} catch (Exception e) {
			return false;
		}
		
		return true;
	}
}
