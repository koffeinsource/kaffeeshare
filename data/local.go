// +build !appengine

// TODO just stubs to compile it without app engine

package data

import "errors"

// StoreItem stores an item in the datastore
func StoreItem(c interface{}, i Item) error {
	return errors.New("not implemented")
}

// GetNewestItems returns the latest number elements for a specific namespace
func GetNewestItems(c interface{}, namespace string, limit int) ([]Item, error) {
	return nil, errors.New("not implemented")
}

// DeleteAllItems deletes all items from datastore
func DeleteAllItems(c interface{}) error {
	return errors.New("not implemented")
}
