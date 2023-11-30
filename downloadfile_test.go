package toolkit

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func TestTools_DownloadStaticFile(t *testing.T) {
	var testTools Tools

	// Create a mock HTTP request
	req, err := http.NewRequest("GET", "/download", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Set the file path and display name
	p := "/path/to/files"
	file := "example.txt"
	displayName := "example_download.txt"

	// Call the DownloadStaticFile method
	testTools.DownloadStaticFile(rr, req, p, file, displayName)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("DownloadStaticFile returned wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	// Check the response header
	expectedHeader := "attachment; filename=\"example_download.txt\""
	if rr.Header().Get("Content-Disposition") != expectedHeader {
		t.Errorf("DownloadStaticFile returned wrong header: got %v, want %v", rr.Header().Get("Content-Disposition"), expectedHeader)
	}

	// Check the response body
	expectedFilePath := path.Join(p, file)
	if rr.Body.String() != expectedFilePath {
		t.Errorf("DownloadStaticFile returned wrong body: got %v, want %v", rr.Body.String(), expectedFilePath)
	}
}
