# AI-Powered Journal Generation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement AI-powered journal generation using pre-aggregated summary data, with secure token storage and support for multiple LLM providers.

**Architecture:** We will introduce an `AppSetting` model for secure token storage in the DB, an `ai` package in `internal/core` for LLM interaction, and new CLI commands and Web handlers to trigger the generation based on the `SummaryReport` data.

**Tech Stack:** Go (net/http), GORM (SQLite/Postgres), HTML templates.

---

### Task 1: Data Models & Database Migration

**Files:**
- Modify: `models/types.go`
- Modify: `db/db.go`

- [ ] **Step 1: Add `AppSetting` model**
Add the `AppSetting` model to `models/types.go`:
```go
// AppSetting represents a key-value store for application settings (like encrypted AI tokens)
type AppSetting struct {
	gorm.Model
	Key   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"` // encrypted if sensitive
}
```

- [ ] **Step 2: Update database migrations**
Add `&models.AppSetting{}` to the `AutoMigrate` call in `db/db.go`'s `MigrateDB` function.

- [ ] **Step 3: Commit**
```bash
git add models/types.go db/db.go
git commit -m "feat: add AppSetting model for secure configuration storage"
```

### Task 2: Crypto Utilities

**Files:**
- Create: `internal/utils/crypto.go`
- Create: `internal/utils/crypto_test.go`

- [ ] **Step 1: Implement AES-GCM encryption/decryption**
Create `internal/utils/crypto.go`. Use a hardcoded fallback salt/key for local installations, but allow overriding via environment variable `MINDLOOP_ENC_KEY`.

```go
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func getEncryptionKey() []byte {
	key := os.Getenv("MINDLOOP_ENC_KEY")
	if key == "" {
		key = "mindloop-default-secret-key-32b!" // 32 bytes
	}
	// ensure 32 bytes
	if len(key) < 32 {
		padding := make([]byte, 32-len(key))
		key = key + string(padding)
	}
	return []byte(key[:32])
}

// Encrypt encrypts a string using AES-GCM and returns a base64 encoded string
func Encrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64 encoded string using AES-GCM
func Decrypt(cryptoText string) (string, error) {
	if cryptoText == "" {
		return "", nil
	}
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
```

- [ ] **Step 2: Write tests for crypto**
Create `internal/utils/crypto_test.go`:
```go
package utils

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	original := "my-secret-api-token"
	encrypted, err := Encrypt(original)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	if encrypted == original {
		t.Fatalf("Encrypted text is same as original")
	}
	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	if decrypted != original {
		t.Fatalf("Expected %s, got %s", original, decrypted)
	}
}
```

- [ ] **Step 3: Run tests and verify**
Run: `go test ./internal/utils -v`
Expected: PASS

- [ ] **Step 4: Commit**
```bash
git add internal/utils/crypto.go internal/utils/crypto_test.go
git commit -m "feat: add crypto utilities for secure token storage"
```

### Task 3: AI Service (Core Logic)

**Files:**
- Create: `internal/core/ai/ai.go`
- Create: `internal/core/ai/prompts.go`

- [ ] **Step 1: Create prompt definitions**
Create `internal/core/ai/prompts.go`:
```go
package ai

const JournalSystemPrompt = `You are Mindloop's analytical, encouraging, and ADHD-friendly AI assistant.
Your goal is to review the user's raw JSON activity data (focus sessions, habits, tasks, side quests, and points) and generate a cohesive, reflective journal entry.

Guidelines:
1. Be concise but comprehensive. Highlight key wins and identify patterns.
2. Structure the summary with clear headers or bullet points for readability.
3. Maintain an encouraging and objective tone.
4. Do not expose raw JSON or IDs to the user. Translate the data into a human-readable narrative.
5. End with a short reflective thought or question to help the user plan their next steps.`
```

- [ ] **Step 2: Create AI interface and settings manager**
Create `internal/core/ai/ai.go`:
```go
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
	s.DB.Where("key = ?", SettingKeyAIProvider).First(&pSetting)
	s.DB.Where("key = ?", SettingKeyAIModel).First(&mSetting)
	s.DB.Where("key = ?", SettingKeyAIToken).First(&tSetting)

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
	return
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
	if s.DB.Where("key = ?", key).First(&setting).Error == gorm.ErrRecordNotFound {
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

	if provider == "openai" {
		return s.generateOpenAI(model, token, string(dataBytes))
	} else if provider == "anthropic" {
		return s.generateAnthropic(model, token, string(dataBytes))
	}
	// Default to gemini format (which is also supported by Ollama if configured right, or we can add ollama specific later)
	return s.generateGemini(model, token, string(dataBytes))
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error (%d): %s", resp.StatusCode, string(body))
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
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
```

- [ ] **Step 3: Commit**
```bash
git add internal/core/ai/
git commit -m "feat: implement AI service and prompt for journal generation"
```

### Task 4: CLI Implementation

**Files:**
- Modify: `cmd/cli/journal.go`
- Modify: `cmd/cli/summary.go`

- [ ] **Step 1: Add generate command to `cmd/cli/journal.go`**
Add the generate command which calls `summaryService.GenerateSummary` and `aiService.GenerateJournal`.

```go
// Add to cmd/cli/journal.go (with appropriate imports like "github.com/snehmatic/mindloop/internal/core/ai" and "github.com/snehmatic/mindloop/internal/core/summary")

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Auto-generate a journal entry using AI",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.GetConfig()
		database, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database: " + err.Error())
			return
		}

		daily, _ := cmd.Flags().GetBool("daily")
		weekly, _ := cmd.Flags().GetBool("weekly")
		yearly, _ := cmd.Flags().GetBool("yearly")

		start, end := utils.GetDateRange("daily") // default
		if weekly {
			start, end = utils.GetDateRange("weekly")
		} else if yearly {
			start, end = utils.GetDateRange("yearly")
		}

		summaryService := summary.NewService(database)
		report, err := summaryService.GenerateSummary(start, end)
		if err != nil {
			utils.PrintErrorln("Failed to generate summary report: " + err.Error())
			return
		}

		utils.PrintInfoln("✨ Generating AI journal entry...")
		aiService := ai.NewService(database)
		generatedText, err := aiService.GenerateJournal(report)
		if err != nil {
			utils.PrintErrorln("Failed to generate journal: " + err.Error())
			return
		}

		fmt.Println("\n" + generatedText + "\n")

		// Interactive save
		fmt.Print("Would you like to save this into the journal? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			journalService := journal.NewService(database)
			title := fmt.Sprintf("AI Summary: %s", report.DateRange)
			err = journalService.CreateJournalEntry(title, generatedText, "reflective")
			if err != nil {
				utils.PrintErrorln("Failed to save journal: " + err.Error())
			} else {
				utils.PrintSuccessln("Saved successfully!")
			}
		}
	},
}

func init() {
	// In the existing init() function of journal.go:
	generateCmd.Flags().BoolP("daily", "d", false, "Generate daily summary")
	generateCmd.Flags().BoolP("weekly", "w", false, "Generate weekly summary")
	generateCmd.Flags().BoolP("yearly", "y", false, "Generate yearly summary")
	journalCmd.AddCommand(generateCmd)
	// ... rest of init
}
```

- [ ] **Step 2: Add recommendation to `cmd/cli/summary.go`**
In `cmd/cli/summary.go`, at the end of the `summaryCmd` Run function, add:
```go
utils.PrintInfoln("\n💡 Tip: For an AI-generated overview, use `mindloop journal generate -d` (or -w, -y)")
```

- [ ] **Step 3: Commit**
```bash
git add cmd/cli/journal.go cmd/cli/summary.go
git commit -m "feat: add CLI support for generating and saving AI journals"
```

### Task 5: Web API & Handlers

**Files:**
- Modify: `api/v1/handlers.go`
- Create: `api/v1/handlers_ai.go` (if needed, or add to `handlers_impl.go`)

- [ ] **Step 1: Add AI API handlers**
Create `api/v1/handlers_ai.go`:
```go
package v1

import (
	"encoding/json"
	"net/http"

	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/utils"
)

type AISettingsRequest struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Token    string `json:"token"`
}

func (h *Handler) GetAISettings(w http.ResponseWriter, r *http.Request) {
	aiService := ai.NewService(h.DB)
	provider, model, token, _ := aiService.GetSettings()
	
	// Don't send token back to UI, just a boolean indicator if it's set
	hasToken := token != ""
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"provider": provider,
		"model":    model,
		"hasToken": hasToken,
	})
}

func (h *Handler) SaveAISettings(w http.ResponseWriter, r *http.Request) {
	var req AISettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	aiService := ai.NewService(h.DB)
	if err := aiService.SaveSettings(req.Provider, req.Model, req.Token); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to save settings: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Settings saved successfully"})
}

func (h *Handler) GenerateAIJournal(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	start, end := utils.GetDateRange(period)
	summaryService := summary.NewService(h.DB)
	report, err := summaryService.GenerateSummary(start, end)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate summary data")
		return
	}

	aiService := ai.NewService(h.DB)
	generatedText, err := aiService.GenerateJournal(report)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "AI Generation failed: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"content": generatedText,
		"title":   "AI Summary: " + report.DateRange,
	})
}
```

- [ ] **Step 2: Register routes in `api/v1/handlers.go`**
Add to `RegisterRoutes`:
```go
	router.HandleFunc("/api/v1/ai/settings", h.GetAISettings).Methods("GET")
	router.HandleFunc("/api/v1/ai/settings", h.SaveAISettings).Methods("POST")
	router.HandleFunc("/api/v1/ai/generate", h.GenerateAIJournal).Methods("GET")
```

- [ ] **Step 3: Commit**
```bash
git add api/v1/
git commit -m "feat: add API handlers for AI settings and journal generation"
```

### Task 6: Web UI Integration

**Files:**
- Modify: `web/templates/settings.html`
- Modify: `web/templates/journal.html`

- [ ] **Step 1: Add AI Settings UI to `web/templates/settings.html`**
Add an AI Configuration card to the settings layout:
```html
<!-- Inside the settings container -->
<div class="card bg-base-100 shadow-md">
    <div class="card-body">
        <h2 class="card-title text-xl mb-4">✨ AI Configuration</h2>
        <form id="aiSettingsForm" class="space-y-4">
            <div class="form-control">
                <label class="label"><span class="label-text">Provider</span></label>
                <select id="aiProvider" class="select select-bordered w-full max-w-xs">
                    <option value="gemini">Google Gemini</option>
                    <option value="openai">OpenAI</option>
                </select>
            </div>
            <div class="form-control">
                <label class="label"><span class="label-text">Model</span></label>
                <input type="text" id="aiModel" placeholder="e.g., gemini-1.5-flash" class="input input-bordered w-full max-w-xs" />
            </div>
            <div class="form-control">
                <label class="label"><span class="label-text">API Token</span></label>
                <input type="password" id="aiToken" placeholder="Enter to update token" class="input input-bordered w-full max-w-xs" />
                <label class="label"><span class="label-text-alt" id="tokenStatus"></span></label>
            </div>
            <button type="submit" class="btn btn-primary mt-4">Save AI Settings</button>
        </form>
    </div>
</div>

<script>
document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/v1/ai/settings')
        .then(res => res.json())
        .then(data => {
            if(data.provider) document.getElementById('aiProvider').value = data.provider;
            if(data.model) document.getElementById('aiModel').value = data.model;
            document.getElementById('tokenStatus').innerText = data.hasToken ? "Token is configured securely." : "Token is not configured.";
        });

    document.getElementById('aiSettingsForm').addEventListener('submit', (e) => {
        e.preventDefault();
        const payload = {
            provider: document.getElementById('aiProvider').value,
            model: document.getElementById('aiModel').value,
            token: document.getElementById('aiToken').value,
        };
        fetch('/api/v1/ai/settings', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(payload)
        }).then(() => {
            alert('AI Settings saved!');
            document.getElementById('aiToken').value = ''; // clear field
        });
    });
});
</script>
```

- [ ] **Step 2: Add Generate Button to `web/templates/journal.html`**
Add a new button near the "New Entry" button:
```html
<!-- Near top actions -->
<div class="dropdown dropdown-end">
  <div tabindex="0" role="button" class="btn btn-secondary m-1">✨ Auto-generate Entry</div>
  <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
    <li><a onclick="generateAIJournal('daily')">Daily Summary</a></li>
    <li><a onclick="generateAIJournal('weekly')">Weekly Summary</a></li>
    <li><a onclick="generateAIJournal('yearly')">Yearly Summary</a></li>
  </ul>
</div>

<!-- Add JS -->
<script>
function generateAIJournal(period) {
    const btn = document.activeElement;
    btn.innerHTML = '<span class="loading loading-spinner"></span> Generating...';
    
    fetch(`/api/v1/ai/generate?period=${period}`)
        .then(res => {
            if(!res.ok) throw new Error("Generation failed (Check API token in settings)");
            return res.json();
        })
        .then(data => {
            // Open the new journal entry modal/editor and pre-fill
            document.getElementById('addJournalModal').showModal();
            document.getElementById('journalTitle').value = data.title;
            // Assuming EasyMDE is used for journal content:
            if(window.easyMDE) {
                window.easyMDE.value(data.content);
            } else {
                document.getElementById('journalContent').value = data.content;
            }
        })
        .catch(err => alert(err.message))
        .finally(() => {
            btn.innerHTML = '✨ Auto-generate Entry';
        });
}
</script>
```

- [ ] **Step 3: Commit**
```bash
git add web/templates/settings.html web/templates/journal.html
git commit -m "feat: add UI for AI configuration and journal generation"
```
