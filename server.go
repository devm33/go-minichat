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
	http.HandleFunc("/", index)
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/_ah/channel/disconnected/", disconnect)
}

func chat(w http.ResponseWriter, r *http.Request) {

}

func index(w http.ResponseWriter, r *http.Request) {
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
		if err := indexTemplate.Execute(w, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

var indexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<div id="messages"></div>
<form id="message-form">
  <input type="text" id="message">
  <button type="submit">Send</button>
</form>
<script src="/_ah/channel/jsapi"></script>
<script>
  window.channel = new goog.appengine.Channel('{{.}}');
</script>
<script src="/static/main.js"></script>
`))
