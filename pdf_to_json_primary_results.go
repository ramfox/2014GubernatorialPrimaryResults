/* Quick script to download and parse the  NYS Gubernatorial Primary Election results pdf */
/* primarily uses the github.com/ledongthuc/pdf package to parse pdfs */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

func main() {
	res, err := http.Get("https://www.elections.ny.gov/NYSBOE/elections/2014/Primary/2014StateLocalPrimaryElectionResults.pdf")
	if err != nil {
		fmt.Println("error getting pdf from url")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading data from body")
		return
	}

	res.Body.Close()

	r, err := pdf.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		fmt.Println("error creating pdf reader from bytes reader")
		return
	}

	page := r.Page(1)
	if page.V.IsNull() {
		fmt.Println("error: first page of pdf has no data")
		return
	}

	texts := page.Content().Text
	fullText := ""

	for _, text := range texts {
		fullText += text.S
	}

	lines := strings.Split(fullText, "\n")

	// meta := lines[0:10]
	votes := lines[10:]
	data := []interface{}{}

	for i := 0; i < len(votes)-1; i += 4 {
		if votes[i] == "Total NYC" || votes[i] == "Total Outside NYC" || votes[i] == "STATEWIDE TOTAL" {
			continue
		}
		row := []interface{}{votes[i]}

		intVotes, err := strToInt(votes[i+1])
		if err != nil {
			fmt.Printf("error converting string votes to int: %s\n", err)
			return
		}
		row = append(row, intVotes)

		intVotes, err = strToInt(votes[i+2])
		if err != nil {
			fmt.Printf("error converting string votes to int: %s\n", err)
			return
		}
		row = append(row, intVotes)

		intVotes, err = strToInt(votes[i+3])
		if err != nil {
			fmt.Printf("error converting string votes to int: %s\n", err)
			return
		}
		row = append(row, intVotes)

		data = append(data, row)
	}
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error marshaling data to json: %s\n", err)
		return
	}
	err = ioutil.WriteFile("2014_gubernatorial_democrat_primary_results.json", json, 0644)
	if err != nil {
		fmt.Printf("error writing json to file: %s\n", err)
		return
	}
}

func strToInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, ",", "", -1)
	return strconv.Atoi(s)
}
