package dl

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/zuchzub/Go/pkg/config"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Laky-64/gologging"
)

const (
	defaultRequestTimeout = 30 * time.Second
	defaultConnectTimeout = 10 * time.Second
	maxRetries            = 2
	initialBackoff        = 1 * time.Second
)

var client = &http.Client{
	Timeout: defaultRequestTimeout,
	Transport: &http.Transport{
		TLSHandshakeTimeout:   defaultConnectTimeout,
		ResponseHeaderTimeout: defaultRequestTimeout,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
	},
}

// sendRequest performs an HTTP request with a given context, method, URL, body, and headers.
// It includes retry logic with exponential backoff for temporary network errors and server-side issues.
// It returns an HTTP response or an error if the request fails after all retries.
func sendRequest(ctx context.Context, method, fullURL string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "*/*")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	var reqErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}

		resp, reqErr = client.Do(req)
		if reqErr == nil {
			if resp.StatusCode < 500 {
				return resp, nil // Success
			}
			if err := resp.Body.Close(); err != nil {
				gologging.WarnF("failed to close response body: %v", err)
			}
			reqErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		} else if isTemporaryError(reqErr) {
			gologging.InfoF("Temporary error on attempt %d/%d: %v", attempt+1, maxRetries, reqErr)
			continue // Retry on temporary errors
		} else {
			break // Do not retry on permanent errors
		}
	}

	if reqErr == nil {
		reqErr = fmt.Errorf("request failed after %d attempts", maxRetries)
	}

	return nil, fmt.Errorf("request failed: %w", reqErr)
}

// isTemporaryError determines if an error is temporary and thus worth retrying.
// It returns true for network timeouts and temporary operational errors.
func isTemporaryError(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}
	return false
}

// generateUniqueName creates a pseudo-random filename using a combination of the current timestamp and a random number.
// It takes a file extension and returns a unique filename.
func generateUniqueName(ext string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999))
	return fmt.Sprintf("%d_%05d%s", time.Now().UnixNano(), n.Int64(), ext)
}

// determineFilename safely determines a valid filename for a download.
// It prioritizes the Content-Disposition header, falls back to the URL path, and generates a unique name if neither is available.
// It returns a secure and sanitized filename.
func determineFilename(urlStr, contentDisp string) string {
	if filename := extractFilename(contentDisp); filename != "" {
		return filepath.Join(config.Conf.DownloadsDir, sanitizeFilename(filename))
	}

	if parsedURL, err := url.Parse(urlStr); err == nil {
		filename := path.Base(parsedURL.Path)
		if filename != "" && filename != "/" && !strings.Contains(filename, "?") {
			return filepath.Join(config.Conf.DownloadsDir, sanitizeFilename(filename))
		}
	}

	return filepath.Join(config.Conf.DownloadsDir, generateUniqueName(".tmp"))
}

// writeToFile writes data from an io.Reader to a specified file.
// It returns an error if file creation or writing fails.
func writeToFile(filename string, data io.Reader) error {
	// #nosec G304 - This is a security risk if the filename is not properly sanitized.
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create the file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, data); err != nil {
		return fmt.Errorf("failed to write to the file: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from a URL and saves it to a local path.
// It supports overwriting existing files and determines the filename automatically if not provided.
// It returns the final file path or an error if the download fails.
func DownloadFile(ctx context.Context, urlStr, fileName string, overwrite bool) (string, error) {
	if urlStr == "" {
		return "", errors.New("an empty URL was provided")
	}

	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create the request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("the request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code received: %d", resp.StatusCode)
	}

	if fileName == "" {
		fileName = determineFilename(urlStr, resp.Header.Get("Content-Disposition"))
	}

	if !overwrite {
		if _, err := os.Stat(fileName); err == nil {
			return fileName, nil // File already exists, no need to download again.
		}
	}

	if err := os.MkdirAll(filepath.Dir(fileName), defaultDownloadDirPerm); err != nil {
		return "", fmt.Errorf("failed to create the directory: %w", err)
	}

	// Download to a temporary .part file to ensure atomicity.
	tempPath := fileName + ".part"
	if err := writeToFile(tempPath, resp.Body); err != nil {
		return "", err
	}

	// Rename the temporary file to its final name upon successful download.
	if err := os.Rename(tempPath, fileName); err != nil {
		return "", fmt.Errorf("failed to rename the temporary file: %w", err)
	}

	return fileName, nil
}
