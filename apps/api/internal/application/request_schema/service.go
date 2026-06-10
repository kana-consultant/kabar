package request_schema

import (
	"context"
	"seo-backend/internal/domain/request_schema"
)

type serviceImpl struct {
	repo request_schema.Repository
}

func NewService(repo request_schema.Repository) request_schema.Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Create(ctx context.Context, req *request_schema.CreateRequest) (*request_schema.Response, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check exists
	exists, err := s.repo.Exists(ctx, req.ProviderID, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, request_schema.ErrDuplicate
	}

	// Default values
	supportsTemp := true
	if req.SupportsTemperature != nil {
		supportsTemp = *req.SupportsTemperature
	}
	supportsStream := true
	if req.SupportsStreaming != nil {
		supportsStream = *req.SupportsStreaming
	}

	entity := request_schema.NewRequestSchema(
		req.ProviderID,
		req.Name,
		req.EndpointPath,
		req.MaxTokensKey,
		req.SystemRoleKey,
		req.ResponseTextPath,
		req.ResponseImagePath,
		req.RequestTemplate,
		&supportsTemp,
		&supportsStream,
	)

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return request_schema.ToResponse(entity), nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (*request_schema.Response, error) {
	if id == "" {
		return nil, request_schema.ErrInvalidID
	}

	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return request_schema.ToResponse(entity), nil
}

func (s *serviceImpl) GetByProviderAndName(ctx context.Context, providerID string, name string) (*request_schema.Response, error) {
	if providerID == "" {
		return nil, request_schema.ErrInvalidProviderID
	}
	if name == "" {
		return nil, request_schema.ErrInvalidName
	}

	entity, err := s.repo.GetByProviderAndName(ctx, providerID, name)
	if err != nil {
		return nil, err
	}

	return request_schema.ToResponse(entity), nil
}

func (s *serviceImpl) GetAll(ctx context.Context, page, limit int) (*request_schema.ListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	entities, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &request_schema.ListResponse{
		Data:       request_schema.ToResponseList(entities),
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *serviceImpl) GetByProvider(ctx context.Context, providerID string) ([]request_schema.Response, error) {
	if providerID == "" {
		return nil, request_schema.ErrInvalidProviderID
	}

	entities, err := s.repo.GetByProvider(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return request_schema.ToResponseList(entities), nil
}

func (s *serviceImpl) Update(ctx context.Context, id string, req *request_schema.UpdateRequest) (*request_schema.Response, error) {
	if id == "" {
		return nil, request_schema.ErrInvalidID
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		exists, err := s.repo.Exists(ctx, existing.ProviderID, *req.Name)
		if err != nil {
			return nil, err
		}
		if exists && existing.Name != *req.Name {
			return nil, request_schema.ErrDuplicate
		}
		existing.Name = *req.Name
	}

	if req.EndpointPath != nil {
		existing.EndpointPath = *req.EndpointPath
	}
	if req.MaxTokensKey != nil {
		existing.MaxTokensKey = req.MaxTokensKey
	}
	if req.SystemRoleKey != nil {
		existing.SystemRoleKey = req.SystemRoleKey
	}
	if req.ResponseTextPath != nil {
		existing.ResponseTextPath = req.ResponseTextPath
	}
	if req.ResponseImagePath != nil {
		existing.ResponseImagePath = req.ResponseImagePath
	}
	if req.RequestTemplate != nil {
		existing.RequestTemplate = req.RequestTemplate
	}
	if req.SupportsTemperature != nil {
		existing.SupportsTemperature = req.SupportsTemperature
	}
	if req.SupportsStreaming != nil {
		existing.SupportsStreaming = req.SupportsStreaming
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return request_schema.ToResponse(existing), nil
}

func (s *serviceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return request_schema.ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}

func (s *serviceImpl) Validate(ctx context.Context, rs *request_schema.RequestSchema) error {
	if rs.ProviderID == "" {
		return request_schema.ErrInvalidProviderID
	}
	if rs.Name == "" {
		return request_schema.ErrInvalidName
	}
	if rs.EndpointPath == "" {
		return request_schema.ErrInvalidEndpointPath
	}
	return nil
}

func (s *serviceImpl) validateCreateRequest(req *request_schema.CreateRequest) error {
	if req.ProviderID == "" {
		return request_schema.ErrInvalidProviderID
	}
	if req.Name == "" {
		return request_schema.ErrInvalidName
	}
	if req.EndpointPath == "" {
		return request_schema.ErrInvalidEndpointPath
	}
	return nil
}
