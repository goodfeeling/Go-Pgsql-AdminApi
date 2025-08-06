package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	domain "github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainFiles "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	shareUtils "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type IUploadController interface {
	Single(ctx *gin.Context)
	Multiple(ctx *gin.Context)
	GetSTSToken(ctx *gin.Context)
}

type STSTokenResponse struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SecurityToken   string `json:"security_token"`
	Expiration      string `json:"expiration"`
	BucketName      string `json:"bucket_name"`
	Region          string `json:"region"`
}

type UploadController struct {
	sysFilesUseCase domainFiles.ISysFilesService
	Logger          *logger.Logger
}

// MultipleUpload
// @Summary multiple files upload
// @Description upload multiple files get files info
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "fileResources" collectionFormat(multi)
// @Success 200 {object} domain.CommonResponse[[]domainFiles.SysFiles]
// @Router /v1/upload/multiple [post]
func (u *UploadController) Multiple(ctx *gin.Context) {
	// 获取多文件表单
	form, err := ctx.MultipartForm()
	if err != nil {
		u.Logger.Error("Failed to get multipart form", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	files := form.File["file"]
	var uploadedFiles []domainFiles.SysFiles

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		ext := filepath.Ext(filename)
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("public", newFilename)

		if err := os.MkdirAll("public", os.ModePerm); err != nil {
			u.Logger.Error("Error creating dir", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			u.Logger.Error("Error save file", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		md5Value, err := shareUtils.CalculateFileMD5(savePath)
		if err != nil {
			u.Logger.Error("calculate file to md5", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		fileInfo := domainFiles.SysFiles{
			FileName:       newFilename,
			FilePath:       savePath,
			FileMD5:        md5Value,
			FileOriginName: filename,
			StorageEngine:  "local",
		}

		res, err := u.sysFilesUseCase.Create(&fileInfo)
		if err != nil {
			u.Logger.Error("insert file info to database", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		uploadedFiles = append(uploadedFiles, *res)
	}

	response := &domain.CommonResponse[[]domainFiles.SysFiles]{
		Data:    uploadedFiles,
		Message: "Upload success",
		Status:  200,
	}

	u.Logger.Info("multiple upload successful", zap.Int("fileCount", len(files)))

	ctx.JSON(http.StatusOK, response)
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
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
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
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		u.Logger.Error("Error save file", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	// calculate md5 file
	md5Value, err := shareUtils.CalculateFileMD5(savePath)
	if err != nil {
		u.Logger.Error("calculate  file to md5", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	// insert to database
	files := domainFiles.SysFiles{
		FileName:       newFilename,
		FilePath:       savePath,
		FileMD5:        md5Value,
		FileOriginName: filename,
		StorageEngine:  "local",
	}
	res, err := u.sysFilesUseCase.Create(&files)
	if err != nil {
		u.Logger.Error("insert file info to database", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
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

// GetSTSToken
// @Summary get sts token with aliyun
// @Description get sts token
// @Tags sts token
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse
// @Router /v1/upload/sts-token [get]
func (u *UploadController) GetSTSToken(ctx *gin.Context) {
	// 从环境变量中获取步骤1.1生成的RAM用户的访问密钥（AccessKey ID和AccessKey Secret）。
	accessKeyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	// 从环境变量中获取步骤1.3生成的RAM角色的RamRoleArn。
	roleArn := os.Getenv("RAM_ROLE_ARN")
	serviceAddress := os.Getenv("SECURITY_SERVICE_ADDRESS")
	// 创建权限策略客户端。
	config := &openapi.Config{
		// 必填，步骤1.1获取到的 AccessKey ID。
		AccessKeyId: tea.String(accessKeyId),
		// 必填，步骤1.1获取到的 AccessKey Secret。
		AccessKeySecret: tea.String(accessKeySecret),
	}

	// Endpoint 请参考 https://api.aliyun.com/product/Sts
	config.Endpoint = tea.String(serviceAddress)
	client, err := sts20150401.NewClient(config)
	if err != nil {
		u.Logger.Error("Failed to create client:", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}
	// 生成唯一的会话名称
	sessionName := fmt.Sprintf("upload-session-%d", time.Now().Unix())

	// 使用RAM用户的AccessKey ID和AccessKey Secret向STS申请临时访问凭证。
	request := &sts20150401.AssumeRoleRequest{
		// 指定STS临时访问凭证过期时间为3600秒。
		DurationSeconds: tea.Int64(3600),
		// 从环境变量中获取步骤1.3生成的RAM角色的RamRoleArn。
		RoleArn: tea.String(roleArn),
		// 指定自定义角色会话名称，这里使用和第一段代码一致的 examplename
		RoleSessionName: tea.String(sessionName),
	}
	response, err := client.AssumeRoleWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		u.Logger.Error("Failed to assume role:", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	// 打印STS返回的临时访问密钥（AccessKey ID和AccessKey Secret）、安全令牌（SecurityToken）以及临时访问凭证过期时间（Expiration）。
	credentials := response.Body.Credentials
	result := controllers.NewCommonResponseBuilder[*STSTokenResponse]().
		Data(&STSTokenResponse{
			AccessKeyId:     *credentials.AccessKeyId,
			AccessKeySecret: *credentials.AccessKeySecret,
			SecurityToken:   *credentials.SecurityToken,
			Expiration:      *credentials.Expiration,
			BucketName:      os.Getenv("OSS_BUCKET_NAME"),
			Region:          os.Getenv("SECURITY_REGION_ID"),
		}).
		Message("success").
		Status(0).
		Build()
	ctx.JSON(http.StatusOK, result)
}

func NewAuthController(sysFilesUseCase domainFiles.ISysFilesService, loggerInstance *logger.Logger) IUploadController {
	return &UploadController{
		sysFilesUseCase: sysFilesUseCase,
		Logger:          loggerInstance,
	}
}
