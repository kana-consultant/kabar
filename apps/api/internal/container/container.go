package container

import (
	"database/sql"
	"log"
	"os"

	// Application Services
	adapterApp "seo-backend/internal/application/adapter"
	apikeyApp "seo-backend/internal/application/apikey"
	authApp "seo-backend/internal/application/auth"
	dashboardApp "seo-backend/internal/application/dashboard"
	draftApp "seo-backend/internal/application/draft"
	generateApp "seo-backend/internal/application/generate"
	historyApp "seo-backend/internal/application/history"
	modelApp "seo-backend/internal/application/model"
	productApp "seo-backend/internal/application/product"
	providerApp "seo-backend/internal/application/provider"
	teamApp "seo-backend/internal/application/team"
	userApp "seo-backend/internal/application/user"
	"seo-backend/internal/helper"

	// Domain Interfaces
	"seo-backend/internal/domain/adapter"
	"seo-backend/internal/domain/apikey"
	"seo-backend/internal/domain/auth"
	"seo-backend/internal/domain/dashboard"
	"seo-backend/internal/domain/draft"
	"seo-backend/internal/domain/generate"
	model "seo-backend/internal/domain/model"
	"seo-backend/internal/domain/product"
	"seo-backend/internal/domain/provider"
	"seo-backend/internal/domain/team"
	"seo-backend/internal/domain/user"

	"seo-backend/internal/scheduler"

	// Repositories
	"seo-backend/internal/database"
	"seo-backend/internal/infrastructure/ai/builder"
	"seo-backend/internal/infrastructure/ai/parser"
	"seo-backend/internal/infrastructure/db/repositories"
	"seo-backend/internal/infrastructure/http/client"

	// Handlers
	historyBuilder "seo-backend/internal/infrastructure/db/query_builder"
	apiKeyHandler "seo-backend/internal/presentation/handlers/apiKey"
	authHandler "seo-backend/internal/presentation/handlers/auth"
	dashboardHandler "seo-backend/internal/presentation/handlers/dashboard"
	draftHandler "seo-backend/internal/presentation/handlers/draft"
	generateHandler "seo-backend/internal/presentation/handlers/generate"
	historyHandler "seo-backend/internal/presentation/handlers/history"
	modelHandler "seo-backend/internal/presentation/handlers/model"
	productHandler "seo-backend/internal/presentation/handlers/product"
	providerHandler "seo-backend/internal/presentation/handlers/provider"
	teamHandler "seo-backend/internal/presentation/handlers/team"
	userHandler "seo-backend/internal/presentation/handlers/user"

	// Utilities
	"seo-backend/internal/pkg/crypto"
	"seo-backend/internal/pkg/jwt"
)

type Container struct {
	DB *sql.DB

	// Repositories
	ProductRepo       product.ProductRepository
	AdapterConfigRepo adapter.AdapterConfigRepository
	APIKeyRepo        apikey.Repository
	GenerateRepo      generate.Repository
	ProviderRepo      provider.Repository
	UserRepo          user.Repository
	AuthRepo          auth.Repository
	dashboardRepo     dashboard.DashboardRepository
	draftRepo         draft.Repository
	teamRepo          team.Repository
	teamMemberRepo    team.MemberRepository
	model             model.Repository

	// Services
	ProductService       product.ProductService
	AdapterConfigService adapter.AdapterConfigService
	APIKeyService        apikey.Service
	GenerateService      generate.Service
	ProviderService      provider.ProviderService
	UserService          user.UserService
	AuthService          auth.AuthService
	dashboardApp         dashboard.DashboardService
	draftService         draft.Service
	teamService          team.Service
	modelService         model.Service

	// Handlers
	ProductHandler   *productHandler.ProductHandler
	APIKeyHandler    *apiKeyHandler.APIKeyHandler
	GenerateHandler  *generateHandler.GenerateHandler
	ProviderHandler  *providerHandler.ProviderHandler
	UserHandler      *userHandler.UserHandler
	AuthHandler      *authHandler.AuthHandler
	HistoryHandler   *historyHandler.HistoryHandler
	DashboardHandler *dashboardHandler.DashboardHandler
	DraftHandler     *draftHandler.DraftHandler
	TeamHandler      *teamHandler.TeamHandler
	ModelHandler     *modelHandler.AIModelHandler

	// Utilities
	Encryptor    crypto.Encryptor
	JWTGenerator *jwt.Generator
}

func NewContainer(db *sql.DB, httpClient *client.HTTPClient,
	builder *builder.PromptBuilder,
	requestBuilder *builder.RequestBuilder,
	responseParser *parser.ResponseParser,
	redisScheduler *scheduler.RedisScheduler) *Container {
	// Utilities
	encryptor := crypto.NewAESEncryptor()

	jwtSecret := getEnv("JWT_SECRET", "default-secret-change-me")
	jwtExpiry := getEnv("JWT_EXPIRY", "24h")
	jwtGenerator, _ := jwt.NewGenerator(jwtSecret, jwtExpiry)

	if err := database.InitRedis(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer database.CloseRedis()
	qb_history := historyBuilder.NewQueryBuilder()

	// ========== EXISTING REPOSITORIES ==========
	dashboardRepo := repositories.NewDashboardRepository(db)
	productRepo := repositories.NewProductRepository(db)
	adapterConfigRepo := repositories.NewAdapterConfigRepository(db)
	apiKeyRepo := repositories.NewAPIKeyRepository(db)
	generateRepo := repositories.NewRepository(db)
	providerRepo := repositories.NewProviderRepository(db)
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	historyRepo := repositories.NewHistoryRepository(db)
	draftRepo := repositories.NewDraftRepository(db)
	modelRepo := repositories.ModelRepository(db)

	// ========== TEAM REPOSITORIES ==========
	teamRepository := repositories.NewTeamRepository(db)
	memberRepository := repositories.NewMemberRepository(db)

	postService := helper.NewPostService(db)

	// ========== EXISTING QUERY BUILDERS ==========
	userQueryBuilder := userApp.NewQueryBuilder()

	// ========== TEAM QUERY BUILDER ==========
	teamQueryBuilder := teamApp.NewQueryBuilder()

	// ========== EXISTING AUTHORIZERS & VALIDATORS ==========
	userAuthorizer := userApp.NewAuthorizer(db)
	userPasswordService := userApp.NewPasswordService()
	userValidator := userApp.NewValidator(userRepo)

	// ========== TEAM AUTHORIZER & VALIDATOR ==========
	teamAuthorizer := teamApp.NewAuthorizer()
	teamValidator := teamApp.NewValidator(memberRepository)

	//httpClient

	// ========== EXISTING SERVICES ==========
	productService := productApp.NewProductService(db, productRepo, adapterConfigRepo)
	adapterConfigService := adapterApp.NewAdapterConfigService(db, adapterConfigRepo)
	apiKeyService := apikeyApp.NewService(db, apiKeyRepo, encryptor)
	generateService := generateApp.NewService(generateRepo, httpClient, builder, requestBuilder, responseParser)
	providerService := providerApp.NewService(db, providerRepo)
	userService := userApp.NewService(db, userRepo, userQueryBuilder, userAuthorizer, userPasswordService, userValidator)
	authService := authApp.NewService(db, authRepo, jwtGenerator)
	historyService := historyApp.NewService(historyRepo, *qb_history)
	dashboardService := dashboardApp.NewDashboardService(dashboardRepo)
	draftService := draftApp.NewService(draftRepo, redisScheduler, postService, productRepo)
	modelService := modelApp.NewService(modelRepo)

	// ========== TEAM SERVICE ==========
	teamService := teamApp.NewService(
		teamRepository,
		memberRepository,
		teamQueryBuilder,
		teamAuthorizer,
		teamValidator,
	)

	// ========== EXISTING HANDLERS ==========
	productHandler := productHandler.NewProductHandler(productService)
	apiKeyHandler := apiKeyHandler.NewAPIKeyHandler(apiKeyService)
	generateHandler := generateHandler.NewGenerateHandler(generateService)
	providerHandler := providerHandler.NewProviderHandler(providerService)
	userHandler := userHandler.NewUserHandler(userService)
	authHandler := authHandler.NewAuthHandler(authService)
	historyHandler := historyHandler.NewHistoryHandler(historyService)
	dashboardHandler := dashboardHandler.NewDashboardHandler(dashboardService)
	draftHandler := draftHandler.NewDraftHandler(draftService)
	modelHandler := modelHandler.NewAIModelHandler(modelService)

	// ========== TEAM HANDLER ==========
	teamHandler := teamHandler.NewTeamHandler(teamService)

	return &Container{
		DB: db,

		// Repositories
		ProductRepo:       productRepo,
		AdapterConfigRepo: adapterConfigRepo,
		APIKeyRepo:        apiKeyRepo,
		GenerateRepo:      generateRepo,
		ProviderRepo:      providerRepo,
		UserRepo:          userRepo,
		AuthRepo:          authRepo,
		teamRepo:          teamRepository,
		teamMemberRepo:    memberRepository,

		// Services
		ProductService:       productService,
		AdapterConfigService: adapterConfigService,
		APIKeyService:        apiKeyService,
		GenerateService:      generateService,
		ProviderService:      providerService,
		UserService:          userService,
		AuthService:          authService,
		dashboardApp:         dashboardService,
		draftService:         draftService,
		teamService:          teamService,

		// Handlers
		ProductHandler:   productHandler,
		APIKeyHandler:    apiKeyHandler,
		GenerateHandler:  generateHandler,
		ProviderHandler:  providerHandler,
		UserHandler:      userHandler,
		AuthHandler:      authHandler,
		HistoryHandler:   historyHandler,
		DashboardHandler: dashboardHandler,
		DraftHandler:     draftHandler,
		TeamHandler:      teamHandler,
		ModelHandler:     modelHandler,

		// Utilities
		Encryptor:    encryptor,
		JWTGenerator: jwtGenerator,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
