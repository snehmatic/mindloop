package models

type TaskStatus string

const (
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusPending   TaskStatus = "pending"
)

type IntentStatus string

const (
	IntentStatusActive IntentStatus = "active"
	IntentStatusPaused IntentStatus = "paused"
	IntentStatusDone   IntentStatus = "done"
)

type FocusSessionStatus string

const (
	FocusSessionStatusActive FocusSessionStatus = "active"
)

type SideQuestStatus string

const (
	SideQuestStatusActive SideQuestStatus = "active"
	SideQuestStatusDone   SideQuestStatus = "done"
)
