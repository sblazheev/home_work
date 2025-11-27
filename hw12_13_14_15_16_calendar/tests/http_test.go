//go:build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	internalgrpc "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/grpc"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common/dto"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb" //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc" //nolint:depguard
)

const (
	httpURL = "http://calendar:8080"
	grpcURL = "calendarGrpc:8081"
)

func TestCalendar(t *testing.T) {
	event := common.NewEvent("", "Test", time.Now(), 2, "Test", 0, 60)

	event1 := common.NewEvent("", "Test 1", event.DateTime.Add(time.Second*-1), 2, "Test 1", 0, 60)
	event2 := common.NewEvent("", "Test 2", event.DateTime.Add(time.Second*10), 1, "Test 2", 0, 8)
	event3 := common.NewEvent("", "Test 3", event.DateTime.Add(time.Second*100), 1, "Test 3", 0, 8)

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	grpcClient := pb.NewEventsServiceClient(grpcConn)

	t.Run("Create event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		dtoEvent, err := storage.MapperEventToDtoEvent(event)
		body, _ := json.Marshal(dtoEvent)

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/event/update", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var res dto.Event
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		assert.NotEmpty(t, res.ID)
		event.ID = res.ID
	})

	t.Run("Create event2", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		dtoEvent, err := storage.MapperEventToDtoEvent(event2)
		body, _ := json.Marshal(dtoEvent)

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/event/update", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var res dto.Event
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		require.NotEmpty(t, res.ID)

		event2.ID = res.ID
	})

	t.Run("Create event1 overlapping", func(t *testing.T) {
		//t.Skip()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		dtoEvent, err := storage.MapperEventToDtoEvent(event1)
		body, _ := json.Marshal(dtoEvent)

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/event/update", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		event.Title = "TitleTitle"
		dtoEvent, err := storage.MapperEventToDtoEvent(event)

		body, _ := json.Marshal(dtoEvent)

		httpReq, err := http.NewRequestWithContext(ctx, "PUT", httpURL+"/event/update", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var res dto.Event
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		assert.Equal(t, *dtoEvent, res)
	})

	t.Run("List event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(ctx, "GET", httpURL+"/event/list", nil)
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var res []dto.Event
		err = json.NewDecoder(resp.Body).Decode(&res)
		require.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})

	t.Run("Delete event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		{
			httpReq, err := http.NewRequestWithContext(ctx, "DELETE", httpURL+"/event/"+event.ID.(string), nil)
			require.NoError(t, err)
			httpReq.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(httpReq)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusNoContent, resp.StatusCode)
		}
		{
			httpReq, err := http.NewRequestWithContext(ctx, "GET", httpURL+"/event/list", nil)
			require.NoError(t, err)
			httpReq.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(httpReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var res []dto.Event
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, 1, len(res))
		}
	})

	t.Run("Notifycation send", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		httpReq, err := http.NewRequestWithContext(ctx, "GET", httpURL+"/event/status/"+event2.ID.(string), nil)
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		for {
			select {
			case <-ticker.C:
				resp, err := client.Do(httpReq)
				require.NoError(t, err)
				defer resp.Body.Close()
				if http.StatusOK == resp.StatusCode {
					var res common.NotificationStatus
					err = json.NewDecoder(resp.Body).Decode(&res)
					if res.Status == common.StatusNotifyDelivered {
						require.Equal(t, common.StatusNotifyDelivered, res.Status)
						return
					}
				}
			}
		}

	})

	t.Run("Create event grpc", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		pbEvent := pb.Event{
			ID:          event3.ID.(string),
			Title:       event3.Title,
			Description: event3.Description,
			DateTime:    uint64(event3.DateTime.Unix()),
			Duration:    event3.Duration,
			UserID:      strconv.Itoa(event3.User),
			NotifyTime:  event3.NotifyTime,
		}
		pbEventCreate, err := grpcClient.CreateEvent(ctx, &pbEvent)
		require.NoError(t, err)
		pbEvent.ID = pbEventCreate.ID
		event3.ID = pbEventCreate.ID
		dtoPbEventCreate, _ := internalgrpc.MapperEventToDtoEvent(pbEventCreate)
		dtoPbEvent, _ := internalgrpc.MapperEventToDtoEvent(&pbEvent)
		require.Equal(t, dtoPbEvent, dtoPbEventCreate)
	})

	t.Run("Update event grpc", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		event3.Title = "Title 3 Update"
		pbEvent := pb.Event{
			ID:          event3.ID.(string),
			Title:       event3.Title,
			Description: event3.Description,
			DateTime:    uint64(event3.DateTime.Unix()),
			Duration:    event3.Duration,
			UserID:      strconv.Itoa(event3.User),
			NotifyTime:  event3.NotifyTime,
		}
		pbEventCreate, err := grpcClient.UpdateEvent(ctx, &pbEvent)
		require.NoError(t, err)
		dtoPbEventCreate, _ := internalgrpc.MapperEventToDtoEvent(pbEventCreate)
		dtoPbEvent, _ := internalgrpc.MapperEventToDtoEvent(&pbEvent)
		require.Equal(t, dtoPbEvent, dtoPbEventCreate)
	})

	t.Run("Get event grpc", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		id := &pb.ID{
			UUID: event3.ID.(string),
		}
		pbEventGet, err := grpcClient.GetEvent(ctx, id)
		require.NoError(t, err)

		pbEvent := pb.Event{
			ID:          event3.ID.(string),
			Title:       event3.Title,
			Description: event3.Description,
			DateTime:    uint64(event3.DateTime.Unix()),
			Duration:    event3.Duration,
			UserID:      strconv.Itoa(event3.User),
			NotifyTime:  event3.NotifyTime,
		}
		dtoPbEventGet, _ := internalgrpc.MapperEventToDtoEvent(pbEventGet)
		dtoPbEvent, _ := internalgrpc.MapperEventToDtoEvent(&pbEvent)
		require.Equal(t, dtoPbEvent, dtoPbEventGet)
	})

	t.Run("Delete event grpc", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		id := &pb.ID{
			UUID: event3.ID.(string),
		}
		response, err := grpcClient.DeleteEvent(ctx, id)
		require.NoError(t, err)
		require.Equal(t, 204, int(response.GetStatus()))
	})

	t.Run("List event grpc", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		events, err := grpcClient.ListEvent(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, 1, len(events.Events))
	})
}
