package email

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gbrayhan/microservices-go/src/application/event/bus"
	"github.com/gbrayhan/microservices-go/src/application/event/model"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type IEmailService interface {
	SendForgetPasswordEmail(userName string) error
}

type EmailUseCase struct {
	Logger         *logger.Logger
	eventBus       bus.EventBus
	RedisClient    *redis.Client
	jwtService     security.IJWTService
	userRepository user.UserRepositoryInterface
}

func NewEmailUseCase(
	userRepository user.UserRepositoryInterface,
	jwtService security.IJWTService,
	eventBus bus.EventBus,
	RedisClient *redis.Client,
	loggerInstance *logger.Logger) IEmailService {
	return &EmailUseCase{
		Logger:         loggerInstance,
		eventBus:       eventBus,
		RedisClient:    RedisClient,
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}

// SendEmail implements IEmailService.
func (e *EmailUseCase) SendForgetPasswordEmail(email string) error {

	user, err := e.userRepository.GetByEmail(email)
	if err != nil || user == nil {
		return errors.New("place enter a valid email")
	}

	return nil
	resetLink, err := e.generateResetLink(email)
	if err != nil {
		return err
	}

	event := &model.ForgetPasswordEvent{
		ID:           uuid.New().String(),
		To:           email,
		Subject:      "Reset Password",
		Body:         "Please click the link to reset your password: " + resetLink,
		RegisteredAt: time.Now(),
	}

	return e.eventBus.Publish(context.Background(), event)
}

func (e *EmailUseCase) generateResetLink(email string) (string, error) {
	token, err := e.generateSecureToken()
	if err != nil {
		return "", err
	}
	//  将token存储到数据库或缓存中
	e.RedisClient.Set(context.Background(), GetEmailTokenKey(email), token, UserTokenExpireDuration)

	// 返回包含token的链接
	return fmt.Sprintf("%s/#/auth/reset-password?token=%s", os.Getenv("SERVER_FRONTEND_URL"), token), nil
}

func (e *EmailUseCase) generateSecureToken() (string, error) {
	token, err := e.jwtService.GenerateJWTToken(0, 0, "reset")
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
