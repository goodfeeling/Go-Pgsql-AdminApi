package di

import (
	"sync"

	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	medicineUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/medicine"
	apiUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/api"
	dictionaryUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/dictionary"
	dictionaryDetailUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/dictionary_detail"
	filesUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/files"
	menuUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/menu"
	menuBtnUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/menu_btn"
	menuGroupUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/menu_group"
	menuParameterUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/menu_parameter"
	operationUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/operation_record"
	roleUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/role"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_parameter"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/casbin_rule"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary_detail"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/operation_records"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/user_role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	apiController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/api"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	dictionaryController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionary"
	dictionaryDetailController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionaryDetail"
	medicineController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	menuController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu"
	menuBtnController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menuBtn"
	menuGroupController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menuGroup"
	menuParameterController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menuParameter"
	operationController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/operation"
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
	AuthController             authController.IAuthController
	UserController             userController.IUserController
	MedicineController         medicineController.IMedicineController
	UploadController           uploadController.IUploadController
	RoleController             roleController.IRoleController
	ApiController              apiController.IApiController
	JWTService                 security.IJWTService
	MenuController             menuController.IMenuController
	OperationController        operationController.IOperationController
	DictionaryController       dictionaryController.IDictionaryController
	DictionaryDetailController dictionaryDetailController.IIDictionaryDetailController
	MenuGroupController        menuGroupController.IMenuGroupController
	MenuBtnController          menuBtnController.IMenuBtnController
	MenuParameterController    menuParameterController.IMenuParameterController
	// repository
	UserRepository             user.UserRepositoryInterface
	MedicineRepository         medicine.MedicineRepositoryInterface
	FilesRepository            files.ISysFilesRepository
	RoleRepository             role.ISysRolesRepository
	ApiRepository              api.ApiRepositoryInterface
	MenuRepository             base_menu.MenuRepositoryInterface
	OperationRepository        operation_records.OperationRepositoryInterface
	DictionaryRepository       dictionary.DictionaryRepositoryInterface
	DictionaryDetailRepository dictionary_detail.DictionaryRepositoryInterface
	MenuGroupRepository        base_menu_group.MenuGroupRepositoryInterface
	MenuBtnRepository          base_menu_btn.MenuBtnRepositoryInterface
	MenuParameterRepository    base_menu_parameter.MenuParameterRepositoryInterface
	// application
	AuthUseCase             authUseCase.IAuthUseCase
	UserUseCase             userUseCase.IUserUseCase
	MedicineUseCase         medicineUseCase.IMedicineUseCase
	FilesUseCase            filesUseCase.ISysFilesService
	RoleUseCase             roleUseCase.ISysRoleService
	ApiUseCase              apiUseCase.ISysApiService
	MenuUseCase             menuUseCase.ISysMenuService
	OperationUseCase        operationUseCase.ISysOperationService
	DictionaryUseCase       dictionaryUseCase.ISysDictionaryService
	DictionaryDetailUseCase dictionaryDetailUseCase.ISysDictionaryService
	MenuGroupUseCase        menuGroupUseCase.ISysMenuGroupService
	MenuBtnUseCase          menuBtnUseCase.IMenuBtnService
	menuParameterUseCase    menuParameterUseCase.IMenuParameterService
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
	userRoleRepo := user_role.NewSysUserRoleRepository(db, loggerInstance)
	medicineRepo := medicine.NewMedicineRepository(db, loggerInstance)
	jwtBlackListRepo := jwt_blacklist.NewUJwtBlacklistRepository(db)
	filesRepo := files.NewSysFilesRepository(db, loggerInstance)
	roleRepo := role.NewSysRolesRepository(db, loggerInstance)
	apiRepo := api.NewApiRepository(db, loggerInstance)
	operationRepo := operation_records.NewOperationRepository(db, loggerInstance)
	dictionaryRepo := dictionary.NewDictionaryRepository(db, loggerInstance)
	dictionaryDetailRepo := dictionary_detail.NewDictionaryRepository(db, loggerInstance)
	menuRepo := base_menu.NewMenuRepository(db, loggerInstance)
	roleMenuRepo := role_menu.NewSysRoleMenuRepository(db, loggerInstance)
	casBinRepo := casbin_rule.NewCasbinRuleRepository(db, loggerInstance)
	menuGroupRepo := base_menu_group.NewMenuGroupRepository(db, loggerInstance)
	menuBtnRepo := base_menu_btn.NewMenuBtnRepository(db, loggerInstance)
	menuParameterRepo := base_menu_parameter.NewMenuParameterRepository(db, loggerInstance)
	roleBtnRepo := role_btn.NewRoleBtnRepository(db, loggerInstance)
	// Initialize use cases with logger
	authUC := authUseCase.NewAuthUseCase(userRepo, roleRepo, jwtService, loggerInstance, jwtBlackListRepo)
	userUC := userUseCase.NewUserUseCase(userRepo, userRoleRepo, loggerInstance)
	medicineUC := medicineUseCase.NewMedicineUseCase(medicineRepo, loggerInstance)
	filesUC := filesUseCase.NewSysFilesUseCase(filesRepo, loggerInstance)
	roleUC := roleUseCase.NewSysRoleUseCase(roleRepo, roleMenuRepo, casBinRepo, menuRepo, roleBtnRepo, loggerInstance)
	apiUC := apiUseCase.NewSysApiUseCase(apiRepo, loggerInstance)
	operationUC := operationUseCase.NewSysOperationUseCase(operationRepo, loggerInstance)
	dictionaryUC := dictionaryUseCase.NewSysDictionaryUseCase(dictionaryRepo, loggerInstance)
	dictionaryDetailUC := dictionaryDetailUseCase.NewSysDictionaryUseCase(dictionaryDetailRepo, loggerInstance)
	menuUC := menuUseCase.NewSysMenuUseCase(menuRepo, roleMenuRepo, userRepo, menuGroupRepo, loggerInstance)
	menuGroupUC := menuGroupUseCase.NewSysMenuGroupUseCase(menuGroupRepo, loggerInstance)
	menuBtnUC := menuBtnUseCase.NewMenuBtnUseCase(menuBtnRepo, loggerInstance)
	menuParameterUC := menuParameterUseCase.NewMenuParameterUseCase(menuParameterRepo, loggerInstance)

	// Initialize controllers with logger
	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	medicineController := medicineController.NewMedicineController(medicineUC, loggerInstance)
	uploadController := uploadController.NewAuthController(filesUC, loggerInstance)
	roleController := roleController.NewRoleController(roleUC, loggerInstance)
	apiController := apiController.NewApiController(apiUC, loggerInstance)
	operationController := operationController.NewOperationController(operationUC, loggerInstance)
	dictionaryController := dictionaryController.NewDictionaryController(dictionaryUC, loggerInstance)
	dictionaryDetailController := dictionaryDetailController.NewIDictionaryDetailController(dictionaryDetailUC, loggerInstance)
	menuController := menuController.NewMenuController(menuUC, loggerInstance)
	menuGroupController := menuGroupController.NewMenuGroupController(menuGroupUC, loggerInstance)
	menuBtnController := menuBtnController.NewMenuBtnController(menuBtnUC, loggerInstance)
	menuParameterController := menuParameterController.NewMenuParameterController(menuParameterUC, loggerInstance)

	return &ApplicationContext{
		DB:     db,
		Logger: loggerInstance,
		// controller
		AuthController:             authController,
		UserController:             userController,
		MedicineController:         medicineController,
		UploadController:           uploadController,
		RoleController:             roleController,
		ApiController:              apiController,
		OperationController:        operationController,
		DictionaryController:       dictionaryController,
		DictionaryDetailController: dictionaryDetailController,
		MenuController:             menuController,
		MenuGroupController:        menuGroupController,
		MenuBtnController:          menuBtnController,
		MenuParameterController:    menuParameterController,
		// repository
		UserRepository:             userRepo,
		MedicineRepository:         medicineRepo,
		FilesRepository:            filesRepo,
		RoleRepository:             roleRepo,
		ApiRepository:              apiRepo,
		OperationRepository:        operationRepo,
		DictionaryRepository:       dictionaryRepo,
		DictionaryDetailRepository: dictionaryDetailRepo,
		MenuRepository:             menuRepo,
		MenuGroupRepository:        menuGroupRepo,
		MenuBtnRepository:          menuBtnRepo,
		MenuParameterRepository:    menuParameterRepo,
		// application
		AuthUseCase:             authUC,
		UserUseCase:             userUC,
		MedicineUseCase:         medicineUC,
		FilesUseCase:            filesUC,
		RoleUseCase:             roleUC,
		ApiUseCase:              apiUC,
		OperationUseCase:        operationUC,
		DictionaryUseCase:       dictionaryUC,
		DictionaryDetailUseCase: dictionaryDetailUC,
		MenuUseCase:             menuUC,
		MenuGroupUseCase:        menuGroupUC,
		MenuBtnUseCase:          menuBtnUC,
		menuParameterUseCase:    menuParameterUC,

		JWTService: jwtService,
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
	// authUC := authUseCase.NewAuthUseCase(mockUserRepo, mockJWTService, loggerInstance, jwtBlackListRepo)
	// userUC := userUseCase.NewUserUseCase(mockUserRepo, loggerInstance)
	medicineUC := medicineUseCase.NewMedicineUseCase(mockMedicineRepo, loggerInstance)

	// Initialize controllers with logger
	// authController := authController.NewAuthController(authUC, loggerInstance)
	// userController := userController.NewUserController(userUC, loggerInstance)
	medicineController := medicineController.NewMedicineController(medicineUC, loggerInstance)

	return &ApplicationContext{
		Logger: loggerInstance,
		// AuthController: authController,
		// UserController:     userController,
		MedicineController: medicineController,
		JWTService:         mockJWTService,
		UserRepository:     mockUserRepo,
		MedicineRepository: mockMedicineRepo,
		// AuthUseCase:        authUC,
		// UserUseCase:        userUC,
		MedicineUseCase: medicineUC,
	}
}
