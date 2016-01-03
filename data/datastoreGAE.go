package data

import (
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Store stores an item in the datastore
func (i *Item) Store(c context.Context) error {
	i.Namespace = strings.ToLower(i.Namespace)
	k := datastore.NewKey(c, "Item", i.Namespace+i.URL, 0, nil)
	_, err := datastore.Put(c, k, i)
	if err != nil {
		log.Errorf(c, "Error while storing item in datastore. Item: %v. Error: %v", i, err)
		return err
	}
	i.DSKey = k.String()
	log.Debugf(c, "Stored item %+v", i)

	if err := clearCache(c, i.Namespace); err != nil {
		log.Infof(c, "Error clearing cache for namespace %v. Error: %v", i.Namespace, err)
	}

	return nil
}

// GetNewestItems returns the latest number elements for a specific namespace
func GetNewestItems(c context.Context, namespace string, limit int, cursor string) ([]Item, string, error) {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery("Item").
		Filter("Namespace =", namespace).
		Order("-CreatedAt").
		Limit(limit)

	return executeQuery(c, q, limit, cursor)
}

// GetNewestItemsByTime returns up to limit numbers of items stored >= the
// give time
func GetNewestItemsByTime(c context.Context, namespace string, limit int, t time.Time, cursor string) ([]Item, string, error) {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery("Item").
		Filter("Namespace =", namespace).
		Filter("CreatedAt >=", t).
		Order("-CreatedAt").
		Limit(limit)

	return executeQuery(c, q, limit, cursor)
}

func executeQuery(c context.Context, q *datastore.Query, limit int, cursor string) ([]Item, string, error) {
	if cursor, err := datastore.DecodeCursor(cursor); err == nil {
		q = q.Start(cursor)
	}

	var is = make([]Item, 0, limit)
	var err error
	t := q.Run(c)
	for {
		var i Item
		_, err = t.Next(&i)
		if err == datastore.Done {
			break
		}

		is = append(is, i)
		if err != nil {
			log.Errorf(c, "Error fetching next item: %v", err)
			return nil, "", err
		}
	}

	if cursor, err := t.Cursor(); err == nil {
		return is, cursor.String(), nil
	}

	return nil, "", err
}

// NamespaceIsEmpty checks if there is already an item in a namespace
func NamespaceIsEmpty(c context.Context, namespace string) (bool, error) {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery("Item").
		Filter("Namespace =", namespace).
		Limit(1).
		KeysOnly()

	k, err := q.GetAll(c, nil)
	b := len(k) == 0

	return b, err
}

// ClearNamespace deletes every entry in a namespace
func ClearNamespace(c context.Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery("Item").
		Filter("Namespace =", namespace).
		KeysOnly()

	k, err := q.GetAll(c, nil)
	if err != nil {
		return err
	}

	clearCache(c, namespace)

	log.Infof(c, "Going to delete %v items in the namespace %v", len(k), namespace)

	return datastore.DeleteMulti(c, k)
}

// DeleteAllItems deletes all items from datastore
func DeleteAllItems(c context.Context) error {
	panic("Are you sure????!!!!")
	/*q := datastore.NewQuery("Item").KeysOnly()

	k, err := q.GetAll(c, nil)
	if err != nil {
		return err
	}

	return datastore.DeleteMulti(c, k)*/
}
