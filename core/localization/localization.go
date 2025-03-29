package localization

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v2"
)

type Localization struct {
	defaultLang string
	langData    map[string]map[string]string
}

// New creates a new Localization instance
func New(defaultLang string) *Localization {
	return &Localization{
		defaultLang: defaultLang,
		langData:    make(map[string]map[string]string),
	}
}

// LoadLanguageFile loads translations from a JSON or YAML file
func (l *Localization) LoadLanguageFile(lang, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read language file: %v", err)
	}

	var translations map[string]string

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse JSON: %v", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse YAML: %v", err)
		}
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	l.langData[lang] = translations
	return nil
}

// GetText returns the translated text for the given key
func (l *Localization) GetText(lang, key string) string {
	// Try requested language first
	if translations, ok := l.langData[lang]; ok {
		if text, ok := translations[key]; ok {
			return text
		}
	}

	// Fall back to default language
	if translations, ok := l.langData[l.defaultLang]; ok {
		if text, ok := translations[key]; ok {
			return text
		}
	}

	// Return key if no translation found
	return key
}

// Middleware creates a Fiber middleware for language selection
func (l *Localization) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get language from query parameter
		lang := c.Query("lang")

		// If not in query, try Accept-Language header
		if lang == "" {
			acceptLang := c.Get("Accept-Language")
			if acceptLang != "" {
				// Parse the Accept-Language header and get the first language
				lang = strings.Split(strings.Split(acceptLang, ",")[0], "-")[0]
			}
		}

		// If still no language, use default
		if lang == "" {
			lang = l.defaultLang
		}

		// Store language in locals for use in handlers
		c.Locals("lang", lang)

		return c.Next()
	}
}
