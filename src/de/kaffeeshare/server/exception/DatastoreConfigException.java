package de.kaffeeshare.server.exception;


/**
 * No datastore interface is defined.
 */
public class DatastoreConfigException extends RuntimeException {

	private static final long serialVersionUID = -2171829811099619421L;

	public DatastoreConfigException() {
		super(Messages.getString("DatastoreConfigException.msg"));
	}
}
