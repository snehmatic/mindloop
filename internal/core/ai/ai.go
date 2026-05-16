package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

const (
	SettingKeyAIProvider = "ai_provider"
	SettingKeyAIModel    = "ai_model"
	SettingKeyAIToken    = "ai_token"
	SettingKeyAIBaseURL  = "ai_base_url"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// GetSettings retrieves the AI configuration from the database
func (s *Service) GetSettings() (provider, model, token, baseURL string, err error) {
	var pSetting, mSetting, tSetting, bSetting models.AppSetting
	s.DB.Where("key = ?", SettingKeyAIProvider).Limit(1).Find(&pSetting)
	s.DB.Where("key = ?", SettingKeyAIModel).Limit(1).Find(&mSetting)
	s.DB.Where("key = ?", SettingKeyAIToken).Limit(1).Find(&tSetting)
	s.DB.Where("key = ?", SettingKeyAIBaseURL).Limit(1).Find(&bSetting)

	provider = pSetting.Value
	model = mSetting.Value
	baseURL = bSetting.Value

	// Token from DB overrides env var, if exists
	envToken := os.Getenv("MINDLOOP_AI_TOKEN")
	if tSetting.Value != "" {
		decrypted, err := utils.Decrypt(tSetting.Value)
		if err == nil {
			token = strings.TrimSpace(decrypted)
		}
	} else if envToken != "" {
		token = strings.TrimSpace(envToken)
	}

	if provider == "" {
		provider = "gemini" // default
	}
	return provider, model, token, baseURL, nil
}

// SaveSettings encrypts the token and saves the configuration
func (s *Service) SaveSettings(provider, model, token, baseURL string) error {
	s.saveOrUpdate(SettingKeyAIProvider, provider)
	s.saveOrUpdate(SettingKeyAIModel, model)
	s.saveOrUpdate(SettingKeyAIBaseURL, baseURL)

	if token == "__CLEAR__" {
		s.saveOrUpdate(SettingKeyAIToken, "")
	} else if token != "" {
		encrypted, err := utils.Encrypt(token)
		if err != nil {
			return err
		}
		s.saveOrUpdate(SettingKeyAIToken, encrypted)
	}
	return nil
}

func (s *Service) saveOrUpdate(key, value string) {
	var setting models.AppSetting
	result := s.DB.Where("key = ?", key).Limit(1).Find(&setting)
	if result.RowsAffected == 0 {
		s.DB.Create(&models.AppSetting{Key: key, Value: value})
	} else {
		setting.Value = value
		s.DB.Save(&setting)
	}
}

// ListModels fetches the available models for the configured provider
func (s *Service) ListModels() ([]string, error) {
	provider, _, token, baseURL, _ := s.GetSettings()
	if token == "" && provider != "custom" {
		return nil, fmt.Errorf("AI token not configured")
	}

	switch provider {
	case "openai", "custom":
		if provider == "custom" && baseURL == "" {
			return nil, fmt.Errorf("base URL is required for custom local providers")
		}
		return s.listOpenAIModels(token, baseURL)
	case "anthropic":
		return nil, fmt.Errorf("anthropic support coming soon")
	default:
		return s.listGeminiModels(token)
	}
}

func (s *Service) listGeminiModels(token string) ([]string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, formatAPIError("Gemini", resp.StatusCode, body)
	}

	var result struct {
		Models []struct {
			Name                       string   `json:"name"`
			SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Models {
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				models = append(models, strings.TrimPrefix(m.Name, "models/"))
				break
			}
		}
	}
	return models, nil
}

func (s *Service) listOpenAIModels(token, baseURL string) ([]string, error) {
	url := "https://api.openai.com/v1/models"
	if baseURL != "" {
		url = strings.TrimSuffix(baseURL, "/") + "/models"
	}
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, formatAPIError("OpenAI", resp.StatusCode, body)
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Data {
		if strings.HasPrefix(m.ID, "gpt-") || strings.HasPrefix(m.ID, "o1-") {
			models = append(models, m.ID)
		}
	}
	return models, nil
}

// TestConnection sends a minimal prompt to verify the configuration works
func (s *Service) TestConnection(provider, model, token, baseURL string) error {
	if token == "" && provider != "custom" {
		return fmt.Errorf("AI token not configured")
	}

	testData := `{"test": true}`

	var err error
	switch provider {
	case "openai", "custom":
		if provider == "custom" && baseURL == "" {
			return fmt.Errorf("base URL is required for custom local providers")
		}
		_, err = s.generateOpenAI(model, token, testData, baseURL)
	case "anthropic":
		_, err = s.generateAnthropic(model, token, testData)
	default:
		_, err = s.generateGemini(model, token, testData)
	}
	return err
}

func (s *Service) GenerateJournal(summary models.SummaryReport) (string, error) {
	provider, model, token, baseURL, _ := s.GetSettings()
	if token == "" && provider != "custom" {
		return "", fmt.Errorf("AI token not configured. Set MINDLOOP_AI_TOKEN or configure via UI settings")
	}

	dataBytes, err := json.Marshal(summary)
	if err != nil {
		return "", err
	}

	switch provider {
	case "openai", "custom":
		if provider == "custom" && baseURL == "" {
			return "", fmt.Errorf("base URL is required for custom local providers")
		}
		return s.generateOpenAI(model, token, string(dataBytes), baseURL)
	case "anthropic":
		return s.generateAnthropic(model, token, string(dataBytes))
	default:
		// Default to gemini format
		return s.generateGemini(model, token, string(dataBytes))
	}
}

func (s *Service) generateGemini(model, token, contextData string) (string, error) {
	if model == "" {
		model = "gemini-1.5-flash-latest"
	}
	model = strings.TrimPrefix(model, "models/")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, token)

	reqBody := map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": map[string]interface{}{"text": JournalSystemPrompt},
		},
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": "Here is my activity summary data:\n" + contextData},
				},
			},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", formatAPIError("Gemini", resp.StatusCode, body)
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no content generated")
}

func (s *Service) generateOpenAI(model, token, contextData, baseURL string) (string, error) {
	if model == "" {
		model = "gpt-4o-mini"
	}
	url := "https://api.openai.com/v1/chat/completions"
	if baseURL != "" {
		url = strings.TrimSuffix(baseURL, "/") + "/chat/completions"
	}

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{"role": "system", "content": JournalSystemPrompt},
			{"role": "user", "content": "Here is my activity summary data:\n" + contextData},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", formatAPIError("OpenAI", resp.StatusCode, body)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no content generated")
}

// Stub for Anthropic to prevent compile errors, implemented simply
func (s *Service) generateAnthropic(model, token, contextData string) (string, error) {
	// Simple stub for v1
	return "", fmt.Errorf("anthropic support coming soon")
}

func formatAPIError(provider string, statusCode int, body []byte) error {
	var errResp struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return fmt.Errorf("%s API error: %s", provider, errResp.Error.Message)
	}
	return fmt.Errorf("%s API error (%d)", provider, statusCode)
}
