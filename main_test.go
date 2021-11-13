package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSuccess(t *testing.T){
	dna := loadDNA()

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(dna)))
	w := httptest.NewRecorder()
	dna.checkCovidDNA(w, req)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected error to be nil got %v", err)
	}

	if string(data) != "[]"  {
		t.Fatalf("expected empty array, got: %v", string(data))
	}
}

func TestChangedDNALength(t *testing.T){
	dna := loadDNA()

	newDNA := dna[10:]

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(newDNA)))
	w := httptest.NewRecorder()
	dna.checkCovidDNA(w, req)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected error to be nil got %v", err)
	}

	if string(data) != "The DNAs are with different length"{
		t.Fatalf("expected 'The DNAs are with different length', recieved:%v", string(data))
	}
}

func TestChangedDNA(t *testing.T){
	dna := loadDNA()

	newDNA := strings.Replace(string(dna), "T", "G", 2)

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(newDNA)))
	w := httptest.NewRecorder()
	dna.checkCovidDNA(w, req)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected error to be nil got %v", err)
	}

	if string(data) != "[30 46]"{
		t.Fatalf("differances doens't look right, got: %v", string(data))
	}
}

func TestStressFunc(t *testing.T){
	for i:=0;i<100000; i++ {
		go TestSuccess(t) //nolint:govet
	}
}

func TestServer(t *testing.T){
	dna := loadDNA()
	url := "http://localhost:8080"
	payload := strings.NewReader(string(dna))
	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		t.Fatal(err)
		return
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
		return
	}

	if string(body) != "[]"  {
		t.Fatalf("expected empty array, got: %v", string(body))
	}
}

func TestStressServer(t *testing.T){
	for i:=0;i<100000; i++ {
		go TestServer(t) //nolint:govet
	}
}