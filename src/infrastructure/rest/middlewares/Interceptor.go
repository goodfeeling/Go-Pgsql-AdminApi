package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	operationRecordsDomain "github.com/gbrayhan/microservices-go/src/domain/sys/operation_records"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	operationRecordsRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/operation_records"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware(db *gorm.DB, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody string
		var resp string
		if c.Request.RequestURI == "/v1/upload/single" {
			reqBody = ""
			resp = ""
		} else {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			buf := make([]byte, 4096)
			num, err := c.Request.Body.Read(buf)
			if err != nil && err.Error() != "EOF" {
				_ = fmt.Errorf("error reading buffer: %s", err.Error())
			}
			reqBody := string(buf[0:num])
			resp = blw.body.String()
			c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(reqBody)))
		}

		loc, _ := time.LoadLocation("America/Mexico_City")

		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		operationRecordsRepository := operationRecordsRepository.NewOperationRepository(db, logger)

		// handle user Id
		userIdRaw, exist := c.Get("user_id")

		var userId int64
		if !exist {
			userId = 0 // no auth user
		} else {
			switch v := userIdRaw.(type) {
			case int:
				userId = int64(v)
			case int64:
				userId = v
			case float64:
				userId = int64(v)
			case string:
				fmt.Sscanf(v, "%d", &userId)
			default:
				userId = 0
			}
		}

		operationRecordsRepository.Create(&operationRecordsDomain.SysOperationRecord{
			IP:           c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.RequestURI,
			Status:       int64(c.Writer.Status()),
			Agent:        c.Request.UserAgent(),
			Body:         reqBody,
			Resp:         resp,
			ErrorMessage: c.Errors.String(),
			UserID:       userId,
			Latency:      latency,
			CreatedAt:    domain.CustomTime{Time: time.Now().In(loc)},
		})

		allDataIO := map[string]any{
			"ruta":          c.FullPath(),
			"request_uri":   c.Request.RequestURI,
			"raw_request":   reqBody,
			"status_code":   c.Writer.Status(),
			"body_response": resp,
			"errors":        c.Errors.Errors(),
			"created_at":    time.Now().In(loc).Format("2006-01-02T15:04:05"),
		}

		_ = fmt.Sprintf("%v", allDataIO)
	}
}
