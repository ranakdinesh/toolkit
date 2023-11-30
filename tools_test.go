package toolkit

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Errorf("RandomString() = %v, want %v", len(s), 10)
	}

}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpg", "image/png", "image/jpeg"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedTypes: []string{"image/jpg", "image/png", "image/jpeg"}, renameFile: true, errorExpected: false},
	{name: "allowed rename", allowedTypes: []string{"image/jpg", "image/png"}, renameFile: true, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {

	for _, e := range uploadTests {

		// set up a pipe to avod buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer writer.Close()
			defer wg.Done()
			//crate a data filed 'file'
			part, err := writer.CreateFormFile("file", "./testdata/testimage.jpg")
			if err != nil {
				t.Error(err)
			}
			f, err := os.Open("../testdata/testimage.jpg")
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error(err)
			}
			err = jpeg.Encode(part, img, nil)
			if err != nil {
				t.Error(err)
			}
			err = jpeg.Encode(part, img, &jpeg.Options{Quality: 1})
			if err != nil {
				t.Error(err)
			}

		}()
		//read from the file which recieve the data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "../testdata/uploads/", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("../testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {

				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}

			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))

		}
		if !e.errorExpected && err != nil {
			t.Error(fmt.Sprintf("%s: error expected but none recived", e.name))
		}

	}
}
func TestTools_UploadOneFile(t *testing.T) {

	// set up a pipe to avod buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		//crate a data filed 'file'
		part, err := writer.CreateFormFile("file", "./testdata/testimage.jpg")
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open("../testdata/testimage.jpg")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error(err)
		}
		err = jpeg.Encode(part, img, nil)
		if err != nil {
			t.Error(err)
		}
		err = jpeg.Encode(part, img, &jpeg.Options{Quality: 1})
		if err != nil {
			t.Error(err)
		}

	}()
	//read from the file which recieve the data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools

	uploadedFile, err := testTools.UploadOneFile(request, "../testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("../testdata/uploads/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {

		t.Errorf("expected file to exist: %s", err.Error())
	}

	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))

}
