// #file:/root/myproject/microapp/src/infrastructure/di/application_context.go
package di

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/application/event/bus"
	"github.com/gbrayhan/microservices-go/src/application/event/factory"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	redisClient "github.com/gbrayhan/microservices-go/src/infrastructure/redis"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_parameter"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/casbin_rule"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary_detail"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/user_role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	sharedUtil "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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

// ApplicationContext holds all application dependencies and services
type ApplicationContext struct {
	DB           *gorm.DB      // database connection
	RedisClient  *redis.Client // redis client
	EventBus     bus.EventBus  // event bus
	Logger       *logger.Logger
	Enforcer     *casbin.Enforcer
	JWTService   security.IJWTService
	Repositories RepositoryContainer

	UserModule             UserModule
	AuthModule             AuthModule
	ApiModule              ApiModule
	UploadModule           UploadModule
	DictionaryModule       DictionaryModule
	DictionaryDetailModule DictionaryDetailModule
	MenuModule             MenuModule
	MenuBtnModule          MenuBtnModule
	MenuGroupModule        MenuGroupModule
	MenuParameterModule    MenuParameterModule
	OperationModule        OperationModule
	RoleModule             RoleModule
	FileModule             FileModule
	ScheduledTaskModule    ScheduledTaskModule
	TaskExecutionLogModule TaskExecutionLogModule
}
type RepositoryContainer struct {
	RoleMenuRepository         role_menu.ISysRoleMenuRepository
	CasBinRepository           casbin_rule.ICasbinRuleRepository
	MenuRepository             base_menu.MenuRepositoryInterface
	RoleBtnRepository          role_btn.ISysRoleBtnRepository
	UserRoleRepository         user_role.ISysUserRoleRepository
	JwtBlacklistRepository     jwt_blacklist.JwtBlacklistRepository
	ApiRepository              api.ApiRepositoryInterface
	DictionaryDetailRepository dictionary_detail.DictionaryRepositoryInterface
	DictionaryRepository       dictionary.DictionaryRepositoryInterface
	MenuGroupRepository        base_menu_group.MenuGroupRepositoryInterface
	MenuBtnRepository          base_menu_btn.MenuBtnRepositoryInterface
	MenuParameterRepository    base_menu_parameter.MenuParameterRepositoryInterface
	RoleRepository             role.ISysRolesRepository
	UserRepository             user.UserRepositoryInterface
	FileRepository             files.ISysFilesRepository
}

// SetupDependencies creates a new application context with all dependencies
func SetupDependencies(loggerInstance *logger.Logger) (*ApplicationContext, error) {
	// Initialize database with logger
	db, err := psql.InitPSQLDB(loggerInstance)
	if err != nil {
		return nil, err
	}

	// create event bus
	eventBus := factory.CreateEventBus(loggerInstance)

	// Initialize Redis client
	redisClientInstance, err := redisClient.InitRedisClient(loggerInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	// 初始化Casbin执行器
	enforcer, err := sharedUtil.InitCasbinEnforcer(db, loggerInstance)
	if err != nil {
		return nil, err
	}

	// 初始化共享的repositories
	repositories := RepositoryContainer{
		RoleMenuRepository:      role_menu.NewSysRoleMenuRepository(db, loggerInstance),
		CasBinRepository:        casbin_rule.NewCasbinRuleRepository(db, loggerInstance),
		MenuRepository:          base_menu.NewMenuRepository(db, loggerInstance),
		RoleBtnRepository:       role_btn.NewRoleBtnRepository(db, loggerInstance),
		UserRoleRepository:      user_role.NewSysUserRoleRepository(db, loggerInstance),
		JwtBlacklistRepository:  jwt_blacklist.NewUJwtBlacklistRepository(db),
		DictionaryRepository:    dictionary.NewDictionaryRepository(db, loggerInstance),
		MenuBtnRepository:       base_menu_btn.NewMenuBtnRepository(db, loggerInstance),
		MenuGroupRepository:     base_menu_group.NewMenuGroupRepository(db, loggerInstance),
		MenuParameterRepository: base_menu_parameter.NewMenuParameterRepository(db, loggerInstance),
		RoleRepository:          role.NewSysRolesRepository(db, loggerInstance),
		UserRepository:          user.NewUserRepository(db, loggerInstance),
		FileRepository:          files.NewSysFilesRepository(db, loggerInstance),
	}

	// Initialize JWT service
	jwtService := security.NewJWTService()

	// 创建应用上下文
	appContext := &ApplicationContext{
		DB:           db,
		RedisClient:  redisClientInstance,
		EventBus:     eventBus,
		Logger:       loggerInstance,
		Enforcer:     enforcer,
		JWTService:   jwtService,
		Repositories: repositories,
	}

	// 定义模块初始化函数切片
	moduleSetupFuncs := []func(*ApplicationContext) error{
		setupUserModule,
		setupAuthModule,
		setupApiModule,
		setupMenuModule,
		setupRoleModule,
		setupDictionaryModule,
		setupDictionaryDetailModule,
		setupMenuGroupModule,
		setupMenuBtnModule,
		setupOperationModule,
		setupUploadModule,
		setupMenuParameterModule,
		setupFileModule,
		setupScheduledTaskModule,
		setupTaskExecutionLogModule,
	}

	for _, setupFunc := range moduleSetupFuncs {
		if err := setupFunc(appContext); err != nil {
			return nil, fmt.Errorf("failed to setup module: %w", err)
		}
	}

	return appContext, nil
}

// close response
func (appContext *ApplicationContext) Close() error {
	// close database connection
	if appContext.DB != nil {
		db, _ := appContext.DB.DB()
		if err := db.Close(); err != nil {
			appContext.Logger.Error("Error closing database connection", zap.Error(err))
		}
	}

	// down redis client
	if appContext.RedisClient != nil {
		if err := appContext.RedisClient.Close(); err != nil {
			appContext.Logger.Error("Error closing Redis connection", zap.Error(err))
		} else {
			appContext.Logger.Info("Redis connection closed successfully")
		}
	}

	return nil
}
