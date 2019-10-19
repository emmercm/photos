package sqlite3

import (
	"net/url"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // use the SQLite 3 driver
	"github.com/pkg/errors"
)

var (
	connection *gorm.DB
	once       sync.Once
)

func connectionString() (string, error) {
	host := os.Getenv("DB_HOST")
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	q := u.Query()

	user := os.Getenv("DB_USERNAME")
	if len(user) > 0 {
		q.Set("_auth", "")
		q.Set("_auth_user", user)
	}

	password := os.Getenv("DB_PASSWORD")
	if len(password) > 0 {
		q.Set("_auth", "")
		q.Set("_auth_pass", password)
	}

	u.RawQuery = q.Encode()

	return url.QueryEscape(u.String()), nil
}

// Connection returns a singleton DB connection
func Connection() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		var u string
		u, err = connectionString()
		if err != nil {
			return
		}

		var db *gorm.DB
		db, err = gorm.Open("sqlite3", u)
		if err != nil {
			return
		}

		connection = db.LogMode(false)
	})

	return connection, errors.Wrap(err, "database open failed")
}
