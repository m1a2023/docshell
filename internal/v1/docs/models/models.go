package models

type Document struct {
	Id         int    `json:"id" db:"id"`
	AuthorId   int    `json:"author_id" db:"author_id"`
	UploaderId int    `json:"uploader_id" db:"uploader_id"`
	Title      string `json:"title" db:"title"`
	Size       int    `json:"size" db:"size"`
	Path       string `json:"path" db:"path"`
	Hash       string `json:"hash" db:"hash"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	ChangedAt  string `json:"changed_at" db:"changed_at"`
}

type ResponseMultipleDocuments struct {
	StatusCode int        `json:"status_code"`
	Documents  []Document `json:"documents"`
}

type ResponseCode struct {
	StatusCode int `json:"status_code"`
}
