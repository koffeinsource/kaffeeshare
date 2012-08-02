package de.kaffeeshare.server.exception;

public class IllegalNamespaceException extends RuntimeException {
	
	private static final long serialVersionUID = -3633614066774348135L;

	public IllegalNamespaceException() {
		super("Trying to use an illegal namespace!");
	}

}
