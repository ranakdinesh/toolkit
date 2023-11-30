package toolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var jsonTests = []struct {
	name               string
	json               string
	errorExpected      bool
	maxSize            int64
	allowUnknownFields bool
}{
	{name: "good json", json: `{"name": "John", "age": 30}`, errorExpected: false, maxSize: 1024, allowUnknownFields: false},
	{name: "badly formatted json", json: `{"name": "John", "age": }`, errorExpected: true, maxSize: 1024, allowUnknownFields: false},
	{name: "incorrect type", json: `{"name": "John", "age":"rana" }`, errorExpected: true, maxSize: 1024, allowUnknownFields: false},
	{name: "two json files", json: `{"name": "John", "age":"rana" }{"name": "John", "age":"rana" }`, errorExpected: true, maxSize: 1024, allowUnknownFields: false},
	{name: "empty body", json: ``, errorExpected: true, maxSize: 1024, allowUnknownFields: false},
	{name: "syntax Error", json: `{"name": "John", "age":3" }{"name": "John", "age":"rana" }`, errorExpected: true, maxSize: 1024, allowUnknownFields: false},
}

func TestTools_ReadJSON(t *testing.T) {
	var testTools Tools

	for _, e := range jsonTests {
		testTools.MaxJsonSize = e.maxSize
		testTools.AllowUnknownFields = e.allowUnknownFields

		var decodedJSON struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(e.json)))
		if err != nil {
			t.Log("Error:", err)
		}
		rr := httptest.NewRecorder()
		err = testTools.ReadJSON(rr, req, &decodedJSON)
		if e.errorExpected && err == nil {
			t.Errorf("Expected error, but got none")

		}
		if !e.errorExpected && err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

	}
}
func TestTools_WriteJSON(t *testing.T) {
	var testTools Tools

	rr := httptest.NewRecorder()

	payload := JSONResponse{
		Error:   false,
		Message: "Hello World",
	}

	headers := make(http.Header)
	headers.Add("foo", "bar")

	err := testTools.WriteJSON(rr, http.StatusOK, payload, headers)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestTools_ErrorJSON(t *testing.T) {
	var testTools Tools

	rr := httptest.NewRecorder()
	err := testTools.ErrorJSON(rr, errors.New("some error"), http.StatusServiceUnavailable)
	if err != nil {

		t.Error(err)

	}

	var payload JSONResponse
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		t.Error("recived error when decoding JSON", err)
	}

	if !payload.Error {
		t.Error("Expected error to be true, but got false")

	}

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("wrong status code returned; expected 503, but got %d", rr.Code)

	}

}
