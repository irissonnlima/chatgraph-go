package output_router_api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"

	dto_file "github.com/irissonnlima/chatgraph-go/adapters/dto/file"
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
)

type fileReturn struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    *dto_file.File `json:"data"`
}

func (r *RouterApi) GetFile(fileID string) (*d_file.File, error) {
	returnBytes, err := r.get("/v1/actions/files/" + fileID)
	if err != nil {
		return nil, err
	}

	var ret fileReturn
	err = json.Unmarshal(returnBytes, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Status && ret.Data != nil {
		domainFile := ret.Data.ToDomain()
		return &domainFile, nil
	}

	if !ret.Status {
		return nil, errors.New(ret.Message)
	}

	return nil, nil
}

func (r *RouterApi) UploadFile(filepath string) (*d_file.File, error) {

	// Read file bytes
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Make hash sha256 of the fileBytes to check the id
	hash := sha256.Sum256(fileBytes)
	hashHex := hex.EncodeToString(hash[:])

	file, err := r.GetFile(hashHex)
	if err != nil {
		return nil, err
	}

	if file != nil && !file.IsEmpty() {
		return file, nil
	}

	// Upload file using multipart/form-data
	responseBytes, err := r.uploadFileMultipart("/v1/actions/files/upload", filepath)
	if err != nil {
		return nil, err
	}

	var ret fileReturn
	err = json.Unmarshal(responseBytes, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Status && ret.Data != nil {
		domainFile := ret.Data.ToDomain()
		return &domainFile, nil
	}

	if !ret.Status {
		return nil, errors.New(ret.Message)
	}

	return nil, nil
}
