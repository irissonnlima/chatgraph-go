package output_router_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	adapter_output "github.com/irissonnlima/chatgraph-go/core/ports/adapters/output"
)

const MAX_RETRIES = 5

type routerReturn struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type RouterApi struct {
	Url      string
	Username string
	Password string
}

func NewRouterApi(url, username, password string) adapter_output.RouterService {
	return &RouterApi{
		Url:      url,
		Username: username,
		Password: password,
	}
}

func (r *RouterApi) post(endpoint string, payload []byte) error {
	var err error

	for i := 0; i < MAX_RETRIES; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i) * 2 * time.Second)
		}

		req, e := http.NewRequest(http.MethodPost, r.Url+endpoint, bytes.NewBuffer(payload))
		if e != nil {
			return e
		}

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(r.Username, r.Password)

		client := &http.Client{}
		resp, e := client.Do(req)
		if e != nil {
			err = e
			log.Printf("[WARN] Request failed: %v. Retrying...", err)
			continue
		}

		body, e := io.ReadAll(resp.Body)
		resp.Body.Close()
		if e != nil {
			err = e
			log.Printf("[WARN] Failed to read body: %v. Retrying...", err)
			continue
		}

		var result routerReturn
		if e := json.Unmarshal(body, &result); e != nil {
			err = e
			log.Printf("[WARN] Failed to unmarshal: %v. Retrying...", err)
			continue
		}

		if result.Status {
			log.Printf("[INFO] %s", result.Message)
			return nil
		}

		err = errors.New(result.Message)
		log.Printf("[ERROR] %s", result.Message)
	}

	return err
}

func (r *RouterApi) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, r.Url+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(r.Username, r.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *RouterApi) uploadFileMultipart(endpoint, filePath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffer to write the multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Create the form file field "content"
	part, err := writer.CreateFormFile("content", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// Copy file content to the form field
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPost, r.Url+endpoint, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(r.Username, r.Password)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
