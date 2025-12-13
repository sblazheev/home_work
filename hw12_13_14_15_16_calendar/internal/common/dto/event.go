package dto

type Event struct {
	ID          string `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Title       string `json:"title" example:"Title" validate:"required"`
	Description string `json:"description" example:"Description"`
	DateTime    uint64 `json:"dateTime" example:"1740000000" validate:"required,numeric,gte=1740000000,lte=2762340413"`
	Duration    uint64 `json:"duration" example:"600" validate:"required"`
	UserID      string `json:"userId" example:"1" validate:"required,min=1"`
	NotifyTime  uint64 `json:"notifyTime" example:"600"`
} // @name Event .
