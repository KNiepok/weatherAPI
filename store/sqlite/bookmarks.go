package sqlite

import (
	"database/sql"
	"github.com/kniepok/weatherAPI"
)

type BookmarkStorage struct {
	db *sql.DB
}

func NewBookmarkStorage(db *sql.DB) (*BookmarkStorage, error) {
	bs := &BookmarkStorage{
		db: db,
	}

	if err := bs.migrate(); err != nil {
		return nil, err
	}

	return bs, nil
}

func (bs *BookmarkStorage) Close() error {
	return bs.db.Close()
}

func (bs *BookmarkStorage) migrate() error {
	_, err := bs.db.Exec("CREATE TABLE IF NOT EXISTS locations (id integer not null primary key, name TEXT);")
	return err
}

func (bs *BookmarkStorage) GetBookmarks() ([]*weatherapi.Location, error) {
	rows, err := bs.db.Query("select name FROM locations ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	locations := make([]*weatherapi.Location, 0)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		locations = append(locations, &weatherapi.Location{Name: name})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (bs *BookmarkStorage) AddBookmark(l *weatherapi.Location) error {
	stmt, err := bs.db.Prepare("INSERT INTO locations(name) VALUES(?)")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(l.Name); err != nil {
		return err
	}
	return nil
}
