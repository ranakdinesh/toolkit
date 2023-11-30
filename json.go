package toolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data, omitempty"`
}

func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := int64(1048576) // 1MB
	if t.MaxJsonSize > 0 {
		maxBytes = t.MaxJsonSize
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)
	if !t.AllowUnknownFields {
		dec.DisallowUnknownFields()
	}
	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return errors.New(fmt.Sprintf("body contains badly-formed JSON (at position %d)", syntaxError.Offset))
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return errors.New(fmt.Sprintf("body contains incorrect JSON type for field %s", unmarshalTypeError.Field))
			}
			return errors.New(fmt.Sprintf("body contains incorrect JSON type (at position %d", unmarshalTypeError.Offset))
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldname := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errors.New("body contains unknown field " + fieldname)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes ", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshalling JSON %s", err.Error())

		default:

			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}

//WriteJSON takes a respose status core and arbitrary data and writes it to the response writer

func (t *Tools) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header().Set(k, v[0])
		}

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil

}

//ErrorJSON takes an error, and opionally a status core and generated and sends a JSON error message

func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return t.WriteJSON(w, statusCode, payload)
}

func (t *Tools) PushJSONToRemote(uri string, data any, client ...*http.Client) (*http.Response, int, error) {
	//create json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err

	}
	//check the custom http clinet

	httpClient := &http.Client{}
	if len(client) > 0 {
		httpClient = client[0]
	}

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}
	//build the request and set the header

	request.Header.Set("Content-Type", "application/json")

	//call the remote uri

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	return response, response.StatusCode, nil

	//send response back

}
