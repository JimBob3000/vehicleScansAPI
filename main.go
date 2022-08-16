package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Mot []MotVehicle

type MotVehicle struct {
	Registration  string     `json:"registration"`
	Make          string     `json:"make"`
	Model         string     `json:"model"`
	FirstUsedDate string     `json:"firstUsedDate"`
	FuelType      string     `json:"fuelType"`
	PrimaryColour string     `json:"primaryColour"`
	MotTests      []MotTests `json:"motTests"`
}

type MotTests struct {
	CompletedDate  string        `json:"completedDate"`
	TestResult     string        `json:"testResult"`
	ExpiryDate     string        `json:"expiryDate"`
	OdometerValue  string        `json:"odometerValue"`
	OdometerUnit   string        `json:"odometerUnit"`
	MotTestNumber  string        `json:"motTestNumber"`
	RfrAndComments []MotComments `json:"rfrAndComments"`
}

type MotComments struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type DVLA struct {
	RegistrationNumber           string `json:"registrationNumber"`
	TaxStatus                    string `json:"taxStatus"`
	TaxDueDate                   string `json:"taxDueDate"`
	ArtEndDate                   string `json:"artEndDate"`
	MotStatus                    string `json:"motStatus"`
	MotExpiryDate                string `json:"motExpiryDate"`
	Make                         string `json:"make"`
	MonthOfFirstDvlaRegistration string `json:"monthOfFirstDvlaRegistration"`
	MonthOfFirstRegistration     string `json:"monthOfFirstRegistration"`
	YearOfManufacture            int    `json:"yearOfManufacture"`
	EngineCapacity               int    `json:"engineCapacity"`
	Co2Emissions                 int    `json:"co2Emissions"`
	FuelType                     string `json:"fuelType"`
	MarkedForExport              bool   `json:"markedForExport"`
	Colour                       string `json:"colour"`
	TypeApproval                 string `json:"typeApproval"`
	Wheelplan                    string `json:"wheelplan"`
	RevenueWeight                int    `json:"revenueWeight"`
	RealDrivingEmissions         string `json:"realDrivingEmissions"`
	DateOfLastV5CIssued          string `json:"dateOfLastV5CIssued"`
	EuroStatus                   string `json:"euroStatus"`
}

func getMotHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vrn := vars["vrn"]
	fmt.Println("Endpoint Hit: getMotHistory, VRN: " + vrn)

	url := "https://beta.check-mot.service.gov.uk/trade/vehicles/mot-tests?registration=" + vrn
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", os.Getenv("MOT_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var responseObject Mot
	json.Unmarshal(body, &responseObject)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(responseObject[0])
}

func getMotRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	fmt.Println("Endpoint Hit: getMotRecords, page: " + page)

	url := "https://beta.check-mot.service.gov.uk/trade/vehicles/mot-tests?page=" + page
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", os.Getenv("MOT_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var responseObject Mot
	json.Unmarshal(body, &responseObject)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(responseObject)
}

func getDvlaRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vrn := vars["vrn"]
	fmt.Println("Endpoint Hit: getDvlaRecord, VRN: " + vrn)

	url := "https://driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles"
	var jsonStr = []byte(`{"registrationNumber":"` + vrn + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("x-api-key", os.Getenv("DVLA_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var responseObject DVLA
	json.Unmarshal(body, &responseObject)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(responseObject)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "API health: Good")`))
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/mot/{vrn}", getMotHistory)
	myRouter.HandleFunc("/motPage/{page}", getMotRecords)
	myRouter.HandleFunc("/dvla/{vrn}", getDvlaRecord)
	myRouter.HandleFunc("/health", healthCheck)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("API Deploying...")
	handleRequests()
}
