package postgres

import (
	e "adviserbot/internal/lib/err"
	"adviserbot/internal/storage"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(conn string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := OpenDB(conn)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	query := `
	CREATE TABLE IF NOT EXISTS url(
		id serial primary key, 
		url text,
		username text
	)	
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func OpenDB(conn string) (*sql.DB, error) {
	const op = "storage.postgres.OpenDB"
	count := 0
	for {
		count++
		db, err := sql.Open("pgx", conn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		if count > 8 {
			return nil, e.Wrap(op, err)
		}
		log.Println("backing off for 2 seconds")
		time.Sleep(time.Second * 2)
	}
}

func (s *Storage) Save(p *storage.Page) error {
	const op = "storage.postgres.Save"
	query := "INSERT INTO url(url, username) VALUES($1, $2)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return e.Wrap(op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.URL, p.UserName)
	if err != nil {
		return e.Wrap(op, err)
	}
	return nil
}

func (s *Storage) PickRandom(username string) (*storage.Page, error) {
	const op = "storage.postgres.PickRandom"
	query := "SELECT * FROM url WHERE username = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	var urls []storage.Page
	for rows.Next() {
		var page storage.Page
		err = rows.Scan(&page.ID, page.URL, page.UserName)
		if err != nil {
			return nil, e.Wrap(op, err)
		}
		urls = append(urls, page)
	}
	l := len(urls)
	if l == 0 {
		return nil, e.Wrap(op, storage.ErrNotFound)
	}
	ind := rand.Intn(l - 1)
	return &urls[ind], nil
}

func (s *Storage) Remove(p *storage.Page) error {
	const op = "storage.postgres.Remove"
	query := "DELETE FROM url WHERE url = $1 and username = $2"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return e.Wrap(op, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.URL, p.UserName)
	if err != nil {
		return e.Wrap(op, err)
	}
	return nil
}

func (s *Storage) IsExists(p *storage.Page) (bool, error) {
	const op = "storage.postgres.IsExists"
	query := "SELECT * FROM url WHERE url = $1 and username = $2"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, e.Wrap(op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(p.URL, p.UserName)
	var page storage.Page
	err = row.Scan(&page.ID, &page.URL, &page.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, e.Wrap(op, err)
	}
	return true, nil
}
