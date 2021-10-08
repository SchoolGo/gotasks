package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var LinkMap = map[string]string{}

type createRequest struct {
	FullUrl string `json:"furl"`
}

type createResponse struct {
	ShortUrl string `json:"shurl"`
}

func create(c echo.Context) error {
	u := new(createRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	short := uuid.NewString()
	LinkMap[short] = u.FullUrl
	//TODO: сделать норм генератор ссылок
	c.JSON(http.StatusOK, createResponse{ShortUrl: "http://localhost:1323/" + short})
	return nil
}

func main() {
	e := echo.New()
	e.GET("/:shurl", func(c echo.Context) error {
		return c.String(http.StatusOK, LinkMap[c.Param("shurl")])
	})
	e.POST("/", create)

	e.Logger.Fatal(e.Start(":1323"))
}
