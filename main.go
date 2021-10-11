package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/pokemosha/gotasks/generator"
	"log"
	"net/http"
)

var db *sqlx.DB

type createRequest struct {
	FullUrl string `json:"furl"`
}

type createResponse struct {
	ShortUrl string `json:"shurl"`
}

type dbURL struct {
	shortURLDB string `db:"shortURL"`
	fullURLDB  string `db:"fullURL"`
}

func create(c echo.Context) error {
	u := new(createRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	g := generator.Generator{}
	short := ""

	//для возврата короткой
	err := db.Get(&short, "SELECT shortURL FROM URL WHERE fullURL = $1", u.FullUrl)
	if err != nil {
		continueFOR := true
		for continueFOR {
			short = g.ShortURL()
			rez, err := db.Exec("insert into URL (shortURL, fullURL) values ($1,$2) on conflict do nothing", short, u.FullUrl)
			if err != nil {
				return err
			}
			aff, err := rez.RowsAffected()
			if err != nil {
				return err
			}
			if aff != 0 {
				continueFOR = false
			}
		}
	}
	c.JSON(http.StatusOK, createResponse{ShortUrl: "http://localhost:1323/" + short})
	return nil
}

func main() {

	var err error
	db, err = sqlx.Connect("postgres", "postgres://ondvpcqesddkvv:23606e00d0d36dec38e4e025822b7b8d87366cacc80c825a095da9147306f22a@ec2-63-33-14-215.eu-west-1.compute.amazonaws.com:5432/d63ugo89dct77n")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/:shurl", func(c echo.Context) error {
		var full string
		err := db.Get(&full, "SELECT fullURL FROM URL WHERE shortURL = $1", c.Param("shurl"))
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusMovedPermanently, full)
	})
	e.POST("/", create)

	e.Logger.Fatal(e.Start(":1323"))
}
