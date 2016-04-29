package data

// A Callback is something that should be executed after sharing an url
type Callback struct {
	URL       string `datastore:"URL,index"`
	Namespace string `datastore:"Namespace,index"`
	DSKey     string `datastore:"-"`
}
