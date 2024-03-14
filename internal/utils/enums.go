package utils

type EventType string

const (
	LogEvent      EventType = "log_event"
	GalleryName   EventType = "gallery_name"
	Progress      EventType = "progress"
	ImageProgress EventType = "image_progress"
	InfoTask      EventType = "info_task"
	Status        EventType = "status"
)
