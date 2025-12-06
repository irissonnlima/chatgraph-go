package d_context

import (
	d_file "chatgraph/core/domain/file"
	"os"
	"path/filepath"
)

func (c *ChatContext[Obs]) LoadFile(filePath string) (*d_file.File, error) {
	if c.Context.Err() != nil {
		return nil, c.Context.Err()
	}

	return c.router.UploadFile(filePath)
}

func (c *ChatContext[Obs]) LoadFileBytes(fileName string, data []byte) (*d_file.File, error) {
	if c.Context.Err() != nil {
		return nil, c.Context.Err()
	}

	// Get the base name and extension from the provided fileName
	ext := filepath.Ext(fileName)
	baseName := fileName[:len(fileName)-len(ext)]

	// Create a temporary file with pattern: baseName-*.ext (e.g., relatorio-abc123.pdf)
	tempFile, err := os.CreateTemp("", baseName+"-*"+ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file afterwards

	// Write the data to the temporary file
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return nil, err
	}
	tempFile.Close()

	// Use the UploadFile method to upload the temporary file
	uploadedFile, err := c.router.UploadFile(tempFile.Name())
	if err != nil {
		return nil, err
	}

	return uploadedFile, nil
}

func (c *ChatContext[Obs]) GetFile(fileID string) (*d_file.File, error) {
	if c.Context.Err() != nil {
		return nil, c.Context.Err()
	}

	return c.router.GetFile(fileID)
}
