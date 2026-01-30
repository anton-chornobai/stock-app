package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"todo-app/model"
)

func GetStocks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome"))
}

func GetSingleStock(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.Write([]byte("No id provided"))
		return
	}
	var stockReq struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&stockReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erorr parsing %v", err), http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "./stocks.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB: %v", err), http.StatusInternalServerError)
		return
	}

	var stock model.Stock
	defer db.Close()
	err = db.QueryRow(
		`SELECT id, name, symbol, created_at FROM stocks WHERE id = ?`, stockReq.ID,
	).Scan(&stock.ID, &stock.Name, &stock.Symbol, &stock.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("DB query error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func PostStock(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./stocks.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var payload model.Stock
	e := json.NewDecoder(r.Body).Decode(&payload)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received payload: %+v\n", payload)
	fmt.Fprintf(w, "Payload received successfully: Name=%s, ID=%v", payload.Name, payload.ID)
}
