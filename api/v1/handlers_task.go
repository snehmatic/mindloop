package v1

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/models"
)

// --- Task Handlers ---

func (mlh *MindloopHandler) HandleTaskList(w http.ResponseWriter, r *http.Request) {
	tasks, err := mlh.task.ListTasks()
	if err != nil {
		log.Error().Err(err).Msg("Error listing tasks")
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	intents, _ := mlh.intent.ListIntents()
	intentMap := make(map[uint]string)
	for _, i := range intents {
		intentMap[i.ID] = i.Name
	}
	sessions, _ := mlh.focus.ListSessions()
	sessionMap := make(map[uint]string)
	for _, s := range sessions {
		sessionMap[s.ID] = s.Title
	}

	var taskViews []models.TaskView
	for _, t := range tasks {
		tv := models.ToTaskView(t)
		if tv.IntentID != nil {
			tv.IntentName = intentMap[*tv.IntentID]
		}
		if tv.FocusSessionID != nil {
			tv.FocusSessionTitle = sessionMap[*tv.FocusSessionID]
		}
		taskViews = append(taskViews, tv)
	}
	// Filter active intents and sessions for the association dropdowns
	var activeIntents []models.Intent
	for _, i := range intents {
		if i.Status == "active" {
			activeIntents = append(activeIntents, i)
		}
	}
	var activeSessions []models.FocusSession
	for _, s := range sessions {
		if s.Status == "active" {
			activeSessions = append(activeSessions, s)
		}
	}

	data := map[string]interface{}{
		"Title":          "Tasks",
		"Tasks":          taskViews,
		"ActiveIntents":  activeIntents,
		"ActiveSessions": activeSessions,
	}

	if success := r.URL.Query().Get("success"); success != "" {
		switch success {
		case "added", "true":
			data["SuccessMessage"] = "Action completed."
		case "done":
			data["SuccessMessage"] = "Task completed successfully! Great job!"
		case "milestone":
			data["SuccessMessage"] = "🏆 MILESTONE REACHED! You are amazing! 🏆"
		}
	}
	if errStr := r.URL.Query().Get("error"); errStr != "" {
		data["ErrorMessage"] = errStr
	}

	mlh.renderTemplate(w, "tasks.html", data)
}

func (mlh *MindloopHandler) HandleTaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	intentIDStr := r.FormValue("intent_id")
	focusIDStr := r.FormValue("focus_id")
	source := r.FormValue("source") // allows redirecting back to /intent or /focus

	var intentID *uint
	if intentIDStr != "" {
		id, err := strconv.ParseUint(intentIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			intentID = &uid
		}
	}

	var focusID *uint
	if focusIDStr != "" {
		id, err := strconv.ParseUint(focusIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			focusID = &uid
		}
	}

	_, err := mlh.task.CreateTask(title, intentID, focusID)
	if err != nil {
		log.Error().Err(err).Msg("Error creating task")
		redirectUrl := "/tasks?error=Failed to create task"
		if source != "" {
			redirectUrl = source + "?error=Failed to create task"
		}
		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
		return
	}

	redirectUrl := "/tasks"
	if source != "" {
		redirectUrl = source
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleTaskComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	source := r.FormValue("source")

	uc := config.UserConfig{}
	_ = uc.ReadFromYAML()
	pointsVal := uc.PointsConfig.Task

	milestoneReached, err := mlh.task.CompleteTask(uint(id), pointsVal)
	if err != nil {
		log.Error().Err(err).Msg("Error completing task")
	}

	successCode := "done"
	if milestoneReached {
		successCode = "milestone"
	}

	redirectUrl := "/tasks?success=" + successCode
	if source != "" {
		redirectUrl = source + "?success=" + successCode
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleTaskDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	source := r.FormValue("source")

	err := mlh.task.DeleteTask(uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Error deleting task")
	}

	redirectUrl := "/tasks"
	if source != "" {
		redirectUrl = source
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleSubtaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	taskIDStr := r.FormValue("task_id")
	title := r.FormValue("title")
	source := r.FormValue("source")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err == nil {
		_, func_err := mlh.task.AddSubTask(uint(taskID), title)
		if func_err != nil {
			log.Error().Err(err).Msg("Error creating subtask")
		}
	}

	redirectUrl := "/tasks"
	if source != "" {
		redirectUrl = source
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleSubtaskComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	source := r.FormValue("source")

	uc := config.UserConfig{}
	_ = uc.ReadFromYAML()
	pointsVal := uc.PointsConfig.SubTask

	milestoneReached, err := mlh.task.CompleteSubTask(uint(id), pointsVal)
	if err != nil {
		log.Error().Err(err).Msg("Error completing subtask")
	}

	successCode := "done"
	if milestoneReached {
		successCode = "milestone"
	}

	redirectUrl := "/tasks?success=" + successCode
	if source != "" {
		redirectUrl = source + "?success=" + successCode
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleSubtaskDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	source := r.FormValue("source")

	err := mlh.task.DeleteSubTask(uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Error deleting subtask")
	}

	redirectUrl := "/tasks"
	if source != "" {
		redirectUrl = source
	}
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleTaskReorder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	idStrs := r.Form["task"]
	var ids []uint
	for _, idStr := range idStrs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) > 0 {
		if err := mlh.task.ReorderTasks(ids); err != nil {
			log.Error().Err(err).Msg("Error reordering tasks")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (mlh *MindloopHandler) HandleSubTaskReorder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	idStrs := r.Form["subtask"]
	var ids []uint
	for _, idStr := range idStrs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) > 0 {
		if err := mlh.task.ReorderSubTasks(ids); err != nil {
			log.Error().Err(err).Msg("Error reordering subtasks")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
