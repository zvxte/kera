package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

func NewHabitsMux(habitStore store.HabitStore, logger *log.Logger) *http.ServeMux {
	h := &habitHandler{
		habitStore: habitStore,
		logger:     logger,
	}

	m := http.NewServeMux()
	m.HandleFunc("POST /{$}", makeHandlerFunc(h.create))
	m.HandleFunc("GET /{$}", makeHandlerFunc(h.getAll))
	m.HandleFunc("PATCH /{id}/title", makeHandlerFunc(h.patchTitle))
	m.HandleFunc("PATCH /{id}/description", makeHandlerFunc(h.patchDescription))
	m.HandleFunc("PATCH /{id}/end", makeHandlerFunc(h.end))
	return m
}

type habitHandler struct {
	habitStore store.HabitStore
	logger     *log.Logger
}

func (h *habitHandler) create(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(model.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	var in struct {
		Title       string          `json:"title"`
		Description string          `json:"description"`
		WeekDays    []model.WeekDay `json:"week_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	habit, err := model.NewHabit(
		in.Title, in.Description, in.WeekDays...,
	)
	if err != nil {
		return newJsonResponse(
			http.StatusUnprocessableEntity,
			newHandlerError(http.StatusUnprocessableEntity, err.Error()),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.habitStore.Create(ctx, habit, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}

func (h *habitHandler) getAll(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(model.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	habits, err := h.habitStore.GetAll(ctx, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	type out struct {
		ID          string            `json:"id"`
		Status      model.HabitStatus `json:"status"`
		Title       string            `json:"title"`
		Description string            `json:"description"`
		WeekDays    []uint            `json:"week_days"`
		StartDate   time.Time         `json:"start_date"`
		EndDate     time.Time         `json:"end_date"`
	}

	outs := make([]out, len(habits))

	for i, habit := range habits {
		weekDays := habit.TrackedWeekDays.WeekDays()

		// Cast []model.WeekDay ([]uint8) into []uint
		// to prevent json encoder from encoding it as base64 string
		weekDaysOut := make([]uint, len(weekDays))
		for i, d := range weekDays {
			weekDaysOut[i] = uint(d)
		}

		outs[i] = out{
			ID:          habit.ID.String(),
			Status:      habit.Status,
			Title:       habit.Title,
			Description: habit.Description,
			WeekDays:    weekDaysOut,
			StartDate:   habit.StartDate,
			EndDate:     habit.EndDate,
		}
	}

	return newJsonResponse(
		http.StatusOK,
		outs,
	)
}

func (h *habitHandler) patchTitle(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(model.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	habitID, err := model.ParseUUID(
		r.PathValue("id"),
	)
	if err != nil {
		return badRequestResponse
	}

	var in struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	if err := model.ValidateTitle(in.Title); err != nil {
		return newJsonResponse(
			http.StatusUnprocessableEntity,
			newHandlerError(http.StatusUnprocessableEntity, err.Error()),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.habitStore.UpdateTitle(ctx, habitID, in.Title, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}

func (h *habitHandler) patchDescription(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(model.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	id, err := model.ParseUUID(
		r.PathValue("id"),
	)
	if err != nil {
		return badRequestResponse
	}

	var in struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	if err := model.ValidateDescription(in.Description); err != nil {
		return newJsonResponse(
			http.StatusUnprocessableEntity,
			newHandlerError(http.StatusUnprocessableEntity, err.Error()),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.habitStore.UpdateDescription(ctx, id, in.Description, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}

func (h *habitHandler) end(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(model.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	id, err := model.ParseUUID(
		r.PathValue("id"),
	)
	if err != nil {
		return badRequestResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.habitStore.End(ctx, id, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}
