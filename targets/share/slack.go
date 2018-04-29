package share

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/share"
	"github.com/mvdan/xurls"
)

/*
token=aESya7pNILPaEMfLGJrS4tJM
team_id=T0001
team_domain=example
channel_id=C2147483705
channel_name=test
timestamp=1355517523.000005
user_id=U2147483697
user_name=Steve
text=googlebot: What is the air-speed velocity of an unladen swallow?
trigger_word=googlebot:
*/

type slackPost struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	// ignoring timestamp
	UserID      string
	UserName    string
	Text        string
	TriggerWord string
}

func parseSlackForm(con *data.Context, r *http.Request) (*slackPost, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	con.Log.Infof("%v", r.Form)

	var sp slackPost

	sp.Token = r.Form.Get("token")
	sp.TeamID = r.Form.Get("team_id")
	sp.TeamDomain = r.Form.Get("team_domain")
	sp.ChannelID = r.Form.Get("channel_id")
	sp.ChannelName = r.Form.Get("channel_name")
	sp.UserID = r.Form.Get("user_id")
	sp.UserName = r.Form.Get("user_name")
	sp.Text = r.Form.Get("text")
	sp.TriggerWord = r.Form.Get("trigger_word")

	return &sp, nil
}

// DispatchSlack receives slack messages
func DispatchSlack(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	// get namespace
	namespace := mux.Vars(r)["namespace"]
	if namespace == "" {
		con.Log.Errorf("No namespace provided by slack.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slack, err := parseSlackForm(con, r)
	if err != nil {
		con.Log.Errorf("Error while parsing form: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	con.Log.Infof("Parsed message: %v", slack)

	if slack.Token != "aESya7pNILPaEMfLGJrS4tJM" {
		con.Log.Errorf("Invalid Token: %v", slack.Token)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if slack.UserName == "slackbot" {
		con.Log.Infof("Got a bot message. I'll ignore this one.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO this c&p. we need to extract that from email parsing...
	var links []string
	l := xurls.Relaxed.FindAllString(slack.Text, -1)
	for _, s := range l {
		if s != "" && !strings.Contains(s, "mailto:") {
			links = append(links, s)
		}
	}
	con.Log.Infof("Found %v urls in message %v,  %v", len(links), slack.Text, links)

	if len(links) == 0 {
		con.Log.Infof("Found no URL.")
		return
	}

	shareURL := links[0]

	if err := share.URL(shareURL, namespace, con); err != nil {
		con.Log.Errorf("Error while sharing an URL. URL: %v. Error: %v", shareURL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
