package de.kaffeeshare.server.exception;

/**
 * Thrown if a reserved namespace is supposed to be used
 */
public class ReservedNamespaceException extends RuntimeException {

	private static final long serialVersionUID = 2826030527124116551L;
	
	public ReservedNamespaceException() {
		super(Messages.getString("ReservedNamespaceException.msg"));
	}

}
