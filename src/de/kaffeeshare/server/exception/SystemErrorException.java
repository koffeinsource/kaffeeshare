package de.kaffeeshare.server.exception;

/**
 * Some unexpected system error.
 */
public class SystemErrorException extends RuntimeException {

	private static final long serialVersionUID = 6358453493328352777L;

	public SystemErrorException() {
		super("Error resulting from unexpected system behaviour!");
	}
}
