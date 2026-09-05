package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"true_shop/internal/httpx"
	"true_shop/internal/middleware"
)

type listing struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	// request scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx, `SELECT id,title,description,price,city,created_at
			FROM listings 
			ORDER BY created_at DESC
			LIMIT 100`)
	if err != nil {
		lh.logger.Error("listings query error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	listings := []listing{}
	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.Id, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("Listing fetch error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		lh.logger.Error("rows error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	lh.logger.Error("Successfully send listings", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listings)

}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	id := r.PathValue("id")
	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("delete failed", "listing_id:", id, "request_id", requestId, "err", err)
		// http.Error(w, "Internal error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, string(httpx.CodeInternalError), "something went wrong")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	var req listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("json decoding failed", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid_body", string(httpx.CodeMalformedJSON))
		return
	}
	var id string
	row := lh.db.QueryRowContext(
		ctx,
		`INSERT INTO listings (title, description, price, city) VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Title, req.Description, req.Price, req.City)
	if err := row.Scan(&id); err != nil {
		lh.logger.Error("failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", string(httpx.CodeInternalError))
		return
	}
	lh.logger.Info("listing created", "request_id", requestId, "listing_id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}
