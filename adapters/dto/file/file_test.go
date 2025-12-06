package dto_file

import (
	"testing"

	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
)

func TestFile_ToDomain(t *testing.T) {
	tests := []struct {
		name     string
		file     File
		wantID   string
		wantType d_file.FileType
		wantURL  string
		wantName string
	}{
		{
			name:     "image file",
			file:     File{ID: "1", Type: "IMAGE", URL: "https://x.com/img.png", Name: "img.png"},
			wantID:   "1",
			wantType: d_file.IMAGE_SEND_TYPE,
			wantURL:  "https://x.com/img.png",
			wantName: "img.png",
		},
		{
			name:     "video file",
			file:     File{ID: "2", Type: "VIDEO", URL: "https://x.com/vid.mp4", Name: "vid.mp4"},
			wantID:   "2",
			wantType: d_file.VIDEO_SEND_TYPE,
			wantURL:  "https://x.com/vid.mp4",
			wantName: "vid.mp4",
		},
		{
			name:     "generic file",
			file:     File{ID: "4", Type: "FILE", URL: "https://x.com/doc.pdf", Name: "doc.pdf"},
			wantID:   "4",
			wantType: d_file.FILE_SEND_TYPE,
			wantURL:  "https://x.com/doc.pdf",
			wantName: "doc.pdf",
		},
		{
			name:     "invalid type defaults to FILE_SEND_TYPE",
			file:     File{ID: "5", Type: "invalid", URL: "https://x.com/x", Name: "x"},
			wantID:   "5",
			wantType: d_file.FILE_SEND_TYPE,
			wantURL:  "https://x.com/x",
			wantName: "x",
		},
		{
			name:     "empty type defaults to FILE_SEND_TYPE",
			file:     File{ID: "6", Type: "", URL: "https://x.com/y", Name: "y"},
			wantID:   "6",
			wantType: d_file.FILE_SEND_TYPE,
			wantURL:  "https://x.com/y",
			wantName: "y",
		},
		{
			name:     "lowercase type defaults to FILE_SEND_TYPE",
			file:     File{ID: "7", Type: "image", URL: "https://x.com/z", Name: "z"},
			wantID:   "7",
			wantType: d_file.FILE_SEND_TYPE,
			wantURL:  "https://x.com/z",
			wantName: "z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.file.ToDomain()
			if got.ID != tt.wantID {
				t.Errorf("File.ToDomain().ID = %v, want %v", got.ID, tt.wantID)
			}
			if got.Type != tt.wantType {
				t.Errorf("File.ToDomain().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.URL != tt.wantURL {
				t.Errorf("File.ToDomain().URL = %v, want %v", got.URL, tt.wantURL)
			}
			if got.Name != tt.wantName {
				t.Errorf("File.ToDomain().Name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}
