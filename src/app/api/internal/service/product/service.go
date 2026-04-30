package product

import (
	"database/sql"
	"fmt"
	"log"

	"seo-backend/internal/database"
	"seo-backend/internal/models"
)

type Service struct {
	db           *sql.DB
	repo         *Repository
	adapterRepo  *AdapterConfigRepo
	queryBuilder *QueryBuilder
	auth         *Authorizer
	txManager    *TransactionManager
	connTester   *ConnectionTester
}

func NewService() *Service {
	db := database.GetDB()

	return &Service{
		db:           db,
		repo:         NewRepository(db),
		adapterRepo:  NewAdapterConfigRepo(db),
		queryBuilder: NewQueryBuilder(),
		auth:         NewAuthorizer(),
		txManager:    NewTransactionManager(db),
		connTester:   NewConnectionTester(),
	}
}

// Public methods
func (s *Service) GetAll(ctx UserContext, filters ProductFilters) ([]models.Product, error) {
	query, args := s.queryBuilder.BuildListQuery(ctx, filters)
	products, err := s.repo.GetAll(query, args)
	if err != nil {
		return nil, err
	}

	// Load adapter configs for each product
	for i := range products {
		if err := s.adapterRepo.LoadForProduct(&products[i]); err != nil {
			log.Printf("Warning: failed to load config for product %s: %v", products[i].ID, err)
		}
	}

	return products, nil
}

func (s *Service) GetByID(id string, ctx UserContext) (*models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	if !s.auth.CanAccess(product, ctx) {
		return nil, fmt.Errorf("access denied")
	}

	if err := s.adapterRepo.LoadForProduct(product); err != nil {
		log.Printf("Warning: failed to load adapter config: %v", err)
	}

	return product, nil
}

func (s *Service) Create(req models.CreateProductRequest, ctx UserContext) (*models.Product, error) {
	productID, err := s.txManager.CreateProduct(req, ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(productID)
}

func (s *Service) Update(id string, updates map[string]interface{}, ctx UserContext) error {
	product, err := s.GetByID(id, ctx)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product not found")
	}

	return s.txManager.UpdateProduct(id, updates, ctx)
}

func (s *Service) Delete(id string, ctx UserContext) error {
	product, err := s.GetByID(id, ctx)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product not found")
	}

	return s.repo.Delete(id)
}

func (s *Service) TestConnection(id string) (*ConnectionTestResult, error) {
	product, err := s.repo.GetProductBasicInfo(id)
	if err != nil {
		return nil, err
	}

	config := s.adapterRepo.GetOrDefault(id)

	result := s.connTester.Test(product, config)

	s.repo.UpdateConnectionStatus(id, result.Success)

	return result, nil
}
