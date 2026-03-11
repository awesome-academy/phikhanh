package user

type SupplementAttachment struct {
	FilePath string `json:"file_path" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

type SupplementRequest struct {
	Attachments []SupplementAttachment `json:"attachments" binding:"required,min=1,dive"`
	Note        string                 `json:"note"`
}
