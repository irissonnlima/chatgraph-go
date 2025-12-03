package output_router_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	adapter_output "chatgraph/core/ports/adapters/output"
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
