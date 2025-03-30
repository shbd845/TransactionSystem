package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"internalTransferSystem/models"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes() {
	http.HandleFunc("/accounts", h.CreateAccountHandler)
	http.HandleFunc("/accounts/", h.GetAccountHandler)
	http.HandleFunc("/transactions", h.CreateTransactionHandler)
}

func (h *Handler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID      int    `json:"account_id"`
		InitialBalance string `json:"initial_balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := models.CreateAccount(h.db, req.AccountID, req.InitialBalance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "account ID not provided", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	account, err := models.GetAccount(h.db, accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (h *Handler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SourceAccountID      int    `json:"source_account_id"`
		DestinationAccountID int    `json:"destination_account_id"`
		Amount               string `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := models.ProcessTransaction(h.db, req.SourceAccountID, req.DestinationAccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
