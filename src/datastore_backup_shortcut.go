package datastore_backup_shortcut

import (
	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ac := appengine.NewContext(r)
	ac.Infof(r.URL.String())

	f(ac, r.URL.String()[61:len(r.URL.String())])
	fmt.Fprint(w, "Done")
}

func f(c appengine.Context, url string) {
	path := fmt.Sprintf("/_ah/datastore_admin/backup.create?name=BackupToCloud&filesystem=gs&%s", url)
	c.Infof("Path : %s", path)
	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		t := taskqueue.Task{
			Path:   path,
			Method: "GET",
		}
		_, err := taskqueue.Add(c, &t, "datastore-backup-shortcut")
		return err
	}, nil)

	if err != nil {
		c.Errorf(err.Error())
	}
}
