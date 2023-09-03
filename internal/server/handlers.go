package server

import (
	"errors"
	"io"
	"net/http"

	"github.com/Nexadis/Storage/internal/storage"
	"github.com/labstack/echo/v4"
)

func (hs *HTTPServer) MountHandlers() {
	hs.GET(APIKV, hs.GetValue)
	hs.PUT(APIKV, hs.PutValue)
	hs.DELETE(APIKV, hs.DeleteValue)
}

func (hs *HTTPServer) GetValue(c echo.Context) error {
	key := c.Param("key")
	val, err := hs.s.Get(key)
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
	value, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = hs.s.Put(key, string(value))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, string(value))
}

func (hs *HTTPServer) DeleteValue(c echo.Context) error {
	key := c.Param("key")
	err := hs.s.Delete(key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, key)
}
