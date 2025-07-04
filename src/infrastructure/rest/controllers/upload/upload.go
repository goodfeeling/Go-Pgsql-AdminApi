package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	domain "github.com/gbrayhan/microservices-go/src/domain"
	domainFiles "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	shareUtils "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IUploadController interface {
	Single(ctx *gin.Context)
	Multiple(ctx *gin.Context)
}

type UploadController struct {
	sysFilesUseCase domainFiles.ISysFilesService
	Logger          *logger.Logger
}

// multiple implements IUploadController.
func (u *UploadController) Multiple(ctx *gin.Context) {
	panic("unimplemented")
}

// SingleUpload
// @Summary single file upload
// @Description upload single file get file info
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "fileResource"
// @Success 200 {object} domain.CommonResponse[domainFiles.SysFiles]
// @Router /v1/upload/single [post]
func (u *UploadController) Single(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		u.Logger.Error("Failed to get file", zap.Error(err))
		fmt.Printf("Error: %v\n", err)
		_ = ctx.Error(err)
		return
	}

	filename := filepath.Base(file.Filename)
	// only name
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// join name
	savePath := filepath.Join("public", newFilename)

	// create save file dir
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		u.Logger.Error("Error creating dir", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		u.Logger.Error("Error save file", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	// calculate md5 file
	md5Value, err := shareUtils.CalculateFileMD5(savePath)
	if err != nil {
		u.Logger.Error("calculate  file to md5", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	// insert to database
	files := domainFiles.SysFiles{
		FileName: newFilename,
		FilePath: savePath,
		FileMD5:  md5Value,
	}
	res, err := u.sysFilesUseCase.Create(&files)
	if err != nil {
		u.Logger.Error("insert file info to database", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	response := &domain.CommonResponse[domainFiles.SysFiles]{
		Data:    *res,
		Message: "Upload success",
		Status:  200,
	}

	u.Logger.Info("upload successful", zap.String("filename", newFilename))

	ctx.JSON(http.StatusOK, response)
}

func NewAuthController(sysFilesUseCase domainFiles.ISysFilesService, loggerInstance *logger.Logger) IUploadController {
	return &UploadController{
		sysFilesUseCase: sysFilesUseCase,
		Logger:          loggerInstance,
	}
}
