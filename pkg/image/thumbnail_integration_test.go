//go:build integration

package image

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

func TestJXLThumbnail(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("/wpBBgATiAIAwAC1nyAAABUqo4wbvJzr+fJDh8W0jesMbbVtYQljs70wSEg4g42LBN8cpHHpBiSikgQ=")
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "image.jxl")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Fatal(err)
	}

	encoder := ThumbnailEncoder{FFMpeg: ffmpeg.NewEncoder(ffmpegPath)}
	thumbnail, err := encoder.GetThumbnail(&models.ImageFile{
		BaseFile: &models.BaseFile{Path: path, Basename: "image.jxl"},
		Format:   "jpegxl",
		Width:    8,
		Height:   8,
	}, 4)
	if err != nil {
		t.Fatal(err)
	}

	config, err := jpeg.DecodeConfig(bytes.NewReader(thumbnail))
	if err != nil {
		t.Fatal(err)
	}
	if config.Width != 4 || config.Height != 4 {
		t.Fatalf("thumbnail dimensions = %dx%d, want 4x4", config.Width, config.Height)
	}
}
