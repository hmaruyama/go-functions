package lunch

import (
	"time"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Parameter struct {
	SubCommand string
	Value string
}

type Restaurant struct {
	ID int64 `datastore:"-"`
	Name string `datastore:"name"`
	CreatedAt time.Time `datastore:"created"`
}

func add(value string) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_NAME"))
	if err != nil {
		return err
	}
	newKey := datastore.IncompleteKey("Restaurant", nil)
	r := Restaurant{
		Name: value,
		CreatedAt: time.Now(),
	}
	if _, err := client.Put(ctx, newKey, &r); err != nil {
		return err
	}
	return nil
}

func list() ([]Restaurant, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_NAME"))
	if err != nil {
		return nil, err
	}
	var r []Restaurant
	q := datastore.NewQuery("Restaurant").Order("-created").Limit(5)
	keys, err := client.GetAll(ctx, q, &r)
	if err != nil {
		return nil, err
	}
	for i := 0; i<len(r); i++ {
		r[i].ID = keys[i].ID
	}
	return r, nil
}

func Lunch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		e := "Method Not Allowed."
		log.Println(e)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(e))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAllError: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	parsed, err := url.ParseQuery(string(b))
	if err != nil {
		log.Printf("ParseQueryError: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if parsed.Get("token") != os.Getenv("SLACK_TOKEN") {
		e := "Unauthorized Token."
		log.Printf(e)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(e))
		return
	}

	p := new(Parameter)
	p.parse(parsed.Get("text"))

	switch p.SubCommand {
	case "add":
		if err := add(p.Value); err != nil {
			log.Printf("DatastorePutError: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(p.Value))
	case "list":
		list, err := list()
		if err != nil {
			log.Printf("DatastoreGetAllError: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sprint(list)))
	default:
		e := "Invalid SubCommand."
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(p.Value))
}

func sprint(list []Restaurant) (s string) {
	for _, r := range list {
		s = s + fmt.Sprintf("[%v] %v\n", r.ID, r.Name)
	}
	return s
}

func (p *Parameter) parse(text string) {
	t := strings.TrimSpace(text)
	if len(t) < 1 {
		return
	}
	s := strings.SplitN(t, " ", 2)
	p.SubCommand = s[0]
	if len(s) == 1 {
		return
	}
	p.Value = s[1]
}