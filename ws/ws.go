package ws

import (
	"cdb/trans"
	"fmt"
	"html/template"
	"net/http"
)

//
//  ROOT
//   +  Database info (HTML)
//     +html
//         + info
//     +api
//         + tran
//         + tran
//         + tran
//
//
//

type ServerEnv struct {
	db      trans.Db
	db_open bool
}

func (c *ServerEnv) open(db_path string) {
	if c.db_open == true {
		return
	}
	if db_path == "" {
		c.db.Open("/tmp")
	} else {
		c.db.Open(db_path)
	}
	c.db_open = true
}

var env ServerEnv

func RootHtmlHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, "")
}

func HtmlInfoHandler(w http.ResponseWriter, r *http.Request) {
	env.open("")

	chunks := env.db.GetTimeChunks()

	t, _ := template.ParseFiles("template/info_templ.html")
	t.Execute(w, chunks.ToString())
	fmt.Fprint(w, chunks.ToString())
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	env.open("")
	f := r.FormValue("from")
	fmt.Fprintln(w, f)
	t := r.FormValue("to")
	fmt.Fprintln(w, t)
	fmt.Fprintf(w, "Hello!!")
	fmt.Fprintf(w, r.URL.Path)
}
