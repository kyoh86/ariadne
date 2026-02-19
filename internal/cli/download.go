package cli

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func downloadWithOptionalSHA512(srcURL, dstPath, expectedSHA512 string) error {
	client := &http.Client{Timeout: 2 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, srcURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("download failed: %s", res.Status)
	}

	tmp := dstPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	hasher := sha512.New()
	writer := io.MultiWriter(f, hasher)
	if _, err := io.Copy(writer, res.Body); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if expectedSHA512 != "" {
		got := hex.EncodeToString(hasher.Sum(nil))
		if !strings.EqualFold(got, expectedSHA512) {
			_ = os.Remove(tmp)
			return fmt.Errorf("sha512 mismatch: got %s expected %s", got, expectedSHA512)
		}
	}

	return os.Rename(tmp, dstPath)
}
