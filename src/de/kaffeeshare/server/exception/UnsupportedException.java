package de.kaffeeshare.server.exception;

public class UnsupportedException extends RuntimeException {

	private static final long serialVersionUID = -2323208065883966089L;
	
	public UnsupportedException() {
		super ("Sorry, this feature is not yet implemented.");
	}

}
