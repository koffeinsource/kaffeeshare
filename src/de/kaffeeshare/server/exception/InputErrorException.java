package de.kaffeeshare.server.exception;

/**
 * Thrown when we have an I/O error.
 */
public class InputErrorException extends RuntimeException {

	private static final long serialVersionUID = -6560942737782782844L;

	public InputErrorException() {
		super ("Input error, not what we expected!");
	}
}
