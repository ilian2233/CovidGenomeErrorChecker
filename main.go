package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DNA string

func loadDNA() DNA {
	content, err := ioutil.ReadFile("covid_dna")
	if err != nil {
		log.Println(err)
	}
	return DNA(content)
}

func (dna DNA)checkCovidDNA(w http.ResponseWriter, r *http.Request){
	input, err := ioutil.ReadAll(r.Body)
	if err != nil{
		_,err = fmt.Fprintf(w, "Failed to read body")
		if err != nil{
			log.Println(err)
		}
		return
	}

	if len(input) != len(dna){
		_,err = fmt.Fprintf(w, "The DNAs are with different length")
		if err != nil{
			log.Println(err)
		}
		return
	}

	var diffs []int

	for i,v := range input{
		if v != dna[i]{
			diffs = append(diffs, i)
		}
	}

	_,err = fmt.Fprintf(w, "%v",diffs)
	if err != nil{
		log.Println(err)
	}
}

func main() {
	dna := loadDNA()

	http.HandleFunc("/", dna.checkCovidDNA)
	log.Fatal(http.ListenAndServe(":8080", nil))
}