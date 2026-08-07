package queue

type ImportMessage struct {
	UID       string                 `json:"uid"`
	Bookmarks []*ImportBookmarkInput `json:"data"`
}

type ImportBookmarkInput struct {
	Description string `json:"description" validate:"lte=255" csv:"description"`
	URL         string `json:"url" validate:"required,url" csv:"url"`
}
