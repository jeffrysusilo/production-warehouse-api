package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

type Item struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Quantity    int    `json:"quantity"`
	Warehouse   string `json:"warehouse"`
	Description string `json:"description"`
}

type Material struct {
	ItemID       string `json:"item_id"`
	QuantityUsed int    `json:"quantity_used"`
}

type Production struct {
	ProductName      string     `json:"product_name"`
	Materials        []Material `json:"materials"`
	QuantityProduced int        `json:"quantity_produced"`
}

var createdItemID string
var createdProductionID string

func TestCreateItem(t *testing.T) {
	item := Item{
		Name:        "Besi",
		Category:    "Bahan Baku",
		Quantity:    1000,
		Warehouse:   "Gudang A",
		Description: "Material dasar",
	}
	body, _ := json.Marshal(item)
	resp, err := http.Post(baseURL+"/items", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != 201 {
		t.Fatalf("Gagal membuat item: %v", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Item Created:", string(respBody))
	var res map[string]interface{}
	json.Unmarshal(respBody, &res)
	createdItemID = res["id"].(string)
}

func TestGetItems(t *testing.T) {
	resp, err := http.Get(baseURL + "/items")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Gagal mengambil daftar item: %v", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Items:", string(respBody))
}

func TestCreateProduction(t *testing.T) {
	if createdItemID == "" {
		t.Skip("Item belum dibuat")
	}
	prod := Production{
		ProductName:      "Meja",
		QuantityProduced: 100,
		Materials: []Material{{
			ItemID:       createdItemID,
			QuantityUsed: 10,
		}},
	}
	body, _ := json.Marshal(prod)
	resp, err := http.Post(baseURL+"/productions", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != 202 {
		t.Fatalf("Gagal membuat produksi: %v", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Production Started:", string(respBody))
	var res map[string]interface{}
	json.Unmarshal(respBody, &res)
	createdProductionID = res["id"].(string)
}

func TestGetProductionByID(t *testing.T) {
	if createdProductionID == "" {
		t.Skip("Produksi belum dibuat")
	}
	resp, err := http.Get(baseURL + "/productions/" + createdProductionID)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Gagal mengambil produksi berdasarkan ID: %v", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Production Detail:", string(respBody))
}

func TestCancelProduction(t *testing.T) {
	if createdProductionID == "" {
		t.Skip("Produksi belum dibuat")
	}
	req, _ := http.NewRequest("POST", baseURL+"/productions/"+createdProductionID+"/cancel", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Gagal membatalkan produksi: %v", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Production Canceled:", string(respBody))
}

func TestGetLogs(t *testing.T) {
	if createdProductionID == "" {
		t.Skip("Produksi belum dibuat")
	}
	resp, err := http.Get(baseURL + "/productions/" + createdProductionID + "/logs")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Gagal mengambil log: %v", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Logs:", string(respBody))
}
