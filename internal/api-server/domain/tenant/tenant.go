package tenant

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type TenantDomain struct {
	TenantRepo repo.TenantRepo
}

func NewTenantDomain(tenantRepo repo.TenantRepo) *TenantDomain {
	return &TenantDomain{
		TenantRepo: tenantRepo,
	}
}

func (ts *TenantDomain) CreateOwnerTenantIfNotExistsTx(ctx context.Context, tx *gorm.DB, account *model.Account, isSetup bool) error {
	tenantJoin, err := ts.TenantRepo.FindTenantJoinByAccount(ctx, account, true, tx)

	if err != nil {
		return err
	}

	if tenantJoin.ID != "" {
		return nil
	}

	name := fmt.Sprintf("%s's Workspace", account.Name)

	tenant, err := ts.CreateTenantTx(ctx, tx, account, &model.Tenant{Name: name}, false)

	if err != nil {
		return err
	}

	encryptPublicKey, err := util.GenerateKeyPair(&util.FileStorage{}, tenant.ID)

	if err != nil {
		return errors.WithCode(code.ErrRSAGenerate, err.Error())
	}

	tenant.EncryptPublicKey = encryptPublicKey

	tenant, err = ts.TenantRepo.UpdateEncryptPublicKey(ctx, tenant, true, tx)

	if err != nil {
		return err
	}

	if _, err := ts.CreateTenantMemberTx(ctx, tx, account, tenant, string(model.OWNER)); err != nil {
		return err
	}

	return nil
}

func (ts *TenantDomain) CreateTenantTx(ctx context.Context, tx *gorm.DB, account *model.Account, tenant *model.Tenant, isSetup bool) (*model.Tenant, error) {
	return ts.TenantRepo.CreateOwnerTenant(ctx, tenant, account, true, tx)
}

func (ts *TenantDomain) CreateTenantMemberTx(ctx context.Context, tx *gorm.DB, account *model.Account, tenant *model.Tenant, role string) (*model.TenantAccountJoin, error) {
	var (
		tenantAccountJoin *model.TenantAccountJoin
	)
	hasRole, err := ts.TenantRepo.HasRoles(ctx, tenant, []model.TenantAccountJoinRole{model.OWNER}, true, tx)

	if err != nil {
		return nil, err
	}

	if hasRole {
		return nil, errors.WithCode(code.ErrTenantAlreadyExist, "tenant %s already exists role %s", tenant.Name, role)
	}

	tenantMember, err := ts.TenantRepo.GetTenantJoinOfAccount(ctx, tenant, account, true, tx)

	if err != nil {
		return nil, err
	}

	if tenantMember.ID != "" {
		tenantMember.Role = role
		tenantMember, err := ts.TenantRepo.UpdateRoleTenantOfAccount(ctx, tenantMember, true, tx)
		if err != nil {
			return nil, err
		}
		return tenantMember, nil
	}

	if tenantAccountJoin, err = ts.TenantRepo.CreateCurrentTenantJoinOfAccount(ctx, tenant, account, model.TenantAccountJoinRole(role), true, tx); err != nil {
		return nil, err
	}

	return tenantAccountJoin, nil
}
