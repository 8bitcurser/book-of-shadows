package storage

import (
	"book-of-shadows/models"
	"bytes"
	"compress/gzip"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type CookiesConfig struct {
	Prefix string
	MaxAge int
}

func NewInvestigatorCookieConfig() *CookiesConfig {
	return &CookiesConfig{
		Prefix: "investigator",
		MaxAge: 3600 * 24 * 30,
	}
}

func (c *CookiesConfig) SaveInvestigatorCookie(w http.ResponseWriter, inv *models.Investigator) {
	id := fmt.Sprintf(
		"%d_%s", time.Now().Unix(), strings.ToLower(strings.ReplaceAll(inv.Name, " ", "_")),
	)
	id = c.Prefix + id
	inv.ID = id
	data, err := inv.ToJSON()
	if err != nil {
		fmt.Errorf("Failed to marshal investigator %s", err)
	}
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	_, err = writer.Write(data)
	if err != nil {
		log.Printf("Failed to compress data: %s", err)
		return
	}
	writer.Close()
	encodedValue := base64.URLEncoding.EncodeToString(compressed.Bytes())
	cookie := &http.Cookie{
		Name:     id,
		Value:    encodedValue,
		Path:     "/",
		MaxAge:   c.MaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func (c *CookiesConfig) GetInvestigatorCookie(r *http.Request, id string) (*models.Investigator, error) {
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(cookie.Name, c.Prefix) {
			if cookie.Name == id {

				data, err := base64.URLEncoding.DecodeString(cookie.Value)
				if err != nil {
					continue // Skip invalid cookies
				}
				reader, err := gzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, fmt.Errorf("failed to create gzip reader: %w", err)
				}
				defer reader.Close()

				decompressed, err := io.ReadAll(reader)
				if err != nil {
					return nil, fmt.Errorf("failed to decompress data: %w", err)
				}
				var character models.Investigator
				if err := json.Unmarshal(decompressed, &character); err != nil {
					continue // Skip invalid character data
				}
				return &character, nil
			}
		}
	}
	return nil, nil
}

func (c *CookiesConfig) DeleteInvestigatorCookie(w http.ResponseWriter, id string) {
	cookieName := id
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
}

func (c *CookiesConfig) ListInvestigators(r *http.Request) (map[string]*models.Investigator, error) {
	characters := make(map[string]*models.Investigator)
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(cookie.Name, c.Prefix) {
			data, err := base64.URLEncoding.DecodeString(cookie.Value)
			if err != nil {
				continue // Skip invalid cookies
			}

			reader, err := gzip.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, fmt.Errorf("failed to create gzip reader: %w", err)
			}
			defer reader.Close()

			decompressed, err := io.ReadAll(reader)
			if err != nil {
				return nil, fmt.Errorf("failed to decompress data: %w", err)
			}

			var character models.Investigator
			if err := json.Unmarshal(decompressed, &character); err != nil {
				continue // Skip invalid character data
			}
			characters[cookie.Name] = &character
		}
	}
	return characters, nil
}

func (c *CookiesConfig) ExportInvestigatorsList(r *http.Request) (string, error) {
	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(cookie.Name, c.Prefix) {
			cookies[cookie.Name] = cookie.Value

		}
	}
	data, err := json.Marshal(cookies)
	if err != nil {
		fmt.Errorf("Failed to marshal investigators %s", err)
	}

	encodedValue := base32.StdEncoding.EncodeToString(data)
	return encodedValue, nil

}

func (c *CookiesConfig) ImportInvestigatorsList(w http.ResponseWriter, encodedData string) error {
	data, err := base32.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return err
	}

	var cookies map[string]string
	if err = json.Unmarshal(data, &cookies); err != nil {

		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Set cookies
	for name, value := range cookies {

		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     "/",
			MaxAge:   c.MaxAge,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
	}
	return nil
}

func (c *CookiesConfig) UpdateInvestigatorCookie(w http.ResponseWriter, id string, inv *models.Investigator) {
	data, err := inv.ToJSON()
	if err != nil {
		fmt.Errorf("Failed to marshal investigator %s", err)
	}
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	_, err = writer.Write(data)
	if err != nil {
		log.Printf("Failed to compress data: %s", err)
		return
	}
	writer.Close()
	encodedValue := base64.URLEncoding.EncodeToString(compressed.Bytes())
	cookie := &http.Cookie{
		Name:     id,
		Value:    encodedValue,
		Path:     "/",
		MaxAge:   c.MaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

}
