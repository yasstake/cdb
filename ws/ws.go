package ws

import (
	"cdb/trans"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
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

func TickerInfoHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")

	t.Execute(w, "")
}

func HtmlInfoHandler(w http.ResponseWriter, r *http.Request) {
	env.open("")

	chunks := env.db.GetTimeChunks()

	t, _ := template.ParseFiles("template/info_templ.html")
	t.Execute(w, chunks)
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

func DbInfoHandler(w http.ResponseWriter, r *http.Request) {
	env.open("")
	f := r.FormValue("from")
	fmt.Fprintln(w, f)
	t := r.FormValue("to")
	fmt.Fprintln(w, t)
	fmt.Fprintf(w, "Hello!!")
	fmt.Fprintf(w, r.URL.Path)

}

func TransactinCSVHandler(w http.ResponseWriter, r *http.Request) {
	env.open("")

	start_time := env.db.GetStartTime()

	session, _ := env.db.CreateSession(start_time)

	reader, err := session.SelectTrans(start_time, start_time.Add(60*time.Minute))

	limit_string := r.FormValue("limit")
	limit := 1000

	if limit_string != "" {
		limit, _ = strconv.Atoi(limit_string)
	}

	if err != nil {
		log.Println(err)
	}

	log.Println("GET transaction csv record", limit)

	for i := 0; i < limit; i++ {
		r, err := reader.ReadTran()

		if err != nil {
			break
		}
		s := strconv.Itoa(int(r.Action)) + "," + strconv.Itoa(int(r.Time_stamp)) + "," +
			strconv.Itoa(int(r.Price)) + "," + strconv.Itoa(int(r.Volume)) + "," +
			strconv.Itoa(int(r.OtherInfo))
		fmt.Fprintln(w, s)
	}

}
