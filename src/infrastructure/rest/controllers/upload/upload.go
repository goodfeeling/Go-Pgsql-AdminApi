package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	domainFiles "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
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

// single implements IUploadController.
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
}

func NewAuthController(sysFilesUseCase domainFiles.ISysFilesService, loggerInstance *logger.Logger) IUploadController {
	return &UploadController{
		sysFilesUseCase: sysFilesUseCase,
		Logger:          loggerInstance,
	}
}
