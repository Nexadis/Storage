package httpserver

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

func NewLoadShedding(depth uint, timeAvailable time.Duration) func(echo.HandlerFunc) echo.HandlerFunc {
	m := sync.RWMutex{}
	CurrentDepth := 0
	ticker := time.NewTicker(timeAvailable)
	go func() {
		for range ticker.C {
			log.Printf("Refresh Timer")
			m.Lock()
			CurrentDepth = 0
			m.Unlock()
		}
	}()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			m.RLock()
			if CurrentDepth > int(depth) {
				m.RUnlock()
				log.Println("Load shedding engaged")
				return c.NoContent(http.StatusServiceUnavailable)
			}
			m.RUnlock()
			log.Printf("Load shedding: depth=%d\n", CurrentDepth)
			m.Lock()
			CurrentDepth += 1
			m.Unlock()
			next(c)
			return nil
		})
	}
}
