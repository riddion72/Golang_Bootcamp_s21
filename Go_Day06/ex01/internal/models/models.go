package models

type Article struct {
	Article_name    string `json:"Article_name"`
	Article_content string `json:"Article_content"`
}

type GetArticleResponse struct {
	Response []Article `json:"Article"`
	Last     int       `json:"Last"`
	Page     int       `json:"Page`
}

type GetArticleRequest struct {
	Page  int `json: "Page"`
	Ofset int `json: "Ofset", omitempty"`
	Limit int `json: "Limit", omitempty`
}

// type SetArticleResponse struct {
// 	Ofset string `json:"Ofset"`
// 	Limit string `json:"Limit"`
// }

type SetArticleRequest struct {
	Request Article
}
