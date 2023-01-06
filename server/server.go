package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotation struct {
	USD struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USD"`
}

func main() {
	http.HandleFunc("/cotacao", CotationHandler)
	http.ListenAndServe(":8080", nil)
}

func CotationHandler(w http.ResponseWriter, r *http.Request) {
	cotation, err := GetDollarCotation()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotation)
}

func GetDollarCotation() (*Cotation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/all/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotation Cotation
	err = json.Unmarshal(body, &cotation)
	if err != nil {
		return nil, err
	}

	err = PersistCotation(&cotation)
	if err != nil {
		return nil, err
	}

	return &cotation, nil
}

func PersistCotation(cotation *Cotation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db, err := sql.Open("sqlite3", "cotation.db")
	if err != nil {
		return err
	}
	defer db.Close()
	const createTable string = "CREATE TABLE IF NOT EXISTS cotation (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,code TEXT,codein TEXT,name TEXT,high TEXT,low TEXT,varBid TEXT,pctChange TEXT,bid TEXT,ask TEXT,timestamp TEXT,create_date TEXT);"
	db.Exec(createTable)

	stm, err := db.PrepareContext(ctx, "INSERT INTO cotation (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(cotation.USD.Code, cotation.USD.Codein,
		cotation.USD.Name, cotation.USD.High,
		cotation.USD.Low, cotation.USD.VarBid,
		cotation.USD.PctChange, cotation.USD.Bid,
		cotation.USD.Ask, cotation.USD.Timestamp,
		cotation.USD.CreateDate)

	if err != nil {
		return err
	}

	return nil
}
