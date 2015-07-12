// +build appengine

package data

import (
	"appengine"
	"appengine/datastore"
)

// StoreItem stores an item in the datastore
func StoreItem(c appengine.Context, i Item) error {
	k := datastore.NewKey(c, "Item", i.Namespace+i.URL, 0, nil)
	_, err := datastore.Put(c, k, &i)
	c.Infof("Stored item %v", i)
	if err != nil {
		c.Errorf("Error while storing item in datastore. Item: %v. Error: %v", i, err)
	}

	return err
}

// GetNewestItems returns the latest number elements for a specific namespace
func GetNewestItems(c appengine.Context, namespace string, limit int) ([]Item, error) {
	q := datastore.NewQuery("Item").
		Filter("Namespace =", namespace).
		Order("-CreatedAt").
		Limit(limit)

	var is []Item
	_, err := q.GetAll(c, &is)

	return is, err
}

// DeleteAllItems deletes all items from datastore
func DeleteAllItems(c appengine.Context) error {
	q := datastore.NewQuery("Item").KeysOnly()

	k, err := q.GetAll(c, nil)
	if err != nil {
		return nil
	}

	return datastore.DeleteMulti(c, k)
}
