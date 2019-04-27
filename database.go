package gpt

import (
	"bytes"
	"fmt"
	"os"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
)

func mapJsonToDbCasing(json string) string {
	buf := bytes.NewBuffer([]byte{})
	for i, ch := range json {
		if i != 0 && unicode.IsUpper(ch) {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToLower(ch))
	}
	return buf.String()
}

func InitDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")))
	db.Mapper = reflectx.NewMapperTagFunc("json", mapJsonToDbCasing, mapJsonToDbCasing)
	if err != nil {
		return db, err
	}
	return db, nil
}
