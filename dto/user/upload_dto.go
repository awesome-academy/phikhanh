package user

// Response upload file
type UploadFileResponse struct {
	FileName     string `json:"file_name" example:"20240115143022_a1b2c3d4.jpg"`
	FileURL      string `json:"file_url" example:"/assets/images/20240115143022_a1b2c3d4.jpg"`
	FileSize     int64  `json:"file_size" example:"1024000"`
	OriginalName string `json:"original_name" example:"document.jpg"`
}
