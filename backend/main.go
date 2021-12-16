package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	chi "github.com/go-chi/chi/v5"
	"github.com/headzoo/surf"
	"golang.org/x/net/html"
)

func main() {
	r := chi.NewRouter()
	r.Get("/*", post)
	http.ListenAndServe(":8080", r)
}

func post(w http.ResponseWriter, r *http.Request) {
	mediumPostURL := chi.URLParam(r, "*")
	bow := surf.NewBrowser()
	mediumOpenError := bow.Open(mediumPostURL)
	if mediumOpenError != nil {
		http.Error(w, fmt.Sprintf("error request medium post: %s", mediumOpenError.Error()), http.StatusInternalServerError)
		return
	}

	author := bow.Find("meta").FilterFunction(func(i int, s *goquery.Selection) bool {
		attr, found := s.Attr("name")
		return found && attr == "author"
	})

	//Waiting for article to load, since it's dynamically loaded via JS.
	var article *goquery.Selection
	for {
		if article != nil {
			break
		}

		time.Sleep(1 * time.Second)
		article = bow.Find("article")
	}

	//Remove author and post meta data from view
	//FIXME Find safer way to do this.
	fmt.Println(article.Find("h1").First().Next().Remove())

	//Since we aren't reusing the stylesheet, the classes are just bloat.
	removeClutter(article)

	//Sections will cause rendering issues. For example h1 and h2
	//inside of a section will look the same in firefox. While this
	//will flatten the post, it won't cause any issues for now.
	article = article.Find("section").Children()

	w.WriteHeader(http.StatusOK)

	//FIXME For some reason toHtml won't work.
	var authorName bytes.Buffer
	html.Render(&authorName, author.Get(0))

	//Body is required for browser to not just render as text.
	w.Write([]byte("<html><head>" + authorName.String() + "</head><body>" + toHtml(article) + "</body></html>"))
}

func toHtml(selection *goquery.Selection) string {
	html, _ := selection.Html()
	return html
}

func removeClutter(selection *goquery.Selection) {
	//Classes and IDs are unnecessary and only waste bandwidth, memory and drive space.
	selection.RemoveClass("")
	selection.RemoveAttr("id")

	children := selection.Children()
	if children.Length() > 0 {
		removeClutter(children)
	}
}
