package sqlite

import (
	"database/sql"
	"github.com/kniepok/weatherAPI"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" //imports sqlite driver
)

func TestBookmarkStorage(t *testing.T) {

	l := &weatherapi.Location{Name: "Washington,DC"}

	f, err := ioutil.TempFile(".", "")
	if err != nil {
		t.Fatalf("could not create file for db: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}

	storage, err := NewBookmarkStorage(db)
	if err != nil {
		t.Fatalf("could not create storage: %v", err)
	}
	defer storage.Close()

	err = storage.AddBookmark(l)
	if err != nil {
		t.Fatalf("could not add location to storage: %v", err)
	}

	locations, err := storage.GetBookmarks()
	if err != nil {
		t.Fatalf("could not get locations: %v", err)
	}

	if l.Name != locations[0].Name {
		t.Fatalf("wanted %v , got %v", l, locations[0])
	}
}
