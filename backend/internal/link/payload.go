package link

type CreateLinkRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type UpdateLinkRequest struct {
	URL  string `json:"url" validate:"required,url"`
	Hash string `json:"hash" validate:"required"`
}

type DeleteLinkRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type GetLinksResponse struct {
	Links []Link `json:"links"`
	Count int64  `json:"count"`
}
