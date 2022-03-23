package main

import (
	"bytes"
	"html/template"
	"log"
	"math"
	"net/url"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/freshman-tech/news-demo/news"
	/*need this for specific api linking
	 stuff i didn't feel like figuring out*/

)

//reading html template for search implementation
var (
	tmpl = template.Must(template.ParseFiles("index.html"))
)

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}


func(s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

func(s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}
//handler for http connection
func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

//obviously search handler
func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		paramies := u.Query()
		searchQuery := paramies.Get("q")
		page := paramies.Get("page")
		if page == "" {
			page = "1"
		}
//newsapi is on newsapi.org
		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}

		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}
//key is in the .env
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets"))
/*forgot i linked these to assets folder and 
forgot to make one, created so many issues*/
mux := http.NewServeMux()
mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
mux.HandleFunc("/search/", searchHandler(newsapi))
mux.HandleFunc("/", indexHandler)
http.ListenAndServe(":"+port, mux)

}
