package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/model/habit"
	"github.com/zvxte/kera/model/uuid"
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
	m.HandleFunc("DELETE /{id}", makeHandlerFunc(h.delete))
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
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	var in struct {
		Title       string          `json:"title"`
		Description string          `json:"description"`
		WeekDays    []habit.WeekDay `json:"week_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	habit, err := habit.New(
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
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
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
		ID          string       `json:"id"`
		Status      habit.Status `json:"status"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		WeekDays    []uint       `json:"week_days"`
		StartDate   time.Time    `json:"start_date"`
		EndDate     time.Time    `json:"end_date"`
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
			StartDate:   time.Time(habit.StartDate),
			EndDate:     time.Time(habit.EndDate),
		}
	}

	return newJsonResponse(
		http.StatusOK,
		outs,
	)
}

func (h *habitHandler) delete(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	habitID, err := uuid.Parse(
		r.PathValue("id"),
	)
	if err != nil {
		return badRequestResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.habitStore.Delete(ctx, habitID, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}

func (h *habitHandler) patchTitle(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	habitID, err := uuid.Parse(
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

	if err := habit.ValidateTitle(in.Title); err != nil {
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
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	id, err := uuid.Parse(
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

	if err := habit.ValidateDescription(in.Description); err != nil {
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
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	id, err := uuid.Parse(
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