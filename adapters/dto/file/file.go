// Package d_file provides file-related domain models for handling media attachments.
// It supports images, videos, audio files, and generic file types.
package dto_file

import d_file "chatgraph/core/domain/file"

// File represents a file attachment in a message.
type File struct {
	// ID is the unique identifier for the file.
	ID string `json:"id"`
	// Type indicates the file type (image, video, audio, or file).
	Type string `json:"type"`
	// URL is the location where the file can be accessed.
	URL string `json:"url"`
	// Name is the filename including extension (e.g., "document.pdf").
	Name string `json:"name"`
}

func (f File) ToDomain() d_file.File {
	t, err := d_file.SendTypeFromString(f.Type)
	if err != nil {
		t = d_file.FILE_SEND_TYPE
	}

	return d_file.File{
		ID:   f.ID,
		Type: t,
		URL:  f.URL,
		Name: f.Name,
	}
}
