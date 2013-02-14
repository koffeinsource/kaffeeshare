package de.kaffeeshare.server.datastore;

import java.util.Arrays;
import java.util.List;
import java.util.logging.Logger;

import com.google.appengine.api.NamespaceManager;

import de.kaffeeshare.server.exception.IllegalNamespaceException;
import de.kaffeeshare.server.exception.ReservedNamespaceException;

/**
 * Namespace class to handle namespaces.
 */
public class NamespaceValidator {
	
	private static final Logger log = Logger.getLogger(NamespaceValidator.class.getName());

	/**
	 * A list of namespaces we want to reserver for special use
	 */
	private static List<String> reservedNamespaces = Arrays.asList("public", "private", "kaffee");
	
	/**
	 * Validates a namespace. An exception if thrown in case it is not valid
	 * @param ns the namespace to be validated
	 */
	private static void validate(String ns) {
		if (ns == null || reservedNamespaces.contains(ns)) {
			log.info("Trying to use a reserved namespace");
			throw new ReservedNamespaceException();
		}
		try {
			NamespaceManager.validateNamespace(ns);
		} catch (IllegalArgumentException e) {
			log.info("Trying to use an illegal namespace: " + ns);
			throw new IllegalNamespaceException();
		}
	}
	
	/**
	 * Checks if a namespace is valid
	 * @param ns the namespace to be validated
	 * @return true if namespace is valid, false otherwise
	 */
	static public boolean isValide(String ns)  {
		try {
			validate(ns);
		} catch (Exception e) {
			return false;
		}
		
		return true;
	}
}
