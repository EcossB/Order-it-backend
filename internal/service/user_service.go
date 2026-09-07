package service

import (
	"context"
	"errors"
	"strings"

	"order-it-backend/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        domain.UserRepository
	tenantRepo  domain.TenantRepository
	roleRepo    domain.UserRoleRepository
	kitchenRepo domain.KitchenRepository
}

func NewUserService(repo domain.UserRepository, tenantRepo domain.TenantRepository, roleRepo domain.UserRoleRepository, kitchenRepo domain.KitchenRepository) *UserService {
	return &UserService{
		repo:        repo,
		tenantRepo:  tenantRepo,
		roleRepo:    roleRepo,
		kitchenRepo: kitchenRepo,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User, rawPassword *string, rawPin *string) (*domain.User, error) {

	if user.TenantId == uuid.Nil {
		return nil, errors.New("invalid tenant id")
	}

	if user.Role == uuid.Nil {
		return nil, errors.New("invalid role id")
	}

	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if user.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*user.Email))
		if email == "" {
			return nil, errors.New("email cannot be empty if provided")
		}
		user.Email = &email
	}

	// Validate Tenant exists
	_, err := s.tenantRepo.GetById(ctx, user.TenantId)
	if err != nil {
		return nil, errors.New("tenant does not exist")
	}

	// Validate Role exists and belongs to Tenant
	role, err := s.roleRepo.GetById(ctx, user.Role)
	if err != nil {
		return nil, errors.New("role does not exist")
	}
	if role.TenantId != user.TenantId {
		return nil, errors.New("role does not belong to this tenant")
	}

	// Validate Kitchen exists if provided
	if user.KitchenId != nil {
		kitchen, err := s.kitchenRepo.GetById(ctx, *user.KitchenId)
		if err != nil {
			return nil, errors.New("kitchen does not exist")
		}
		if kitchen.TenantId != user.TenantId {
			return nil, errors.New("kitchen does not belong to this tenant")
		}
	} else if role.Description == "CHEF" {
		return nil, errors.New("kitchen is required for CHEF role")
	}

	// Validate Name uniqueness
	existingUser, err := s.repo.GetByTenantIdAndName(ctx, user.TenantId, user.Name)
	if err == nil && existingUser != nil {
		return nil, errors.New("a user with this name already exists in this tenant")
	}

	// Validate Email uniqueness
	if user.Email != nil {
		existingEmailUser, err := s.repo.GetByTenantIdAndEmail(ctx, user.TenantId, *user.Email)
		if err == nil && existingEmailUser != nil {
			return nil, errors.New("a user with this email already exists in this tenant")
		}
	}

	// Hash Password and Pin
	if rawPassword != nil && *rawPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*rawPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		hashStr := string(hash)
		user.PasswordHash = &hashStr
	}

	if rawPin != nil && *rawPin != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*rawPin), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash pin")
		}
		hashStr := string(hash)
		user.PinHash = &hashStr
	}

	if user.PasswordHash == nil && user.PinHash == nil {
		return nil, errors.New("at least a password or a pin must be provided")
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetById(ctx, id)
}

func (s *UserService) GetAllByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.User, error) {
	if tenantId == uuid.Nil {
		return nil, errors.New("invalid tenant id")
	}
	return s.repo.GetAllByTenantId(ctx, tenantId)
}

func (s *UserService) Update(ctx context.Context, user *domain.User, rawPassword *string, rawPin *string) error {
	if user.Id == uuid.Nil {
		return errors.New("invalid user id")
	}

	// Validate Name
	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return errors.New("name cannot be empty")
	}

	if user.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*user.Email))
		if email == "" {
			return errors.New("email cannot be empty if provided")
		}
		user.Email = &email
	}

	// Validate Role
	role, err := s.roleRepo.GetById(ctx, user.Role)
	if err != nil {
		return errors.New("role does not exist")
	}
	if role.TenantId != user.TenantId {
		return errors.New("role does not belong to this tenant")
	}

	// Validate Kitchen
	if user.KitchenId != nil {
		kitchen, err := s.kitchenRepo.GetById(ctx, *user.KitchenId)
		if err != nil {
			return errors.New("kitchen does not exist")
		}
		if kitchen.TenantId != user.TenantId {
			return errors.New("kitchen does not belong to this tenant")
		}
	} else if role.Description == "CHEF" {
		return errors.New("kitchen is required for CHEF role")
	}

	// Retrieve original user to keep hashes if not updating
	originalUser, err := s.repo.GetById(ctx, user.Id)
	if err != nil {
		return errors.New("user not found")
	}

	user.PasswordHash = originalUser.PasswordHash
	user.PinHash = originalUser.PinHash

	// Update Hashes if new ones are provided
	if rawPassword != nil && *rawPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*rawPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		hashStr := string(hash)
		user.PasswordHash = &hashStr
	}

	if rawPin != nil && *rawPin != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*rawPin), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash pin")
		}
		hashStr := string(hash)
		user.PinHash = &hashStr
	}

	// Ensure at least one authentication method remains
	if user.PasswordHash == nil && user.PinHash == nil {
		return errors.New("at least a password or a pin must be provided")
	}

	// Validate Name uniqueness
	existingUser, err := s.repo.GetByTenantIdAndName(ctx, user.TenantId, user.Name)
	if err == nil && existingUser != nil && existingUser.Id != user.Id {
		return errors.New("a user with this name already exists in this tenant")
	}

	// Validate Email uniqueness
	if user.Email != nil {
		existingEmailUser, err := s.repo.GetByTenantIdAndEmail(ctx, user.TenantId, *user.Email)
		if err == nil && existingEmailUser != nil && existingEmailUser.Id != user.Id {
			return errors.New("a user with this email already exists in this tenant")
		}
	}

	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid user id")
	}
	return s.repo.Delete(ctx, id)
}
