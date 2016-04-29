package data

import (
	"strings"

	"google.golang.org/appengine/datastore"
)

// NamespaceIsEmpty checks if there is already an item in a namespace
func NamespaceIsEmpty(con *Context, namespace string) (bool, error) {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery(itemTable).
		Filter("Namespace =", namespace).
		Limit(1).
		KeysOnly()

	k, err := q.GetAll(con.C, nil)
	b := len(k) == 0

	return b, err
}

// ClearNamespace deletes every entry in a namespace
func ClearNamespace(con *Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery(itemTable).
		Filter("Namespace =", namespace).
		KeysOnly()

	k, err := q.GetAll(con.C, nil)
	if err != nil {
		return err
	}

	clearCache(con, namespace)

	con.Log.Infof("Going to delete %v items in the namespace %v", len(k), namespace)

	return datastore.DeleteMulti(con.C, k)
}
