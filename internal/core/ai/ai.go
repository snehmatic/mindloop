package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

const (
	SettingKeyAIProvider = "ai_provider"
	SettingKeyAIModel    = "ai_model"
	SettingKeyAIToken    = "ai_token"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// GetSettings retrieves the AI configuration from the database
func (s *Service) GetSettings() (provider, model, token string, err error) {
	var pSetting, mSetting, tSetting models.AppSetting
	s.DB.Where("key = ?", SettingKeyAIProvider).Limit(1).Find(&pSetting)
	s.DB.Where("key = ?", SettingKeyAIModel).Limit(1).Find(&mSetting)
	s.DB.Where("key = ?", SettingKeyAIToken).Limit(1).Find(&tSetting)

	provider = pSetting.Value
	model = mSetting.Value

	// Token from DB overrides env var, if exists
	envToken := os.Getenv("MINDLOOP_AI_TOKEN")
	if tSetting.Value != "" {
		decrypted, err := utils.Decrypt(tSetting.Value)
		if err == nil {
			token = decrypted
		}
	} else if envToken != "" {
		token = envToken
	}

	if provider == "" {
		provider = "gemini" // default
	}
	return provider, model, token, nil
}

// SaveSettings encrypts the token and saves the configuration
func (s *Service) SaveSettings(provider, model, token string) error {
	s.saveOrUpdate(SettingKeyAIProvider, provider)
	s.saveOrUpdate(SettingKeyAIModel, model)
	
	if token != "" {
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

// GenerateJournal generates a journal entry based on the summary report
func (s *Service) GenerateJournal(summary models.SummaryReport) (string, error) {
	provider, model, token, _ := s.GetSettings()
	if token == "" {
		return "", fmt.Errorf("AI token not configured. Set MINDLOOP_AI_TOKEN or configure via UI settings")
	}

	dataBytes, err := json.Marshal(summary)
	if err != nil {
		return "", err
	}

	switch provider {
	case "openai":
		return s.generateOpenAI(model, token, string(dataBytes))
	case "anthropic":
		return s.generateAnthropic(model, token, string(dataBytes))
	default:
		// Default to gemini format
		return s.generateGemini(model, token, string(dataBytes))
	}
}

func (s *Service) generateGemini(model, token, contextData string) (string, error) {
	if model == "" {
		model = "gemini-1.5-flash"
	}
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
		return "", fmt.Errorf("gemini API error (%d): %s", resp.StatusCode, string(body))
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

func (s *Service) generateOpenAI(model, token, contextData string) (string, error) {
	if model == "" {
		model = "gpt-4o-mini"
	}
	url := "https://api.openai.com/v1/chat/completions"
	
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
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openAI API error (%d): %s", resp.StatusCode, string(body))
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
