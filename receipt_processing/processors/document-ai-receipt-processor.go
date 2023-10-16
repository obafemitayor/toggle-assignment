package processors

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ReceiptDetails struct {
	Text string `json:"text"`
}

type ReceiptResponse struct {
	Document ReceiptDetails `json:"document"`
}

type DocumentAIReceiptProcessor struct{}

func (rp *DocumentAIReceiptProcessor) ExtractDetailsFromReceipt(filePath string) string {
	payload := buildPayload(filePath)

	response := getDetailsOfReceiptFromDocumentAI(payload)

	return buildResponse(string(response))
}

func buildPayload(filePath string) *strings.Reader {
	payload := strings.NewReader(`{
		"skipHumanReview": true,
		"rawDocument": {
		"mimeType": "application/pdf",
		"content": "` + convertReceiptToBase64String(filePath) + `"
	},
	"fieldMask": "text"
	}`)
	return payload
}

func getDetailsOfReceiptFromDocumentAI(payload *strings.Reader) []byte {
	method := "POST"

	client := &http.Client{}

	req, err := http.NewRequest(method, getRequestUrl(), payload)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+getAuthToken()+"")

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return body
}

func convertReceiptToBase64String(filePath string) string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(fileContent)
}

func getRequestUrl() string {
	projectID := "235872245316"
	location := "eu"
	processorID := "5cfedb6edd696fb"

	return fmt.Sprintf("https://eu-documentai.googleapis.com/v1/projects/%s/locations/%s/processors/%s:process", projectID, location, processorID)
}

func getAuthToken() string {
	return "ya29.a0AfB_byAno_Fu8Pjgmrg2faHU5bzk-nSlfd-kppth7EmU1k-5dxXGJVpBmqkmF0cJl8M688TXxCCoITIS_GK4k_bVLc8VHtqDdPW5xrZGUzVFFKoit5V1wQchnmOIcO5FTwgdI27YlczLLwHNljpkxCG-zO9koWdg_WWXW7i2CvUaCgYKAbcSARESFQGOcNnCrlZ40z3owojC4AnzWpcSLA0178"
}

func buildResponse(documentResponse string) string {
	var receiptResponse ReceiptResponse

	respErr := json.Unmarshal([]byte(documentResponse), &receiptResponse)

	if respErr != nil {
		fmt.Println("Error:", respErr)
		return ""
	}

	return receiptResponse.Document.Text
}
