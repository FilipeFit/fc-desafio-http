package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
	err := GetCotacao()
	if err != nil {
		log.Fatal(err)
	}
}

func GetCotacao() error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var cotation Cotation
	err = json.Unmarshal(body, &cotation)
	if err != nil {
		return err
	}
	err = SaveCotationFile(cotation.USD.Bid)
	if err != nil {
		return err
	}

	return nil
}

func SaveCotationFile(bid string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	len, err := file.WriteString("Dólar: " + bid)

	if err != nil {
		return err
	}

	log.Println("Fle Created: ", file.Name())
	log.Println("Written bytes: ", len)

	return nil
}
