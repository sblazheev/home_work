package internalhttp

import (
	"net/http"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"            //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
)

type HTTPHandler struct {
	app    app.App
	logger common.LoggerInterface
}

func NewHandler(app app.App, logger common.LoggerInterface) *HTTPHandler {
	return &HTTPHandler{app, logger}
}

func (h *HTTPHandler) helloWorldHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hello-world"))
}
