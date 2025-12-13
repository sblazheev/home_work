package common

type Notification struct {
	ID          string `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
	Time        uint64 `json:"time" validate:"required"`
	NotifyTime  uint64 `json:"notifyTime" validate:"required"`
}
