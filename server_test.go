package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetCellPhones(t *testing.T) {
	phoneList := GetCellPhones()
	expected := 2
	observed := len(phoneList)
	if expected != observed {
		t.Fatalf("Expected: %q, Observed: %q", expected, observed)
	}
}


func TestGetCellPhoneById(t *testing.T) {
	var p *CellPhone
	
	// Validate that a valid pointer is returned for a valid id
	expected := 0
	p = GetCellPhoneById(expected)
	if p == nil {
		t.Fatal("Expected: valid cell phone pointer, Observed: nil pointer")
	} 

	// Validate the result
	observed := p.Id
	if expected != observed {
		t.Fatalf("Expected: %q, Observed: %q", expected, observed)
	}

	// Validate that a nil pointer is returned for invalid id 
	expected = -1
	p = GetCellPhoneById(expected)
	if p != nil {
		t.Fatal("Expected: nil cell phone pointer, Observed: valid cell phone pointer")
	} 
}
func TestGetCellPhonesByMake(t *testing.T) {
	makes := [...]string{"Motorola", "Apple"}
	for _, expected := range makes {
		for _, phone := range GetCellPhonesByMake(expected) {
			observed := phone.Make
			if expected != observed {
				t.Fatalf("Expected: %q, Observed: %q", expected, observed)
			}
		}
	}
}

func TestGetCellPhonesByOS(t *testing.T) {
	makes := [...]string{"android", "ios"}
	for _, expected := range makes {
		for _, phone := range GetCellPhonesByOS(expected) {
			observed := phone.OS
			if expected != observed {
				t.Fatalf("Expected: %q, Observed: %q", expected, observed)
			}
		}
	}
}

func TestGetIpPortText(t *testing.T) {
	ip := "127.0.0.1"
	port := 4321

	expected := fmt.Sprintf("%s:%d", ip, port)
	observed := GetIpPortText(ip, port)
	if expected != observed {
		t.Fatalf("Expected: %q, Observed: %q", expected, observed)
	}
}

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	testServer := httptest.NewServer(newMultiplexer())

	return testServer.URL, func() { testServer.Close() }
}

type GetTestCase struct {
	name            string
	path            string
	expectedCode    int
	expectedContent []string
}

func TestGet(t *testing.T) {
	testCases := []GetTestCase{
		{
			name:            "health",
			path:            "/health",
			expectedCode:    http.StatusOK,
			expectedContent: []string{"Service is healthy"},
		},
		{
			name:            "utc",
			path:            "/utc",
			expectedCode:    http.StatusOK,
			expectedContent: []string{"UTC"},
		},
		{
			name:            "phones",
			path:            "/api/phones",
			expectedCode:    http.StatusOK,
			expectedContent: []string{"id", "make", "model", "os", "releaseDate", "image"},
		},
		{
			name:            "phones",
			path:            "/api/phones/0",
			expectedCode:    http.StatusOK,
			expectedContent: []string{"id", "make", "model", "os", "releaseDate", "image"},
		},
		{
			name:            "phones",
			path:            "/api/phones/make/Motorola",
			expectedCode:    http.StatusOK,
			expectedContent: []string{"id", "make", "model", "os", "releaseDate", "image", "Motorola"},
		},
		{
			name:            "phones",
			path:            "/api/phones/os/android",
			expectedCode:    http.StatusOK,
			expectedContent: []string{"id", "make", "model", "os", "releaseDate", "image", "Motorola",},
		},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)

			res, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}
			defer res.Body.Close()

			// Test for expected status code:
			if res.StatusCode != tc.expectedCode {
				t.Fatalf("Expected: %q, Observed: %q", http.StatusText(tc.expectedCode), http.StatusText(res.StatusCode))
			}

			// Test for expected content:
			switch { // No break required in Go, unlike C/C++/Javascript
			// Test for expected text/plain content:
			case strings.Contains(res.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(res.Body); err != nil {
					t.Error(err)
				}
				for _, expected := range tc.expectedContent {
					if !strings.Contains(string(body), expected) {
						t.Errorf("Expected %q, Observed %q.", expected, string(body))
					}
				}

			//Test for expected application/json content:
			case strings.Contains(res.Header.Get("Content-Type"), "application/json"):
				if body, err = io.ReadAll(res.Body); err != nil {
					t.Error(err)
				}
				for _, expected := range tc.expectedContent {
					if !strings.Contains(string(body), expected) {
						t.Errorf("Expected %q, Observed %q.", expected, string(body))
					}
				}

			default:
				t.Fatalf("Unsupported Content-Type: %q", res.Header.Get("Content-Type"))
			}
		})
	}
}

type PostTestCase struct {
	name            string
	path            string
	contentType     string
	requestBody     string
	expectedCode    int
	expectedContent []string
}

func TestPost(t *testing.T) {
	// Post request contains url, content type, body
	// Post response contains response body, status code

	testCases := []PostTestCase{
		{
			name:        "phones",
			path:        "/api/phones",
			contentType: "application/json",
			requestBody: `{
				"id": 0,
				"make": "Motorola",
				"model": "g power 2021",
				"os": "android",
				"releaseDate": "01/11/2021",
				"image": "A Motorola phone image"
			}`,
			expectedCode:    http.StatusCreated,
			expectedContent: []string{"id", "make", "model", "os", "releaseDate", "image"},
		},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// For the http Post method, the body argument must implement the
			// io.reader reader.
			reqBody := strings.NewReader(tc.requestBody)

			res, err := http.Post(url+tc.path, tc.contentType, reqBody)
			if err != nil {
				t.Error(err)
			}
			defer res.Body.Close()

			// Test for expected status code:
			if res.StatusCode != tc.expectedCode {
				t.Fatalf("Expected: %q, Observed: %q", http.StatusText(tc.expectedCode), http.StatusText(res.StatusCode))
			}

			// Test for expected response content:
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}
			for _, expected := range tc.expectedContent {
				if !strings.Contains(string(body), expected) {
					t.Errorf("Expected %q, Observed %q.", expected, string(body))
				}
			}
		})
	}
}
