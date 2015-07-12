package check

import (
	"encoding/json"
	"net/http"

	"github.com/koffeinsource/notreddit/data"
	"github.com/koffeinsource/notreddit/targets"
	"github.com/koffeinsource/notreddit/targets/startpage"

	"appengine"
)

const (
	statusOk    = "ok"
	statusError = "error"
	statusInUse = "use"
)

// DispatchJSON executes all commands for the www target
func DispatchJSON(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// get namespace
	namespace := targets.GetNamespace(r, "/k/check/json/")
	if namespace == "" {
		startpage.Dispatch(w, r)
		return
	}

	var temp map[string]string

	inUse, err := data.NamespaceIsEmpty(c, namespace)

	if err != nil {
		temp = map[string]string{"status": statusError}
	} else {

		if inUse {
			temp = map[string]string{"status": statusOk}
		} else {
			temp = map[string]string{"status": statusInUse}
		}

	}

	b, err := json.Marshal(temp)
	if err != nil {
		c.Errorf("Error at marshalling for check/json. Namespace: %v. Error: %v", namespace, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
	c.Infof("JSON result for namespace %v is %v", namespace, string(b))
}
