package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"book-of-shadows/internal/config"
	"book-of-shadows/internal/errors"
	"book-of-shadows/models"
)

// CookieStore handles investigator storage using browser cookies
type CookieStore struct {
	config      *config.CookieConfig
	exportStore ExportStore
}

// NewCookieStore creates a new CookieStore instance
func NewCookieStore(cfg *config.CookieConfig, exportStore ExportStore) *CookieStore {
	return &CookieStore{
		config:      cfg,
		exportStore: exportStore,
	}
}

// SaveInvestigator saves an investigator to a cookie
func (cs *CookieStore) SaveInvestigator(w http.ResponseWriter, inv *models.Investigator) (string, error) {
	if inv == nil {
		return "", errors.ErrInvalidData
	}

	// Generate ID
	id := cs.generateInvestigatorID(inv.Name)
	inv.ID = id

	// Encode investigator data
	encodedValue, err := cs.encodeInvestigator(inv)
	if err != nil {
		return "", fmt.Errorf("failed to encode investigator: %w", err)
	}

	// Check cookie size limit (4KB)
	if len(encodedValue) > 4096 {
		return "", errors.ErrCookieTooLarge
	}

	// Create and set cookie
	cookie := cs.createCookie(id, encodedValue)
	http.SetCookie(w, cookie)

	return id, nil
}

// GetInvestigator retrieves an investigator from a cookie
func (cs *CookieStore) GetInvestigator(r *http.Request, id string) (*models.Investigator, error) {
	if id == "" {
		return nil, errors.ErrInvalidData
	}

	cookie, err := r.Cookie(id)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, errors.ErrCookieNotFound
		}
		return nil, fmt.Errorf("failed to get cookie: %w", err)
	}

	investigator, err := cs.decodeInvestigator(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode investigator: %w", err)
	}

	return investigator, nil
}

// UpdateInvestigator updates an existing investigator cookie
func (cs *CookieStore) UpdateInvestigator(w http.ResponseWriter, id string, inv *models.Investigator) error {
	if id == "" || inv == nil {
		return errors.ErrInvalidData
	}

	encodedValue, err := cs.encodeInvestigator(inv)
	if err != nil {
		return fmt.Errorf("failed to encode investigator: %w", err)
	}

	if len(encodedValue) > 4096 {
		return errors.ErrCookieTooLarge
	}

	cookie := cs.createCookie(id, encodedValue)
	http.SetCookie(w, cookie)

	return nil
}

// DeleteInvestigator removes an investigator cookie
func (cs *CookieStore) DeleteInvestigator(w http.ResponseWriter, id string) error {
	if id == "" {
		return errors.ErrInvalidData
	}

	cookie := &http.Cookie{
		Name:     id,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: cs.config.HttpOnly,
		Secure:   cs.config.Secure,
		SameSite: http.SameSite(cs.config.SameSite),
	}

	http.SetCookie(w, cookie)
	return nil
}

// ListInvestigators returns all investigators stored in cookies
func (cs *CookieStore) ListInvestigators(r *http.Request) (map[string]*models.Investigator, error) {
	investigators := make(map[string]*models.Investigator)

	for _, cookie := range r.Cookies() {
		if !strings.HasPrefix(cookie.Name, cs.config.Prefix) {
			continue
		}

		investigator, err := cs.decodeInvestigator(cookie.Value)
		if err != nil {
			// Skip invalid cookies instead of failing completely
			continue
		}

		investigators[cookie.Name] = investigator
	}

	return investigators, nil
}

// ExportInvestigatorsList exports all investigators for sharing
func (cs *CookieStore) ExportInvestigatorsList(r *http.Request) (string, error) {
	cookies := make(map[string]string)

	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(cookie.Name, cs.config.Prefix) {
			cookies[cookie.Name] = cookie.Value
		}
	}

	if len(cookies) == 0 {
		return "", errors.ErrNotFound
	}

	data, err := json.Marshal(cookies)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cookies: %w", err)
	}

	encodedValue := base32.StdEncoding.EncodeToString(data)

	// Save to export store
	uuid, err := cs.exportStore.SaveExport(encodedValue)
	if err != nil {
		return "", fmt.Errorf("failed to save export: %w", err)
	}

	return uuid, nil
}

// ImportInvestigatorsList imports investigators from a shared export
func (cs *CookieStore) ImportInvestigatorsList(w http.ResponseWriter, uuid string) error {
	if uuid == "" {
		return errors.ErrInvalidData
	}

	encodedData, err := cs.exportStore.GetExport(uuid)
	if err != nil {
		return fmt.Errorf("failed to get export: %w", err)
	}

	data, err := base32.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	var cookies map[string]string
	if err := json.Unmarshal(data, &cookies); err != nil {
		return fmt.Errorf("failed to unmarshal cookies: %w", err)
	}

	for name, value := range cookies {
		cookie := cs.createCookie(name, value)
		http.SetCookie(w, cookie)
	}

	return nil
}

// Helper methods

// generateInvestigatorID creates a unique ID for an investigator
func (cs *CookieStore) generateInvestigatorID(name string) string {
	timestamp := time.Now().Unix()
	safeName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	return fmt.Sprintf("%s%d_%s", cs.config.Prefix, timestamp, safeName)
}

// createCookie creates a properly configured HTTP cookie
func (cs *CookieStore) createCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   cs.config.MaxAge,
		HttpOnly: cs.config.HttpOnly,
		Secure:   cs.config.Secure,
		SameSite: http.SameSite(cs.config.SameSite),
	}
}

// encodeInvestigator compresses and encodes investigator data
func (cs *CookieStore) encodeInvestigator(inv *models.Investigator) (string, error) {
	// Marshal to JSON
	data, err := inv.ToJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal investigator: %w", err)
	}

	// Compress with gzip
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(data); err != nil {
		gzWriter.Close()
		return "", fmt.Errorf("failed to compress data: %w", err)
	}
	if err := gzWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Encode with base64
	encoded := base64.URLEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

// decodeInvestigator decodes and decompresses investigator data
func (cs *CookieStore) decodeInvestigator(encodedData string) (*models.Investigator, error) {
	// Decode from base64
	data, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Decompress from gzip
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	// Unmarshal from JSON
	var investigator models.Investigator
	if err := json.Unmarshal(decompressed, &investigator); err != nil {
		return nil, fmt.Errorf("failed to unmarshal investigator: %w", err)
	}

	// Populate SpecialArchetypeRules from the Archetypes map (not serialized with json:"-")
	if investigator.Archetype != nil && investigator.Archetype.Name != "" {
		if archetype, exists := models.Archetypes[investigator.Archetype.Name]; exists {
			investigator.Archetype.SpecialArchetypeRules = archetype.SpecialArchetypeRules
		}
	}

	return &investigator, nil
}