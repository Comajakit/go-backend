package utils

import (
	db "go-backend/database"
	"go-backend/database/models"

	"github.com/google/uuid"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserPortByIDAndName(userID uuid.UUID, portName string) (*models.UserPort, error) {
	var port models.UserPort
	err := db.DB.Where("user_id = ? AND port_name = ?", userID, portName).First(&port).Error
	if err != nil {
		return nil, err
	}
	return &port, nil
}

func GetPortStrategyByPortID(portID uuid.UUID) (*models.PortStrategyProfile, error) {
	var portStrategy models.PortStrategyProfile
	err := db.DB.Where("user_port_id = ?", portID).First(&portStrategy).Error
	if err != nil {
		return nil, err
	}
	return &portStrategy, nil
}

func GetStrategyPairPercentage(strategyID uuid.UUID) ([]models.ThemePercentagePair, error) {
	var pair []models.ThemePercentagePair
	err := db.DB.Where("strategy_profile_id = ?", strategyID).Find(&pair).Error
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func GetStocksByPortID(portID uuid.UUID) ([]models.PortStock, error) {
	var stocks []models.PortStock
	err := db.DB.Where("user_port_id = ?", portID).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func GetStockByPortIDAndSymbol(portID uuid.UUID, symbol string) (models.PortStock, error) {
	var stock models.PortStock
	err := db.DB.Where("user_port_id = ? AND stock_symbol = ?", portID, symbol).First(&stock).Error
	if err != nil {
		return stock, err
	}
	return stock, nil
}
