package user

import (
	"net/http"

	userDto "phikhanh/dto/user"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadFile godoc
// @Summary      Upload file
// @Description  Upload file (images, documents) lên server
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file formData file true "File to upload"
// @Success      200  {object}  utils.APIResponse{data=userDto.UploadFileResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /upload [post]
func (c *UploadController) UploadFile(ctx *gin.Context) {
	// Get file from form
	file, err := ctx.FormFile("file")
	if err != nil {
		svcErr := utils.NewBadRequestError("No file uploaded")
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	// Upload file
	result, err := utils.UploadFile(file)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		svcErr := utils.NewInternalServerError(err)
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	response := userDto.UploadFileResponse{
		FileName:     result.FileName,
		FileURL:      result.FileURL,
		FileSize:     result.FileSize,
		OriginalName: result.OriginalName,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "File uploaded successfully", response)
}
