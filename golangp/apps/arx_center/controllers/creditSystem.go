package controllers

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"

	"github.com/toheart/functrace"
	"gorm.io/gorm"
)

type CreditSystem struct {
	DB *gorm.DB
}

func NewCreditSystem(DB *gorm.DB) *CreditSystem {
	defer functrace.Trace([]interface {
	}{DB})()
	return &CreditSystem{DB}
}

func (cs *CreditSystem) GetBalanceByUserEmail(email string) (int64, error) {
	defer functrace.Trace([]interface {
	}{cs, email})()
	var user models.User
	err := cs.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return 0, err
	}
	logger.Info("Got balance of user", email, "as", user.Balance)
	return user.Balance, nil
}

func (cs *CreditSystem) UpdateBalanceByUserEmail(email string, amount int64) (int64, error) {
	defer functrace.Trace([]interface {
	}{cs, email, amount})()
	var user models.User
	err := cs.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return 0, err
	}
	newBalance := user.Balance + amount
	err = cs.DB.Model(&user).Update("balance", newBalance).Error
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}
