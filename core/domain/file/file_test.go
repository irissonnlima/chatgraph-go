package d_file

import (
	"testing"
)

func TestFileType_String(t *testing.T) {
	tests := []struct {
		name     string
		fileType FileType
		want     string
	}{
		{
			name:     "IMAGE_SEND_TYPE returns IMAGE",
			fileType: IMAGE_SEND_TYPE,
			want:     "IMAGE",
		},
		{
			name:     "VIDEO_SEND_TYPE returns VIDEO",
			fileType: VIDEO_SEND_TYPE,
			want:     "VIDEO",
		},
		{
			name:     "FILE_SEND_TYPE returns FILE",
			fileType: FILE_SEND_TYPE,
			want:     "FILE",
		},
		{
			name:     "AUDIO_SEND_TYPE returns UNKNOWN",
			fileType: AUDIO_SEND_TYPE,
			want:     "UNKNOWN",
		},
		{
			name:     "unknown type returns UNKNOWN",
			fileType: FileType(99),
			want:     "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileType.String()
			if got != tt.want {
				t.Errorf("FileType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendTypeFromString(t *testing.T) {
	tests := []struct {
		name      string
		sendType  string
		want      FileType
		wantError bool
	}{
		{
			name:      "IMAGE returns IMAGE_SEND_TYPE",
			sendType:  "IMAGE",
			want:      IMAGE_SEND_TYPE,
			wantError: false,
		},
		{
			name:      "VIDEO returns VIDEO_SEND_TYPE",
			sendType:  "VIDEO",
			want:      VIDEO_SEND_TYPE,
			wantError: false,
		},
		{
			name:      "FILE returns FILE_SEND_TYPE",
			sendType:  "FILE",
			want:      FILE_SEND_TYPE,
			wantError: false,
		},
		{
			name:      "invalid type returns error",
			sendType:  "INVALID",
			want:      -1,
			wantError: true,
		},
		{
			name:      "empty string returns error",
			sendType:  "",
			want:      -1,
			wantError: true,
		},
		{
			name:      "lowercase image returns error",
			sendType:  "image",
			want:      -1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SendTypeFromString(tt.sendType)
			if (err != nil) != tt.wantError {
				t.Errorf("SendTypeFromString() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("SendTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		file File
		want bool
	}{
		{
			name: "empty file returns true",
			file: File{},
			want: true,
		},
		{
			name: "EmptyFile constant returns true",
			file: EmptyFile,
			want: true,
		},
		{
			name: "file with empty ID returns true",
			file: File{
				ID:   "",
				Type: IMAGE_SEND_TYPE,
				URL:  "https://example.com/image.png",
				Name: "image.png",
			},
			want: true,
		},
		{
			name: "file with ID returns false",
			file: File{
				ID:   "123",
				Type: IMAGE_SEND_TYPE,
				URL:  "https://example.com/image.png",
				Name: "image.png",
			},
			want: false,
		},
		{
			name: "file with only ID returns false",
			file: File{
				ID: "abc",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.file.IsEmpty()
			if got != tt.want {
				t.Errorf("File.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Extension(t *testing.T) {
	tests := []struct {
		name string
		file File
		want string
	}{
		{
			name: "returns lowercase extension for png",
			file: File{
				ID:   "123",
				Name: "image.png",
			},
			want: ".png",
		},
		{
			name: "returns lowercase extension for uppercase PNG",
			file: File{
				ID:   "123",
				Name: "image.PNG",
			},
			want: ".png",
		},
		{
			name: "returns lowercase extension for mixed case",
			file: File{
				ID:   "123",
				Name: "document.PdF",
			},
			want: ".pdf",
		},
		{
			name: "returns empty string for file without extension",
			file: File{
				ID:   "123",
				Name: "filename",
			},
			want: "",
		},
		{
			name: "returns empty string for empty name",
			file: File{
				ID:   "123",
				Name: "",
			},
			want: "",
		},
		{
			name: "handles multiple dots in filename",
			file: File{
				ID:   "123",
				Name: "file.backup.tar.gz",
			},
			want: ".gz",
		},
		{
			name: "handles hidden files",
			file: File{
				ID:   "123",
				Name: ".gitignore",
			},
			want: ".gitignore",
		},
		{
			name: "handles file ending with dot",
			file: File{
				ID:   "123",
				Name: "file.",
			},
			want: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.file.Extension()
			if got != tt.want {
				t.Errorf("File.Extension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmptyFile(t *testing.T) {
	// Verify EmptyFile is a zero-value File
	expected := File{}
	if EmptyFile != expected {
		t.Errorf("EmptyFile should be equal to zero-value File")
	}

	// Verify EmptyFile.IsEmpty() returns true
	if !EmptyFile.IsEmpty() {
		t.Errorf("EmptyFile.IsEmpty() should return true")
	}
}

func TestFileTypeConstants(t *testing.T) {
	// Verify the iota values are as expected
	tests := []struct {
		name     string
		fileType FileType
		want     int
	}{
		{
			name:     "IMAGE_SEND_TYPE is 0",
			fileType: IMAGE_SEND_TYPE,
			want:     0,
		},
		{
			name:     "VIDEO_SEND_TYPE is 1",
			fileType: VIDEO_SEND_TYPE,
			want:     1,
		},
		{
			name:     "AUDIO_SEND_TYPE is 2",
			fileType: AUDIO_SEND_TYPE,
			want:     2,
		},
		{
			name:     "FILE_SEND_TYPE is 3",
			fileType: FILE_SEND_TYPE,
			want:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.fileType) != tt.want {
				t.Errorf("FileType constant = %v, want %v", int(tt.fileType), tt.want)
			}
		})
	}
}

func TestFile_Bytes_EmptyURL(t *testing.T) {
	file := File{
		ID:   "123",
		Name: "test.txt",
		URL:  "",
	}

	_, err := file.Bytes()
	if err == nil {
		t.Error("File.Bytes() should return error for empty URL")
	}
	if err.Error() != "file URL is empty" {
		t.Errorf("File.Bytes() error = %v, want 'file URL is empty'", err)
	}
}

func TestFile_Bytes_InvalidURL(t *testing.T) {
	file := File{
		ID:   "123",
		Name: "test.txt",
		URL:  "not-a-valid-url",
	}

	_, err := file.Bytes()
	if err == nil {
		t.Error("File.Bytes() should return error for invalid URL")
	}
}
