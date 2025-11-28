// Package d_file provides file-related domain models for handling media attachments.
// It supports images, videos, audio files, and generic file types.
package d_file

import (
	"fmt"
	"path/filepath"
	"strings"
)

// EmptyFile represents a zero-value File for comparison purposes.
var EmptyFile = File{}

// FileType represents the type of file being sent or received.
type FileType int

// File type constants.
const (
	// IMAGE_SEND_TYPE represents an image file (e.g., PNG, JPG, GIF).
	IMAGE_SEND_TYPE FileType = iota
	// VIDEO_SEND_TYPE represents a video file (e.g., MP4, AVI).
	VIDEO_SEND_TYPE
	// AUDIO_SEND_TYPE represents an audio file (e.g., MP3, WAV).
	AUDIO_SEND_TYPE
	// FILE_SEND_TYPE represents a generic file (e.g., PDF, DOC).
	FILE_SEND_TYPE
)

// String returns the string representation of the FileType.
func (st FileType) String() string {
	switch st {
	case IMAGE_SEND_TYPE:
		return "IMAGE"
	case VIDEO_SEND_TYPE:
		return "VIDEO"
	case FILE_SEND_TYPE:
		return "FILE"
	default:
		return "UNKNOWN"
	}
}

// SendTypeFromString converts a string to a FileType.
// Returns an error if the string does not match a valid FileType.
func SendTypeFromString(sendType string) (FileType, error) {
	switch sendType {
	case "IMAGE":
		return IMAGE_SEND_TYPE, nil
	case "VIDEO":
		return VIDEO_SEND_TYPE, nil
	case "FILE":
		return FILE_SEND_TYPE, nil
	default:
		return -1, fmt.Errorf("invalid send type: %s", sendType)
	}
}

// File represents a file attachment in a message.
type File struct {
	// ID is the unique identifier for the file.
	ID string
	// Type indicates the file type (image, video, audio, or file).
	Type FileType
	// URL is the location where the file can be accessed.
	URL string
	// Name is the filename including extension (e.g., "document.pdf").
	Name string
}

// IsEmpty returns true if the File has no ID or equals EmptyFile.
func (f File) IsEmpty() bool {
	return f == EmptyFile || f.ID == ""
}

// Extension returns the lowercase file extension including the dot.
// For example, "document.PDF" returns ".pdf".
func (f File) Extension() string {
	return strings.ToLower(filepath.Ext(f.Name))
}
