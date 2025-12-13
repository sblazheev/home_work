package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"                                    //nolint:depguard
	"github.com/google/uuid"                                                    //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"        //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"     //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common/dto" //nolint:depguard
	httpSwagger "github.com/swaggo/http-swagger/v2"                             //nolint:depguard
)

type HTTPHandler struct {
	app    app.App
	logger common.LoggerInterface
	mux    *http.ServeMux
}

type JSONErrorResponse struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Detail  string  `json:"description,omitempty"`
} // @name JSONError .

func NewHandler(app app.App, logger common.LoggerInterface) *HTTPHandler {
	mux := http.NewServeMux()

	handler := &HTTPHandler{app, logger, mux}

	mux.HandleFunc("/ping", handler.pingHandler)
	mux.HandleFunc("POST /event/update", handler.createEventHandler)
	mux.HandleFunc("PUT /event/update", handler.updateEventHandler)
	mux.HandleFunc("GET /event/status/{uuid}", handler.getEventNotificationStatusHandler)
	mux.HandleFunc("GET /event/{uuid}", handler.getEventHandler)
	mux.HandleFunc("DELETE /event/{uuid}", handler.deleteEventHandler)
	mux.HandleFunc("GET /event/list", handler.listEventHandler)
	mux.HandleFunc("GET /event/listday", handler.listEventByDayHandler)
	mux.HandleFunc("GET /event/listweek", handler.listEventByWeekHandler)
	mux.HandleFunc("GET /event/listmonth", handler.listEventByMonthHandler)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
	return handler
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func JSONError(httpcode int, code string, messageError string, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if httpcode == 503 {
		w.Header().Set("Retry-After", "600")
	}
	w.WriteHeader(httpcode)
	json.NewEncoder(w).Encode(
		JSONErrorResponse{
			Code:    &code,
			Message: &messageError,
			Detail:  err.Error(),
		},
	)
}

// @Summary      Создать событие
// @Description  Создать событие
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        data body dto.Event true  "Создание события"
// @Success      201  {object} dto.Event
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /event/update [post] .
func (h *HTTPHandler) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var dtoEvent *dto.Event
	if err := json.NewDecoder(r.Body).Decode(&dtoEvent); err != nil {
		h.logger.Debug("createEventHandler-Invalid request body", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid request body", err, w)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoEvent)
	if err != nil {
		h.logger.Debug("createEventHandler-Invalid format Event", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid format Event", err, w)
		return
	}
	dtoEvent, err = h.app.CreateEvent(dtoEvent)
	if err != nil {
		if errors.Is(err, common.ErrEventConflictOverlap) {
			h.logger.Debug("createEventHandler-ConflictOverlap", "err", err)
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "ConflictOverlap", err, w)
		} else {
			h.logger.Error("createEventHandler", "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("createEventHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(dtoEventJSON)
}

// @Summary      Измененить событие
// @Description  Измененить событие
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        data body dto.Event true  "Изменение события"
// @Success      200  {object} dto.Event
// @Router       /event/update [put] .
func (h *HTTPHandler) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var dtoEvent dto.Event
	if err := json.NewDecoder(r.Body).Decode(&dtoEvent); err != nil {
		h.logger.Debug("updateEventHandler-Invalid request body", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid request body", err, w)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoEvent)
	if err != nil {
		h.logger.Debug("updateEventHandler-Invalid format Event", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid format Event", err, w)
		return
	}
	err = h.app.UpdateEvent(&dtoEvent)
	if err != nil {
		if errors.Is(err, common.ErrEventConflictOverlap) {
			h.logger.Debug("updateEventHandler-ConflictOverlap", "err", err)
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "ConflictOverlap", err, w)
		} else {
			h.logger.Error("updateEventHandler", "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("updateEventHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.Write(dtoEventJSON)
}

// @Summary      Получить событие
// @Description  Получить событие
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param    	 uuid    path  string true  "UUID события"
// @Success      200  {object} dto.Event
// @Router       /event/{uuid} [get] .
func (h *HTTPHandler) getEventHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uuid")
	err := uuid.Validate(uid)
	if err != nil {
		h.logger.Debug("getEventHandler-Uuid invalid format", "uuid", uid, "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Uuid invalid format", err, w)
		return
	}
	dtoEvent, err := h.app.GetEvent(uid)
	if err != nil {
		if errors.Is(err, common.ErrEventNotFound) {
			h.logger.Debug("getEventHandler-NotFound", "uuid", uid, "err", err)
			JSONError(http.StatusNotFound, strconv.Itoa(http.StatusNotFound), "NotFound", err, w)
		} else {
			h.logger.Error("getEventHandler", "uuid", uid, "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("getEventHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.Write(dtoEventJSON)
}

// @Summary      Получить событие
// @Description  Получить событие
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param    	 uuid    path  string true  "UUID события"
// @Success      200  {object} common.NotificationStatus
// @Failure		 404
// @Router       /event/status/{uuid} [get] .
func (h *HTTPHandler) getEventNotificationStatusHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uuid")
	err := uuid.Validate(uid)
	if err != nil {
		h.logger.Debug("getEventNotificationStatusHandler-Uuid invalid format", "uuid", uid, "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Uuid invalid format", err, w)
		return
	}
	status, err := h.app.GetNotificationStatus(context.Background(), uid)
	if err != nil {
		if errors.Is(err, common.ErrEventNotFound) {
			h.logger.Debug("getEventNotificationStatusHandler-NotFound", "uuid", uid, "err", err)
			JSONError(http.StatusNotFound, strconv.Itoa(http.StatusNotFound), "NotFound", err, w)
		} else {
			h.logger.Error("getEventNotificationStatusHandler", "uuid", uid, "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	eventJSON, err := json.Marshal(status)
	if err != nil {
		h.logger.Error("getEventNotificationStatusHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if status == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write(eventJSON)
}

// @Summary      Удалить событие
// @Description  Удалить событие
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param    	 uuid    path  string true  "Удалить событие"
// @Success      204
// @Router       /event/{uuid} [delete] .
func (h *HTTPHandler) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uuid")
	err := uuid.Validate(uid)
	if err != nil {
		h.logger.Debug("deleteEventHandler-Uuid invalid format", "uuid", uid, "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "NotFound", err, w)
		return
	}
	err = h.app.DeleteEvent(uid)
	if err != nil {
		h.logger.Error("deleteEventHandler", "uuid", uid, "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Список событий
// @Description  Список событий
// @Tags         Event
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Event
// @Router       /event/list [get] .
func (h *HTTPHandler) listEventHandler(w http.ResponseWriter, _ *http.Request) {
	dtoEvent, err := h.app.ListEvent()
	if err != nil {
		h.logger.Error("listEventHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("listEventHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.Write(dtoEventJSON)
}

// @Summary      Список событий для пользователя за сутки
// @Description  Список событий для пользователя за сутки
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param    	 user    query  string true  "ID пользователя"
// @Param    	 date    query  string true  "Дата YYYY-MM-DD"
// @Success      200  {array} dto.Event
// @Router       /event/listday [get] .
func (h *HTTPHandler) listEventByDayHandler(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	dtoEvent, err := h.app.ListEventUserByDay(r.URL.Query().Get("user"), r.URL.Query().Get("date"))
	if err != nil {
		if errors.Is(err, common.ErrQueryRequest) {
			h.logger.Debug("listEventByDayHandler", "err", err)
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
				"Bad Request", err, w)
		} else {
			h.logger.Error("listEventByDayHandler", "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("listEventByDayHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.Write(dtoEventJSON)
}

// @Summary      Список событий для пользователя за неделю
// @Description  Список событий для пользователя за неделю
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param    	 user    query  string true  "ID пользователя"
// @Param    	 date    query  string true  "Дата первого дня недели в формате YYYY-MM-DD"
// @Success      200  {array} dto.Event
// @Router       /event/listweek [get] .
func (h *HTTPHandler) listEventByWeekHandler(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	dtoEvent, err := h.app.ListEventUserByWeek(r.URL.Query().Get("user"), r.URL.Query().Get("date"))
	if err != nil {
		if errors.Is(err, common.ErrQueryRequest) {
			h.logger.Debug("listEventByWeekHandler", "err", err)
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
				"Bad Request", err, w)
		} else {
			h.logger.Error("listEventByWeekHandler", "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("listEventByWeekHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.Write(dtoEventJSON)
}

// @Summary      Список событий для пользователя за месяц
// @Description  Список событий для пользователя за месяц
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param    	 user    query  string true  "ID пользователя"
// @Param    	 date    query  string true  "Дата первого дня месяц в формате YYYY-MM-DD"
// @Success      200  {array} dto.Event
// @Router       /event/listmonth [get] .
func (h *HTTPHandler) listEventByMonthHandler(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	dtoEvent, err := h.app.ListEventUserByMonth(r.URL.Query().Get("user"), r.URL.Query().Get("date"))
	if err != nil {
		if errors.Is(err, common.ErrQueryRequest) {
			h.logger.Debug("listEventByMonthHandler", "err", err)
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
				"Bad Request", err, w)
		} else {
			h.logger.Error("listEventByMonthHandler", "err", err)
			JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
				"Service Unavailable", common.ErrServiceUnavailable, w)
		}
		return
	}
	dtoEventJSON, err := json.Marshal(dtoEvent)
	if err != nil {
		h.logger.Error("listEventByMonthHandler-Json marshal", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.Write(dtoEventJSON)
}
