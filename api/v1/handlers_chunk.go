package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/snehmatic/mindloop/internal/core/ai"
)

func (mlh *MindloopHandler) HandleChunk(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	itemType := r.FormValue("type")
	idStr := r.FormValue("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var itemName string
	switch itemType {
	case "intent":
		i, err := mlh.intent.GetIntent(idStr)
		if err != nil || i == nil {
			http.Error(w, "Intent not found", http.StatusNotFound)
			return
		}
		itemName = i.Name
	case "task":
		t, err := mlh.task.GetTask(uint(id))
		if err != nil || t == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		itemName = t.Title
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	aiSvc := ai.NewService(mlh.journal.DB) // Using DB from existing service

	res, err := aiSvc.GenerateChunker(itemName)
	if err != nil {
		http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var steps []string
	cleanRes := strings.TrimSpace(res)
	if strings.HasPrefix(cleanRes, "```json") {
		cleanRes = strings.TrimPrefix(cleanRes, "```json")
		cleanRes = strings.TrimSuffix(cleanRes, "```")
	} else if strings.HasPrefix(cleanRes, "```") {
		cleanRes = strings.TrimPrefix(cleanRes, "```")
		cleanRes = strings.TrimSuffix(cleanRes, "```")
	}

	err = json.Unmarshal([]byte(cleanRes), &steps)
	if err != nil {
		http.Error(w, "Failed to parse AI response", http.StatusInternalServerError)
		return
	}

	// Create the tasks/subtasks and return HTML snippets for HTMX
	w.Header().Set("Content-Type", "text/html")
	for _, step := range steps {
		switch itemType {
		case "intent":
			uid := uint(id)
			t, err := mlh.task.CreateTask(step, &uid, nil)
			if err == nil {
				// Render task HTML
				html := fmt.Sprintf(`<div class="card p-sm animate-slide-up" style="display: flex; justify-content: space-between; align-items: center; border-left: 2px solid var(--primary);">
					<div style="display: flex; align-items: center; gap: 0.5rem;">
						<i data-lucide="circle" style="color: var(--text-secondary); width: 18px;"></i>
						%s
					</div>
					<form action="/tasks/complete" method="POST" class="mb-0">
						<input type="hidden" name="id" value="%d">
						<input type="hidden" name="source" value="/intent">
						<button type="submit" class="btn btn-sm btn-success-outline" style="padding: 0.25rem 0.5rem;">Done</button>
					</form>
				</div>`, t.Title, t.ID)
				_, _ = w.Write([]byte(html))
			}
		case "task":
			st, err := mlh.task.AddSubTask(uint(id), step)
			if err == nil {
				html := fmt.Sprintf(`<div class="subtask-item animate-slide-up" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; padding: 0.5rem; background: var(--bg-secondary); border-radius: var(--radius-sm); border-left: 2px solid var(--primary);">
					<div style="display: flex; align-items: center; gap: 0.5rem;">
						<i data-lucide="circle" style="width: 14px; height: 14px; color: var(--text-muted);"></i>
						<span style="font-size: 0.9em;">%s</span>
					</div>
					<form action="/tasks/subtask/complete" method="POST" style="margin: 0;">
						<input type="hidden" name="id" value="%d">
						<button type="submit" class="btn btn-sm btn-secondary" style="padding: 0.2rem 0.5rem; font-size: 0.75rem;"><i data-lucide="check" style="width:12px;height:12px; margin-right: 2px;"></i> Done</button>
					</form>
				</div>`, st.Title, st.ID)
				_, _ = w.Write([]byte(html))
			}
		}
	}
}
