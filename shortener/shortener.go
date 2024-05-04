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

func NewUrlStorage(storage Storage) *UrlShortener {
	return &UrlShortener{storage: storage}
}

func (us *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url := us.storage.Get(key)
	if url == "" {
		http.NotFound(w, r)
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

	key := us.storage.Put(url)
	fmt.Fprintf(w, "%s", key)
}
