package minichat

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type ActiveUser struct {
	Userid string
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/_ah/channel/disconnected/", disconnect)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// Check if user is logged in
	if u != nil {
		// Create unique chat channel for user and save to active list
		token, err := channel.Create(c, user.user_id())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Add user to active users
		a := ActiveUser{
			Userid: u.ID,
		}
		key := datastore.NewKey(c, ActiveUser, u.ID, 0, nil)
		_, err := datastore.Put(c, key, &a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
