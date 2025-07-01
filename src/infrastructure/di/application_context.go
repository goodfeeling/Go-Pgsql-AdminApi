package di

import (
	"sync"

	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	medicineUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	filesUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/files"
	roleUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/role"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	medicineController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	roleController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/role"
	uploadController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/upload"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"gorm.io/gorm"
)

// ApplicationContext holds all application dependencies and services
type ApplicationContext struct {
	DB     *gorm.DB
	Logger *logger.Logger
	// controller
	AuthController     authController.IAuthController
	UserController     userController.IUserController
	MedicineController medicineController.IMedicineController
	UploadController   uploadController.IUploadController
	RoleController     roleController.IRoleController

	JWTService security.IJWTService
	// repository
	UserRepository     user.UserRepositoryInterface
	MedicineRepository medicine.MedicineRepositoryInterface
	FilesRepository    files.ISysFilesRepository
	RoleRepository     role.ISysRolesRepository
	// application
	AuthUseCase     authUseCase.IAuthUseCase
	UserUseCase     userUseCase.IUserUseCase
	MedicineUseCase medicineUseCase.IMedicineUseCase
	FilesUseCase    filesUseCase.ISysFilesService
	RoleUseCase     roleUseCase.ISysRoleService
}

var (
	loggerInstance *logger.Logger
	loggerOnce     sync.Once
)

func GetLogger() *logger.Logger {
	loggerOnce.Do(func() {
		loggerInstance, _ = logger.NewLogger()
	})
	return loggerInstance
}

// SetupDependencies creates a new application context with all dependencies
func SetupDependencies(loggerInstance *logger.Logger) (*ApplicationContext, error) {
	// Initialize database with logger
	db, err := psql.InitPSQLDB(loggerInstance)
	if err != nil {
		return nil, err
	}

	// Initialize JWT service (manages its own configuration)
	jwtService := security.NewJWTService()

	// Initialize repositories with logger
	userRepo := user.NewUserRepository(db, loggerInstance)
	medicineRepo := medicine.NewMedicineRepository(db, loggerInstance)
	jwtBlackListRepo := jwt_blacklist.NewUJwtBlacklistRepository(db)
	filesRepo := files.NewSysFilesRepository(db, loggerInstance)
	roleRepo := role.NewSysRolesRepository(db, loggerInstance)

	// Initialize use cases with logger
	authUC := authUseCase.NewAuthUseCase(userRepo, jwtService, loggerInstance, jwtBlackListRepo)
	userUC := userUseCase.NewUserUseCase(userRepo, loggerInstance)
	medicineUC := medicineUseCase.NewMedicineUseCase(medicineRepo, loggerInstance)
	filesUC := filesUseCase.NewSysFilesUseCase(filesRepo, loggerInstance)
	roleUC := roleUseCase.NewSysFilesUseCase(roleRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	medicineController := medicineController.NewMedicineController(medicineUC, loggerInstance)
	uploadController := uploadController.NewAuthController(filesUC, loggerInstance)
	roleController := roleController.NewRoleController(roleUC, loggerInstance)

	return &ApplicationContext{
		DB:     db,
		Logger: loggerInstance,
		// controller
		AuthController:     authController,
		UserController:     userController,
		MedicineController: medicineController,
		UploadController:   uploadController,
		RoleController:     roleController,
		JWTService:         jwtService,
		// repository
		UserRepository:     userRepo,
		MedicineRepository: medicineRepo,
		FilesRepository:    filesRepo,
		RoleRepository:     roleRepo,
		// application
		AuthUseCase:     authUC,
		UserUseCase:     userUC,
		MedicineUseCase: medicineUC,
		FilesUseCase:    filesUC,
		RoleUseCase:     roleUC,
	}, nil
}

// NewTestApplicationContext creates an application context for testing with mocked dependencies
func NewTestApplicationContext(
	mockUserRepo user.UserRepositoryInterface,
	mockMedicineRepo medicine.MedicineRepositoryInterface,
	mockJWTService security.IJWTService,
	loggerInstance *logger.Logger,
	jwtBlackListRepo jwt_blacklist.JwtBlacklistRepository,
) *ApplicationContext {
	// Initialize use cases with mocked repositories and logger
	authUC := authUseCase.NewAuthUseCase(mockUserRepo, mockJWTService, loggerInstance, jwtBlackListRepo)
	userUC := userUseCase.NewUserUseCase(mockUserRepo, loggerInstance)
	medicineUC := medicineUseCase.NewMedicineUseCase(mockMedicineRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	medicineController := medicineController.NewMedicineController(medicineUC, loggerInstance)

	return &ApplicationContext{
		Logger:             loggerInstance,
		AuthController:     authController,
		UserController:     userController,
		MedicineController: medicineController,
		JWTService:         mockJWTService,
		UserRepository:     mockUserRepo,
		MedicineRepository: mockMedicineRepo,
		AuthUseCase:        authUC,
		UserUseCase:        userUC,
		MedicineUseCase:    medicineUC,
	}
}
