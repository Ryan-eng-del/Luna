package dao

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"gorm.io/gorm"
)

type AccountDao struct {
	db *gorm.DB
}

func NewAccountDao(db *gorm.DB) *AccountDao {
	return &AccountDao{db: db}
}

var _ repo.AccountRepo = (*AccountDao)(nil)

func (ad *AccountDao) GetAccountByEmail(context context.Context, email string) (*model.Account, error) {
	var account model.Account

	if err := ad.db.Limit(1).Find(&account, "email = ?", email).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &account, nil
}

func (ad *AccountDao) CreateAccount(context context.Context, account *model.Account) (*model.Account, error) {
	if err := ad.db.Create(account).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return account, nil
}

func (ad *AccountDao) UpdateAccountIpAddress(context context.Context, account *model.Account) error {
	if err := ad.db.Model(account).Where("id = ?", account.ID).Select("last_login_at", "last_login_ip").Updates(account).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ad *AccountDao) UpdateAccountStatus(context context.Context, account *model.Account) error {
	if err := ad.db.Model(account).Where("id = ?", account.ID).Update("status", account.Status).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}
