package httpserver

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Nexadis/Storage/internal/storage"
	"github.com/labstack/echo/v4"
)

func (hs *HTTPServer) MountHandlers() {
	loadShedding := NewLoadShedding(2, 10*time.Second)
	hs.Use(loadShedding)
	hs.GET(APIKV, hs.GetValue)
	hs.PUT(APIKV, hs.PutValue)
	hs.DELETE(APIKV, hs.DeleteValue)
}

func (hs *HTTPServer) GetValue(c echo.Context) error {
	key := c.Param("key")
	user := GetUser(c.Request())
	log.Printf("Got %s with k=%s", c.Request().Method, key)
	val, err := hs.db.Get(user, key)
	if err != nil {
		if errors.Is(err, storage.ErrorNoSuchKey) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, val)
}

func (hs *HTTPServer) PutValue(c echo.Context) error {
	key := c.Param("key")
	user := GetUser(c.Request())
	v, err := io.ReadAll(c.Request().Body)
	value := string(v)
	defer c.Request().Body.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	log.Printf("Got %s with k=%s v=%s", c.Request().Method, key, value)
	err = hs.db.Put(user, key, value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	hs.l.WritePut(user, key, value)

	return c.String(http.StatusCreated, value)
}

func (hs *HTTPServer) DeleteValue(c echo.Context) error {
	key := c.Param("key")
	user := GetUser(c.Request())
	log.Printf("Got %s with k=%s", c.Request().Method, key)
	err := hs.db.Delete(user, key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	hs.l.WriteDelete(user, key)

	return c.String(http.StatusOK, key)
}
