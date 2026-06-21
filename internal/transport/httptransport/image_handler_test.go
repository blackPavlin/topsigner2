package httptransport_test

import (
	"testing"
)

func TestImageHandlerGetImages(t *testing.T) {
	t.Parallel()

	tests := []struct{ name string }{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
		})
	}
}

func TestImageHandlerUploadImage(t *testing.T) {
	t.Parallel()

	tests := []struct{ name string }{}

	// handler := httptransport.NewHandler(&config.Config{}, image.New())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
		})
	}
}

func TestImageHandlerDeleteImageByID(t *testing.T) {
	t.Parallel()
}
