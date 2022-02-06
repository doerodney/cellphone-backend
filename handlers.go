package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	reGetAllPhones    = regexp.MustCompile("^/api/phones/?$")
	reGetPhoneById    = regexp.MustCompile("^/api/phones/([0-9]+)$")
	reGetPhonesByMake = regexp.MustCompile("^/api/phones/make/([a-zA-Z]+)$")
	reGetPhonesByOS   = regexp.MustCompile("^/api/phones/os/([a-zA-Z]+)$")
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("\nMethod: %q, Path: %q\n", r.Method, r.URL.Path)
	if r.Method == http.MethodGet {
		if r.URL.Path == "/health" {
			HealthHandler(w, r)
		} else if r.URL.Path == "/utc" {
			UtcHandler(w, r)
		} else if reGetAllPhones.MatchString(r.URL.Path) {
			GetCellPhonesHandler(w, r)
		} else if reGetPhoneById.MatchString(r.URL.Path) {
			GetCellPhoneByIdHandler(w, r)
		} else if reGetPhonesByMake.MatchString(r.URL.Path) {
			GetCellPhonesByMakeHandler(w, r)
		} else if reGetPhonesByOS.MatchString(r.URL.Path) {
			GetCellPhonesByOsHandler(w, r)
		}
	} else if r.Method == http.MethodPost {
		if reGetAllPhones.MatchString(r.URL.Path) {
			PostCellPhoneHandler(w, r)
		}
	}
}

// curl --request GET http://localhost:4000/health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service is healthy")
}

// curl --request GET http://localhost:4000/utc
func UtcHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	utcTimestamp := time.Now().UTC().String()
	fmt.Fprintf(w, "%s", utcTimestamp)
}

/*

curl --data '{"id": 0, "make": "Motorola", "model": "g power 2021", "os": "android", "releaseDate": "01/11/2021", "image": "A Motorola phone image"}' --header "Content-Type: application/json" --request POST http://localhost:4000/api/phones

*/
func PostCellPhoneHandler(w http.ResponseWriter, r *http.Request) {
	// Verify the request content type:
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var phone CellPhone

	// Unmarshall the JSON request body to a CellPhone struct:
	var unmarshallErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&phone)
	if err != nil {
		if errors.As(err, &unmarshallErr) {
			errorResponse(w, fmt.Sprintf("bad request:  incorrect type provided for field %s", unmarshallErr.Field), http.StatusBadRequest)
		} else {
			errorResponse(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadGateway)
		}
	}

	// TODO Store the struct instance:

	// Marshall the Go CellPhone struct to JSON text for the response:
	jsonTxt, err := json.Marshal(phone)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonTxt)
}

// TODO This is a stub for phone storage:
func GetCellPhones() []CellPhone {
	phones := []CellPhone{
		{
			Id:          0,
			Make:        "Motorola",
			Model:       "g power 2021",
			OS:          "android",
			ReleaseDate: "01/11/2021",
			Image:       "A Motorola phone image",
		},
		{
			Id:          1,
			Make:        "Apple",
			Model:       "iPhone 13",
			OS:          "ios",
			ReleaseDate: "11/11/2021",
			Image:       "An iPhone image",
		},
	}

	return phones
}

func GetCellPhoneById(id int) *CellPhone {
	var p *CellPhone = nil
	for _, phone := range GetCellPhones() {
		if phone.Id == id {
			p = &phone
			break
		}
	}
	return p
}

func GetCellPhonesByMake(make string) []CellPhone {
	var filtered []CellPhone
	for _, phone := range GetCellPhones() {
		if phone.Make == make {
			filtered = append(filtered, phone)
		}
	}
	return filtered
}

func GetCellPhonesByOS(os string) []CellPhone {
	var filtered []CellPhone
	for _, phone := range GetCellPhones() {
		if phone.OS == os {
			filtered = append(filtered, phone)
		}
	}
	return filtered
}

// curl --request GET http://localhost:4000/api/phones
func GetCellPhonesHandler(w http.ResponseWriter, r *http.Request) {
	phone := GetCellPhones()
	jsonTxt, err := json.Marshal(phone)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTxt)
}

// curl --request GET http://localhost:4000/api/phones/0
func GetCellPhoneByIdHandler(w http.ResponseWriter, r *http.Request) {
	match := reGetPhoneById.FindStringSubmatch(r.URL.Path)
	id, err := strconv.Atoi(match[1])
	if err != nil {
		errorResponse(w, fmt.Sprintf("id parameter cannot be converted to integer: %s", err.Error()), http.StatusBadRequest)
	}
	phone := GetCellPhoneById(int(id))
	jsonTxt, err := json.Marshal(phone)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTxt)
}

// curl --request GET http://localhost:4000/api/phones/make/Motorola
func GetCellPhonesByMakeHandler(w http.ResponseWriter, r *http.Request) {
	match := reGetPhonesByMake.FindStringSubmatch(r.URL.Path)
	make := match[1]

	phone := GetCellPhonesByMake(make)
	jsonTxt, err := json.Marshal(phone)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTxt)
}

// curl --request GET http://localhost:4000/api/phones/os/android
func GetCellPhonesByOsHandler(w http.ResponseWriter, r *http.Request) {
	match := reGetPhonesByOS.FindStringSubmatch(r.URL.Path)
	os := match[1]

	phone := GetCellPhonesByOS(os)
	jsonTxt, err := json.Marshal(phone)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTxt)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
