package shortener

import (
	"fmt"
	"net/http"
)

const addForm = `
<html><body>
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
</html></body>
`

type UrlShortener struct {
	storage Storage
}

func NewUrlShortnener(storage Storage) *UrlShortener {
	return &UrlShortener{storage: storage}
}

func (us *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	if key == "" {
		http.NotFound(w, r)
		return
	}
	var url string
	if err := us.storage.Get(&key, &url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (us *UrlShortener) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprint(w, addForm)
		return
	}

	var key string
	if err := us.storage.Put(&url, &key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", key)
}
