package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

// Tools is the type used to instatiate this module, Any variable of this type has access to to all the methods with the reciver (t *tools) *Tools
type Tools struct {
	MaxFileSize        int64
	AllowedFileTypes   []string
	MaxJsonSize        int64
	AllowUnknownFields bool
}

// RandomString is used to generate a random string of n length
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomString)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}

// UploadedFile is a struct used to save infomration about an uploaded file
type UploadFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}
	//if user does not provide a max file size, set it to 10MB
	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 10
	}

	var uploadedFiles []*UploadFile
	err := r.ParseMultipartForm(t.MaxFileSize)
	if err != nil {
		return nil, errors.New("file size too large")
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {

			uploadedFiles, err = func(uploadedFiles []*UploadFile) ([]*UploadFile, error) {
				var uploadedFile UploadFile
				infile, err := fileHeader.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()
				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				//check to see if the file type is permietted

				allowed := false

				fileType := http.DetectContentType(buff)
				if len(t.AllowedFileTypes) > 0 {

					for _, allowedFileType := range t.AllowedFileTypes {
						if strings.EqualFold(fileType, allowedFileType) {
							allowed = true
						}
					}
				} else {
					allowed = false
				}
				if !allowed {
					return nil, errors.New("File type not allowed")
				}
				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}
				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(fileHeader.Filename))
				} else {
					uploadedFile.NewFileName = fileHeader.Filename

				}
				uploadedFile.OriginalFileName = fileHeader.Filename
				var outfile *os.File
				defer outfile.Close()
				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize

				}
				uploadedFiles = append(uploadedFiles, &uploadedFile)
				return uploadedFiles, nil

			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}

		}

	}
	return uploadedFiles, nil
}

func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadFile, error) {

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 10
	}
	files, err := t.UploadFiles(r, uploadDir, rename...)
	if err != nil {
		return nil, err
	}
	return files[0], nil

}
