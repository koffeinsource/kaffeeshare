/*******************************************************************************
 * Copyright 2013 Jens Breitbart, Daniel Klemm, Dennis Obermann
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 ******************************************************************************/
package de.kaffeeshare.server.exception;

/**
 * Thrown if no datastore interface is defined.
 */
public class DatastoreConfigException extends RuntimeException {

	private static final long serialVersionUID = -2171829811099619421L;

	/**
	 * Void constructor.
	 */
	public DatastoreConfigException() {
		super("No datastore interface is defined (See config.properties)!");
	}
}
