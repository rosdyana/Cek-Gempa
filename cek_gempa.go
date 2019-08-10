package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type InfoGempa struct {
	XMLName   xml.Name `xml:"Infogempa"`
	InfoGempa []Gempa  `xml:"Gempa"`
}

type Gempa struct {
	XMLName    xml.Name `xml:"Gempa"`
	Tanggal    string   `xml:"Tanggal"`
	Posisi     string   `xml:"Posisi"`
	Kedalaman  string   `xml:"Kedalaman"`
	Point      Points   `xml:"point"`
	Magnitude  string   `xml:"Magnitude"`
	Keterangan string   `xml:"Keterangan"`
	Dirasakan  string   `xml:"Dirasakan"`
}

type Points struct {
	XMLName     xml.Name `xml:"point"`
	Coordinates string   `xml:"coordinates"`
}

func Usage() {
	// custom usage (help) output here if needed
	fmt.Println("")
	fmt.Println("Application Flags:")
	flag.PrintDefaults()
	fmt.Println("")
}

func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

var maxResults int

func main() {
	flag.IntVar(&maxResults, "max-results", 5, "max results.")
	flag.Parse()
	flag.Usage = Usage

	var bmkgURL string
	bmkgURL = "http://dataweb.bmkg.go.id/inatews/gempadirasakan.xml?decache=0.755139865912497"
	if data, err := getContent(bmkgURL); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		var infoGempa InfoGempa
		xml.Unmarshal(data, &infoGempa)
		if maxResults > len(infoGempa.InfoGempa) {
			log.Printf("Maximum results is: %v", len(infoGempa.InfoGempa))
		} else {
			for i := 0; i < maxResults; i++ {
				fmt.Println(fmt.Sprintf("--- Gempa Terkini [%d] ---", i+1))
				fmt.Println("Tanggal : " + infoGempa.InfoGempa[i].Tanggal)
				fmt.Println("Coordinates : " + infoGempa.InfoGempa[i].Point.Coordinates)
				fmt.Println("Posisi : " + infoGempa.InfoGempa[i].Posisi)
				fmt.Println("Kedalaman : " + infoGempa.InfoGempa[i].Kedalaman)
				fmt.Println("Magnitude : " + infoGempa.InfoGempa[i].Magnitude)
				fmt.Println("Keterangan : " + infoGempa.InfoGempa[i].Keterangan)
				fmt.Println("Dirasakan : " + infoGempa.InfoGempa[i].Dirasakan)
				fmt.Println("--------------------------")
				fmt.Println()
			}
		}

	}
}
