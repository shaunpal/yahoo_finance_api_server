package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

// Equity struct to format parse equity data
type Equity struct {
	Name                 string `json: "name"`
	Exchange             string `json: "exchange"`
	Value                string `json: "value"`
	Performance          string `json: "performance"`
	Status               string `json: "status"`
	ClosePrice           string `json: "close_price"`
	OpenPrice            string `json: "open_price"`
	Bid                  string `json: "bid"`
	Ask                  string `json: "ask"`
	DayRange             string `json: "day_range"`
	Week52Range          string `json: "52_week_range"`
	Volume               string `json: "volume"`
	AverageVolume        string `json: "average_volume"`
	MarketCapitalization string `json:"market_cap"`
	Beta                 string `json: "beta"`
	PERatio              string `json: "pe_ratio"`
	EPS                  string `json: "eps"`
	EarningsDate         string `json: "earnings_date"`
	YearlyTargetEstimate string `json: "yearly_target_estimate"`
}

// Index struct to format parse index data
type Index struct {
	Name          string `json:"name"`
	Exchange      string `json:"exchange"`
	Value         string `json:"value"`
	Performance   string `json:"performance"`
	Status        string `json:"status"`
	ClosePrice    string `json:"close_price"`
	OpenPrice     string `json:"open_price"`
	DayRange      string `json:"day_range"`
	Week52Range   string `json:"52_week_range"`
	Volume        string `json:"volume"`
	AverageVolume string `json:"average_volume"`
}

// ErrorHandler for formatting error
type ErrorHandler struct {
	ErrorType string `json:"errorType"`
	Message   string `json:"message"`
}

// GetEquity function to retrieve public traded company's equity stock data
func GetEquity(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var data []string
	url := fmt.Sprintf("https://sg.finance.yahoo.com/quote/%s", strings.ToUpper(vars["symbol"]))
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find quote-header-info div tag
	doc.Find("#quote-header-info").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("h1").Text()
		// w.Write([]byte(title))
		data = append(data, title)

		// For each item found, filter text
		s.Find("span").Each(func(i int, s *goquery.Selection) {
			if s.Text() != "Add to watchlist" {
				data = append(data, s.Text())
			}
		})

	})

	doc.Find("#quote-summary").Each(func(i int, s *goquery.Selection) {
		// For each item found, get span tag
		s.Find("td").Each(func(i int, s *goquery.Selection) {

			// get inner text of each td tag
			data = append(data, s.Text())
		})
	})

	if len(data) == 0 {
		errorMsg := fmt.Sprintf(`{ "errorType": "Not Found", "message": "Index Symbol (%s) may be invalid" }`, vars["symbol"])
		var errorhandler ErrorHandler
		err := json.Unmarshal([]byte(errorMsg), &errorhandler)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(errorMsg))
	} else {
		progress := make(chan *Equity)
		go formatEquityData(data, progress)
		result, _ := json.Marshal(<-progress)
		w.Write([]byte(result))
	}
}

// GetIndex function to retrieve index data e.g. ^STI, ^DJI
func GetIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var data []string
	url := fmt.Sprintf("https://sg.finance.yahoo.com/quote/^%s?p=^%s", strings.ToUpper(vars["symbol"]), strings.ToUpper(vars["symbol"]))
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load HTML body document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find quoto-header-info div tag
	doc.Find("#quote-header-info").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the inner text
		title := s.Find("h1").Text()
		// w.Write([]byte(title))
		data = append(data, title)
		s.Find("span").Each(func(i int, s *goquery.Selection) {

			// get inner text of each span tag, exclude string filter
			if s.Text() != "Add to watchlist" {
				data = append(data, s.Text())
			}
		})
	})

	doc.Find("#quote-summary").Each(func(i int, s *goquery.Selection) {
		// For each item found, get span tag
		s.Find("td").Each(func(i int, s *goquery.Selection) {

			// get inner text of each td tag
			data = append(data, s.Text())
		})
	})
	if len(data) == 0 {
		errorMsg := fmt.Sprintf(`{ "errorType": "Not Found", "message": "Index Symbol (%s) may be invalid" }`, vars["symbol"])
		var errorhandler ErrorHandler
		err := json.Unmarshal([]byte(errorMsg), &errorhandler)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(errorMsg))
	} else {
		progress := make(chan *Index)
		go formatIndexData(data, progress)
		result, _ := json.Marshal(<-progress)
		w.Write([]byte(result))
	}
}

func formatIndexData(data []string, progress chan *Index) {
	if len(data) != 0 && len(data) < 18 {
		raw := &Index{
			Name:          data[0],
			Exchange:      data[1],
			Value:         data[2],
			Performance:   data[3],
			Status:        data[4],
			ClosePrice:    data[6],
			OpenPrice:     data[8],
			Volume:        data[10],
			DayRange:      data[12],
			Week52Range:   data[14],
			AverageVolume: data[16],
		}
		progress <- raw
	}
}

func formatEquityData(data []string, progress chan *Equity) {
	if len(data) != 0 && len(data) < 33 {
		raw := &Equity{
			Name:                 data[0],
			Exchange:             data[1],
			Value:                data[2],
			Performance:          data[3],
			Status:               data[4],
			ClosePrice:           data[6],
			OpenPrice:            data[8],
			Bid:                  data[10],
			Ask:                  data[12],
			DayRange:             data[14],
			Week52Range:          data[16],
			Volume:               data[18],
			AverageVolume:        data[20],
			MarketCapitalization: data[22],
			Beta:                 data[24],
			PERatio:              data[26],
			EPS:                  data[28],
			EarningsDate:         data[30],
			YearlyTargetEstimate: data[32],
		}
		progress <- raw
	}
}
