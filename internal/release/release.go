// Package release fetches gitt release artifacts from GitHub.
package release

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	repoOwner   = "foreverfl"
	repoName    = "gitt"
	binName     = "gitt"
	httpTimeout = 30 * time.Second
)

// LatestTag returns the tag name of the latest non-draft, non-prerelease
// release on GitHub.
func LatestTag() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api status %d", resp.StatusCode)
	}
	var body struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode release: %w", err)
	}
	if body.TagName == "" {
		return "", errors.New("empty tag_name in release")
	}
	return body.TagName, nil
}

// Download streams the release tarball for tag and writes the gitt
// binary to outPath with mode 0755. outPath should sit on the same
// filesystem as the final install location so the caller can atomically
// os.Rename it onto the running binary.
func Download(tag, outPath string) error {
	asset := fmt.Sprintf("%s_%s_%s.tar.gz", binName, runtime.GOOS, runtime.GOARCH)
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", repoOwner, repoName, tag, asset)

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download %s: %w", asset, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: status %d", asset, resp.StatusCode)
	}

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("tar: %w", err)
		}
		if header.Typeflag != tar.TypeReg || filepath.Base(header.Name) != binName {
			continue
		}
		f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			return fmt.Errorf("create %s: %w", outPath, err)
		}
		if _, err := io.Copy(f, tr); err != nil {
			_ = f.Close()
			_ = os.Remove(outPath)
			return fmt.Errorf("extract: %w", err)
		}
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("binary %s not found in tarball", binName)
}
