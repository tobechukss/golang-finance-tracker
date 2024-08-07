package records

import (
	"finance-crud-app/internal/services/auth"
	"finance-crud-app/internal/types"
	"finance-crud-app/internal/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	recordStore types.RecordStore
	userStore   types.UserStore
}

func NewHandler(recordStore types.RecordStore, userStore types.UserStore) *Handler {
	return &Handler{recordStore: recordStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/record", auth.JWTAuthMiddleWare(h.handleGetRecord, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/record", auth.JWTAuthMiddleWare(h.handlePostRecord, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/record/{recordID}", auth.JWTAuthMiddleWare(h.handleGetRecordById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/record/{recordID", auth.JWTAuthMiddleWare(h.handleDeleteRecord, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetRecord(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIDFromContext(r.Context())

	records, err := h.recordStore.GetUserRecords(userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error retrievig records %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"data": records})
}

func (h *Handler) handlePostRecord(w http.ResponseWriter, r *http.Request) {
	var record types.PostRecordPayload
	if err := utils.ParseJSON(r, &record); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	if err := utils.Validate.Struct(record); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	new_record := types.Record{
		Description: record.Description,
		Amount:      record.Amount,
	}

	if record.Category != "" {
		new_record.Category = record.Category
	}

	userId := auth.GetUserIDFromContext(r.Context())

	log.Printf("userid %v", userId)

	created_record_id, err := h.recordStore.CreateUserRecord(userId, new_record)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("cannot create record %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]int{"record_id": created_record_id})
}

func (h *Handler) handleGetRecordById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordId, ok := vars["recordID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("record id not provided"))
		return
	}

	userId := auth.GetUserIDFromContext(r.Context())

	accessible := h.recordStore.CheckRecordBelongsToUser(userId, recordId)
	if !accessible {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error accessing record"))
		return
	}

	records, err := h.recordStore.GetRecordById(recordId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error retrievig records %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{"data": records})
}

func (h *Handler) handleDeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordId, ok := vars["recordID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("record id not provided"))
		return
	}
	userId := auth.GetUserIDFromContext(r.Context())

	err := h.recordStore.UserDeleteRecord(recordId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not remove record %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, nil)
}
