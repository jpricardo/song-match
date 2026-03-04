package ytbutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DownloadTrack(videoURL string) ([]byte, string, string, error) {
	tmpDir, err := os.MkdirTemp("", "ytdl-*")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outTemplate := filepath.Join(tmpDir, "audio.%(ext)s")
	expectedWavPath := filepath.Join(tmpDir, "audio.wav")

	cmd := exec.Command("yt-dlp",
		"-4",
		"--no-playlist",
		"--print", "%(title)s|||%(thumbnail)s",
		"--no-simulate",
		"-f", "bestaudio",
		"-x",
		"--audio-format", "wav",
		"-o", outTemplate,
		videoURL,
	)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, "", "", fmt.Errorf("yt-dlp failed: %v\nLogs: %s", err, errBuf.String())
	}

	output := strings.TrimSpace(outBuf.String())
	parts := strings.Split(output, "|||")

	title := parts[0]
	thumbnailURL := ""
	if len(parts) > 1 {
		thumbnailURL = parts[1]
	}

	wavData, err := os.ReadFile(expectedWavPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to read output WAV: %w", err)
	}

	return wavData, title, thumbnailURL, nil
}
