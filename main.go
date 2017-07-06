package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type (
	store interface {
		Create(p *post) error
		Delete(ID int) error
		FindByID(ID int) (post, error)
	}

	app struct {
		DB store
	}

	post struct {
		ID      int    `db:"id"`
		Title   string `db:"title"`
		Slug    string `db:"slug"`
		Content string `db:"content"`
	}
)

func main() {
	db := mustNewDB()

	a := &app{db}

	http.HandleFunc("/posts/view/", viewPost(a))
	http.HandleFunc("/posts/delete/", deletePost(a))
	http.HandleFunc("/posts/create", createPost(a))
	http.ListenAndServe(":3000", nil)
}

type db struct {
	*sqlx.DB
}

func (d *db) Create(p *post) error {
	//for simplicity sake, we aren't checking for the existence of a simlar post
	stmt, err := d.Preparex("INSERT INTO posts(title, slug, content) VALUES(?,?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.MustExec(p.Title, p.Slug, p.Content).RowsAffected()

	return err
}

func (d *db) Delete(ID int) error {
	stmt, err := d.Preparex("DELETE FROM posts WHERE id=?")

	if err != nil {
		return err
	}

	_, err = stmt.MustExec(ID).RowsAffected()

	return err
}

func (d *db) FindByID(ID int) (post, error) {
	var p post
	var err error

	stmt, err := d.Preparex("SELECT * FROM posts WHERE id=?")

	if err != nil {
		return p, err
	}

	err = stmt.QueryRowx(ID).StructScan(&p)
	return p, err
}

func mustNewDB() *db {
	con := sqlx.MustConnect("sqlite3", "demo.sqlite")
	return &db{con}
}

func createPost(a *app) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		p := &post{}

		if err := parseJSON(r.Body, p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := a.DB.Create(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Your post was successfully created"))
	}
}

func deletePost(a *app) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cleanID, err := strconv.Atoi(r.URL.Path[14:])

		if err != nil || cleanID == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if err := a.DB.Delete(cleanID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(
			fmt.Sprintf("The post with ID %d was successfully deleted", cleanID)))
	}
}

func viewPost(a *app) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cleanID, err := strconv.Atoi(r.URL.Path[12:])

		if err != nil || cleanID == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		p, err := a.DB.FindByID(cleanID)

		if err != nil {
			http.Error(w, fmt.Sprintf("The post with the ID, %d does not exist", cleanID), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(p)
	}
}

func parseJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
