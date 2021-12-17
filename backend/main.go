package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	chi "github.com/go-chi/chi/v5"
	"github.com/headzoo/surf"
	"github.com/kennygrant/sanitize"
	"golang.org/x/net/html"
)

var cacheDir string

func main() {
	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println("Couldn't find cache directory.")
		panic(err)
	}
	cacheDir = filepath.Join(baseCacheDir, "unmedium")
	cacheDirError := os.MkdirAll(cacheDir, 0777)
	if cacheDirError != nil {
		log.Printf("Error creating cache direction %s:\n\t%s\n", cacheDir, cacheDirError)
	}
	log.Printf("Caching into %s\n", cacheDir)

	r := chi.NewRouter()
	r.Get("/*", post)
	http.ListenAndServe(":8080", r)
}

func getPost(postURL string) (*goquery.Selection, *goquery.Selection, error) {
	baseName := sanitize.BaseName(postURL)
	postFileName := filepath.Join(cacheDir, baseName)
	file, openError := os.Open(postFileName)
	if openError != nil {
		bow := surf.NewBrowser()
		mediumOpenError := bow.Open(postURL)
		if mediumOpenError != nil {
			return nil, nil, mediumOpenError
		}

		dom := bow.Dom()
		//Waiting for article to load, since it's dynamically loaded via JS.
		//This has to happen before caching, since we might cache it without the article otherwise.
		var article *goquery.Selection
		for {
			if article != nil {
				break
			}

			time.Sleep(100 * time.Millisecond)
			article = dom.Find("article")
		}

		if os.IsNotExist(openError) {
			log.Printf("Caching %s\n", postURL)
			cacheWriteError := os.WriteFile(postFileName, []byte(toHtml(dom)), 0666)
			if cacheWriteError != nil {
				log.Printf("Error caching %s:\n\t%s\n", postURL, cacheWriteError)
			}
		} else {
			log.Printf("Error caching %s:\n\t%s\n", postURL, openError)
		}

		return dom, article, nil
	}

	defer file.Close()

	log.Printf("Cache hit for %s\n", postURL)
	document, readError := goquery.NewDocumentFromReader(file)
	if readError != nil {
		return nil, nil, readError
	}
	return document.Selection, document.Find("article"), nil
}

func post(w http.ResponseWriter, r *http.Request) {
	mediumPostURL := chi.URLParam(r, "*")
	document, article, postRetrievalError := getPost(mediumPostURL)
	if postRetrievalError != nil {
		http.Error(w, fmt.Sprintf("error retrieving medium post: %s", postRetrievalError.Error()), http.StatusInternalServerError)
		return
	}

	authorTag := document.Find("meta[name=author]").First()

	//Remove author and post meta data from view
	//FIXME Find safer way to do this.
	article.Find("h1").First().Next().Remove()

	//Since we aren't reusing the stylesheet, the classes are just bloat.
	removeClutter(article)

	//Sections will cause rendering issues. For example h1 and h2
	//inside of a section will look the same in firefox. While this
	//will flatten the post, it won't cause any issues for now.
	article.Find("section > *").Unwrap()
	article.Find("section").Remove()

	var html bytes.Buffer
	html.WriteString("<html><head>")
	if authorTag.Length() == 1 {
		html.WriteString(renderNode(authorTag.Get(0)))
	}
	//Body is required for browser to not just render as text.
	html.WriteString("</head><body>")
	html.WriteString(toHtml(article))
	html.WriteString("</body></html>")

	w.WriteHeader(http.StatusOK)
	w.Write(html.Bytes())
}

func renderNode(node *html.Node) string {
	var authorName bytes.Buffer
	html.Render(&authorName, node)
	return authorName.String()
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
