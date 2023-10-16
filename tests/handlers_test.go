package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"

	"github.com/toggle-assignment/routes"
)

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Payload struct {
	body        *bytes.Buffer
	contentType string
}

var processedReceiptUUID string

func buildPayload(t *testing.T, fileName string, filePath string) Payload {
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	//Create multipart header
	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", fileName))
	fileHeader.Set("Content-Type", "multipart/form-data")
	writer, err := multipartWriter.CreatePart(fileHeader)
	if err != nil {
		t.Fatal(err)
	}
	//Copy file to file multipart writer
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(writer, file)
	// close the writer before making the request
	multipartWriter.Close()

	var contentType = multipartWriter.FormDataContentType()

	var response = Payload{body, contentType}
	return response
}

func makeHttpRequest(t *testing.T, method string, url string, payload *Payload) *httptest.ResponseRecorder {
	r := routes.SetupRoutes()

	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else if method == "POST" && payload != nil {
		req, err = http.NewRequest(method, url, payload.body)
		req.Header.Add("Content-Type", payload.contentType)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	return rr
}

func TestGetRouteShouldReturnIsProcessingWhenReceiptIsProcessing(t *testing.T) {
	var apiResponse Response

	rr := makeHttpRequest(t, "GET", "/receipt/5a2ac852-1489-47ca-9487-82194087f934", nil)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	jsonData := rr.Body.String()

	json.Unmarshal([]byte(jsonData), &apiResponse)

	expectedResponse := "receipt is still processing"

	if apiResponse.Message != expectedResponse {
		t.Errorf("Expected response body %s, but got %s", expectedResponse, jsonData)
	}
}

func TestGetRouteShouldReturnNotFoundWhenReceiptDoesNotExist(t *testing.T) {
	rr := makeHttpRequest(t, "GET", "/receipt/5a2ac852", nil)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
	}
}

func TestPostRouteShouldReturnBadRequestWhenFileIsNotPresentInRequest(t *testing.T) {
	var apiResponse Response

	rr := makeHttpRequest(t, "POST", "/receipt", nil)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}

	jsonData := rr.Body.String()

	json.Unmarshal([]byte(jsonData), &apiResponse)

	expectedResponse := "File not found in the request."

	if apiResponse.Error != expectedResponse {
		t.Errorf("Expected response body %s, but got %s", expectedResponse, jsonData)
	}
}

func TestPostRouteShouldReturnBadRequestWhenFileIsNotSupported(t *testing.T) {
	var apiResponse Response

	var payload = buildPayload(t, "example.txt", "./test_files/example.txt")

	rr := makeHttpRequest(t, "POST", "/receipt", &payload)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}

	jsonData := rr.Body.String()

	json.Unmarshal([]byte(jsonData), &apiResponse)

	expectedResponse := "Invalid file format. Supported formats: PNG, JPG, PDF."

	if apiResponse.Error != expectedResponse {
		t.Errorf("Expected response body %s, but got %s", expectedResponse, jsonData)
	}
}

func TestPostRouteShouldSuccessWhenCorrectReceiptIsUploaded(t *testing.T) {

	var payload = buildPayload(t, "test_receipt.pdf", "./test_files/test_receipt.pdf")

	rr := makeHttpRequest(t, "POST", "/receipt", &payload)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	response := rr.Body.String()

	json.Unmarshal([]byte(response), &processedReceiptUUID)

	if processedReceiptUUID == "" {
		t.Errorf("Did not match expected response")
	}
	var apiResponse Response

	url := fmt.Sprintf("/receipt/%s", processedReceiptUUID)

	newReq := makeHttpRequest(t, "GET", url, nil)

	jsonData := newReq.Body.String()

	json.Unmarshal([]byte(jsonData), &apiResponse)

	expectedResponse := "receipt is still processing"

	if apiResponse.Message != expectedResponse {
		t.Errorf("Expected response body %s, but got %s", expectedResponse, jsonData)
	}
}
