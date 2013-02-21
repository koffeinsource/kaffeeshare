package de.kaffeeshare.server.exception;

/**
 * Thrown if there is a problem with the data store.
 */
public class DSError extends RuntimeException {

	private static final long serialVersionUID = -6898169722430572339L;

	public DSError() {
		super (Messages.getString("DSError.msg"));
	}
}
