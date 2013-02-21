package de.kaffeeshare.server.exception;

/**
 * Thrown if a reserved namespace name is not legal 
 */
public class IllegalNamespaceException extends RuntimeException {
	
	private static final long serialVersionUID = -3633614066774348135L;

	public IllegalNamespaceException() {
		super(Messages.getString("IllegalNamespaceException.msg"));
	}

}
