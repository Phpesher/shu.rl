package main

/**
 * Copy. 2020-2020
 * Nikita ( Phpesher )
**/

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	conf "shu.rl/conf"
	m "shu.rl/models"
	g "shu.rl/pkg/gens"
)

/**
 * Temporary storage
 * urls : temporary storage data into memory
 * db : data base "object"
 * sh : short url in memory
 * sc : source url in memory
**/
var (
	urls = make(map[string]*m.Url)
	db   = m.NewDatabase("mysql", "mysql", "localhost")
	sh   = ""
	sc   = ""
	c    = *conf.NewConfig(":8080", "localhost")
)

/**
 * This is main method shu.rl server
**/
func main() {
	// For correct css
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("./www/"))))

	router := mux.NewRouter()

	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/about", AboutHandler)
	router.HandleFunc("/s/{match}", RedirectOnShortToSourceUrl)
	router.HandleFunc("/short", ShortHandler)
	router.HandleFunc("/shortUrl", ShortUrlHandler)

	router.NotFoundHandler = http.HandlerFunc(ErrorHandler)

	// WEB handle
	http.Handle("/", router)

	port := c.ServerPort

	fmt.Printf("Listening on port %s\n", port)
	fmt.Printf("http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", port), nil))

}

/**
 * Web Handler shorting template
**/
func ShortHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./www/tmp/index.html", "./www/tmp/header.html", "./www/tmp/footer.html")

	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
	}

	err = t.ExecuteTemplate(w, "index", urls)
}

/**
 * Work with model and save short and source url in memory and database
**/
func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	id        := g.GenerateId()
	SourceUrl := r.FormValue("url")
	ShortUrl  := g.GenerateShortUrl()

	if FindSourceUrlInDb(SourceUrl, "no") != "" {
		rows, err := db.DATABASE.Query("select * from " + c.DataBaseTable + " where u_id = ?", FindSourceUrlInDb(SourceUrl, "no"))

		for rows.Next() {
			var id int
			var uId string
			var sourceUrl string
			var shortUrl string

			err = rows.Scan(&id, &uId, &sourceUrl, &shortUrl)

			if err != nil {
				fmt.Println(err.Error())
			}
			sh = "http://" + c.ServerHost + c.ServerPort + "/s/" + shortUrl
			sc = sourceUrl

			http.Redirect(w, r, sh, 302)
		}

		if err != nil {
			fmt.Println(err.Error())
		}

		defer rows.Close()
	} else {
		newUrl := m.NewUrl(id, SourceUrl, ShortUrl)
		urls[newUrl.Id] = newUrl

		for _, i := range urls {
			_, err := db.DATABASE.Exec("insert into " +  c.DataBaseName + "." + c.DataBaseTable + " (id, u_id, source_url, short_url) values(NULL, ?, ?, ?)", i.Id, i.SourceUrl, i.NewUrl)

			sh = "http://" + c.ServerHost + c.ServerPort + "/s/" + i.NewUrl
			sc = i.SourceUrl

			if err != nil {
				fmt.Println(err.Error())
			}
		}

		http.Redirect(w, r, "http://" + c.ServerHost + c.ServerPort + "/short", 302)
	}
}

/**
 * This function find "source" in data base and if find then return he
**/
func FindSourceUrlInDb(sUrl, findShort string) string {
	rows, err := db.DATABASE.Query("select * from " + c.DataBaseTable + " where source_url = ?", sUrl)

	for rows.Next() {
		var id int
		var uId string
		var sourceUrl string
		var shortUrl string

		err = rows.Scan(&id, &uId, &sourceUrl, &shortUrl)

		if err != nil {
			fmt.Println(err.Error())
		}

		if findShort == "yes" {
			return shortUrl
		} else {
			return uId
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()
	return ""
}

/**
 * If requested short url exists into data base, when redirect to source url
**/
func RedirectOnShortToSourceUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println(sc, sh)
	http.Redirect(w, r, sc, 302)
}

/**
 * About page web handler
**/
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./www/tmp/about.html", "./www/tmp/header.html", "./www/tmp/footer.html")

	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
	}

	err = t.ExecuteTemplate(w, "about", nil)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./www/tmp/404.html")

	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
	}

	err = t.ExecuteTemplate(w, "404", nil)
}

/**
 * Index WEB handler
**/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/short", 302)
	}
}