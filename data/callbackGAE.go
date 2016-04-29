package data

import (
	"strings"

	"google.golang.org/appengine/datastore"
)

const callbackTable = "Callback"

// Store stores an item in the datastore
func (c *Callback) Store(con *Context) error {
	c.Namespace = strings.ToLower(c.Namespace)
	k := datastore.NewKey(con.C, callbackTable, c.Namespace+c.URL, 0, nil)
	_, err := datastore.Put(con.C, k, c)
	if err != nil {
		con.Log.Errorf("Error while storing callback in datastore. Callback: %v. Error: %v", c, err)
		return err
	}
	c.DSKey = k.String()
	con.Log.Debugf("Stored callback %+v", c)

	/* TODO
	if err := clearCache(con, i.Namespace); err != nil {
		con.Log.Infof("Error clearing cache for namespace %v. Error: %v", i.Namespace, err)
	}*/

	return nil
}

// GetCallbacksByNamespace returns all callbacks for a specific namespace
func GetCallbacksByNamespace(con *Context, namespace string) ([]Callback, error) {
	namespace = strings.ToLower(namespace)
	q := datastore.NewQuery(callbackTable).
		Filter("Namespace =", namespace)

	return executeCallbackQuery(con, q, 10)
}

func executeCallbackQuery(con *Context, q *datastore.Query, limit int) ([]Callback, error) {
	var is = make([]Callback, 0, limit)
	var err error
	t := q.Run(con.C)
	for {
		var i Callback
		_, err = t.Next(&i)
		if err == datastore.Done {
			break
		}

		is = append(is, i)
		if err != nil {
			con.Log.Errorf("Error fetching next item: %v", err)
			return nil, err
		}
	}

	return is, nil
}
