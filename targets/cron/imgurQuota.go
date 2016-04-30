package cron

import (
	"net/http"
	"time"

	"github.com/koffeinsource/kaffeeshare/data"
	"github.com/koffeinsource/kaffeeshare/imgurClient"
)

// ImgurQuota will query imgur and print the quota to the log
func ImgurQuota(w http.ResponseWriter, r *http.Request) {
	con := data.MakeContext(r)

	imgurclient := imgurclient.Get(con)

	rl, err := imgurclient.GetRateLimit()
	if err != nil {
		con.Log.Errorf("Error in GetRateLimit: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prctClient := (float64)(rl.ClientRemaining) / (float64)(rl.ClientLimit) * 100.0
	con.Log.Infof("Client limit: %v of %v available (%v%%)", rl.ClientRemaining, rl.ClientLimit, prctClient)

	prctUser := (float64)(rl.UserRemaining) / (float64)(rl.UserLimit) * 100.0
	con.Log.Infof("User limit: %v of %v available (%v%%)", rl.UserRemaining, rl.UserLimit, prctUser)

	con.Log.Infof("User reset: %v, i.e. in %v", rl.UserReset, rl.UserReset.Sub(time.Now()))

	w.WriteHeader(http.StatusOK)
}
