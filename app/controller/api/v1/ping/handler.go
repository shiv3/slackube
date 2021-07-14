package ping

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Handler interface {
		GetPing(c echo.Context) error
	}

	handlerImpl struct {
		pingMessage string
	}
)

func NewHandlerImpl(pingMessage string) *handlerImpl {
	return &handlerImpl{pingMessage: pingMessage}
}

type PingResponce struct {
	Status string
}

func (h handlerImpl) GetPing(c echo.Context) error {
	err := c.JSON(http.StatusOK, PingResponce{Status: h.pingMessage})
	if err != nil {
		return c.Blob(http.StatusInternalServerError, "", nil)
	}
	return nil
}
