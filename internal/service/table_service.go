package service

import (
	"context"
	"errors"
	"order-it-backend/internal/domain"
	"strings"

	"github.com/google/uuid"
)

type TableService struct {
	repo             domain.TableRepository
	tenantRepository domain.TenantRepository
}

func NewTableService(repo domain.TableRepository, tenantRepo domain.TenantRepository) *TableService {
	return &TableService{repo: repo, tenantRepository: tenantRepo}
}

func (t *TableService) CreateTable(ctx context.Context, table *domain.Table) (*domain.Table, error) {

	if table.TenantId == uuid.Nil {
		return nil, errors.New("el id de la organización es invalido")
	}

	table_number := strings.TrimSpace(table.TableNumber)

	if table_number == "" {
		return nil, errors.New("el numero de la mesa debe de ser valido")
	}

	err := t.repo.Create(ctx, table)

	if err != nil {
		return nil, err
	}

	return table, nil

}

func (t *TableService) GetTableById(ctx context.Context, id uuid.UUID) (*domain.Table, error) {

	if id == uuid.Nil {
		return nil, errors.New("el id de la mesa inválido")
	}

	return t.repo.GetById(ctx, id)
}

func (t *TableService) GetTableByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*domain.Table, error) {

	if tenantId == uuid.Nil {
		return nil, errors.New("el id de la organización es invalido")
	}

	return t.repo.GetAllByTenantId(ctx, tenantId)

}

func (t *TableService) UpdateTableStatus(ctx context.Context, table *domain.Table) error {

	//Validamos que el id de la mesa y de la organizacion no sean nulos.
	if table.Id == uuid.Nil || table.TenantId == uuid.Nil {
		return errors.New("el ID de la cocina y la organización son obligatorios")
	}

	// validamos que el estatus que esta pasando sea el correcto.
	if table.Status != "FREE" && table.Status != "OCCUPIED" && table.Status != "WAITING_FOOD" {
		return errors.New("el estatus de la mesa no es valido")
	}

	_, err := t.tenantRepository.GetById(ctx, table.TenantId)

	if err != nil {
		return errors.New("el id de la organización es invalido")
	}

	return t.repo.UpdateStatus(ctx, table)
}

func (t *TableService) UpdateTable(ctx context.Context, table *domain.Table) error {

	//Validamos que el id de la mesa y de la organizacion no sean nulos.
	if table.Id == uuid.Nil || table.TenantId == uuid.Nil {
		return errors.New("el ID de la cocina y la organización son obligatorios")
	}

	// validamos que el estatus que esta pasando sea el correcto.
	tableNumber := strings.TrimSpace(table.TableNumber)

	if tableNumber == "" {
		return errors.New("el numero de la mesa de ser valido")
	}

	_, err := t.tenantRepository.GetById(ctx, table.TenantId)

	if err != nil {
		return errors.New("el id de la organización es invalido")
	}

	existingTable, err := t.repo.GetByTenantIdAndName(ctx, table.TenantId, table.TableNumber)

	if err != nil {
		return err
	}

	if existingTable != nil && existingTable.Id != table.Id {
		return errors.New("ya existe otra mesa con este nombre en la organización")
	}

	return t.repo.Update(ctx, table)
}

func (t *TableService) DeleteTableById(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id de cocina inválido")
	}
	return t.repo.Delete(ctx, id)
}
