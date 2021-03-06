package minichat

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/channel"
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

func disconnect(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	userid := r.FormValue("from")
	if userid != "" {
		key := datastore.NewKey(c, "ActiveUser", userid, 0, nil)
		err := datastore.Delete(c, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if r.Body != nil {
		message, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt_msg := fmt.Sprintf("%v: %v", u.String(), string(message))
		var actives []ActiveUser
		_, err = datastore.NewQuery("ActiveUser").GetAll(c, &actives)
		for _, active := range actives {
			channel.Send(c, active.Userid, fmt_msg)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	// Check if user is logged in
	if u != nil {
		// Create unique chat channel for user and save to active list
		token, err := channel.Create(c, u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Add user to active users
		a := ActiveUser{
			Userid: u.ID,
		}
		key := datastore.NewKey(c, "ActiveUser", u.ID, 0, nil)
		_, err = datastore.Put(c, key, &a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := indexTemplate.Execute(w, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		l, err := user.LoginURL(c, "/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, l, 302)
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
