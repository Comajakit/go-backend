package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "go-backend/database"
	"go-backend/database/models"
	itf "go-backend/interfaces"
	"go-backend/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePort(c *gin.Context) {
	var req itf.CreatePortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	username, err := getNameFromToken(c)
	if err != nil {
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err = db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		// Handle the error, e.g., return an error response
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	port := models.UserPort{
		PortName: req.PortName,
		UserID:   user.ID,
	}

	if err := db.DB.Create(&port).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := map[string]interface{}{
		"ID":       port.ID,
		"portName": port.PortName,
		"userID":   user.ID,
	}
	c.JSON(http.StatusCreated, response)

}

func AddPortStrategy(c *gin.Context) {
	var req itf.CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	username, err := getNameFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	port, err := getUserPortByIDAndName(user.ID, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User port not found"})
		return
	}

	portStrategy := models.PortStrategyProfile{
		StrategyName: req.StrategyName,
		OwnerID:      user.ID,
		UserPortID:   port.ID,
	}

	if err := db.DB.Create(&portStrategy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalValue := 0.0
	for _, theme := range req.Theme {
		for _, themeValue := range theme {
			floatValue, err := strconv.ParseFloat(themeValue, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theme value"})
				return
			}

			totalValue += floatValue
		}
	}
	var themeInfoList []types.ThemeInfo
	for _, theme := range req.Theme {
		for themeName, themeValue := range theme {
			floatValue, err := strconv.ParseFloat(themeValue, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theme value"})
				return
			}

			percent := floatValue / totalValue

			themePair := models.ThemePercentagePair{
				StrategyProfileID: portStrategy.ID,
				Theme:             themeName,
				Percentage:        percent,
			}
			if err := db.DB.Create(&themePair).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			percentageStr := strconv.FormatFloat(percent*100, 'f', 2, 64)
			themeInfo := types.ThemeInfo{
				ThemeName:  themeName,
				Percentage: percentageStr + "%",
			}
			themeInfoList = append(themeInfoList, themeInfo)
		}
	}

	response := map[string]interface{}{
		"ID":           portStrategy.ID,
		"portName":     port.PortName,
		"strategyName": portStrategy.StrategyName,
		"userID":       user.ID,
		"portID":       port.ID,
		"totalTheme":   len(themeInfoList),
		"themeInfo":    themeInfoList,
	}

	c.JSON(http.StatusCreated, response)
}

func AddStock(c *gin.Context) {
	var req itf.AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	username, err := getNameFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	port, err := getUserPortByIDAndName(user.ID, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User port not found"})
		return
	}

	stocks, err := getStocksByPortID(port.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can't find stocks"})
		return
	}

	stocksLen := len(stocks)
	totalSum := 0.0
	if stocksLen > 0 {
		for _, stock := range stocks {
			totalSum += stock.Total
		}
	}

	for _, stock := range req.Stock {
		averagePrice, err := strconv.ParseFloat(stock.AveragePrice, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid average price value"})
			return
		}
		totalSum += averagePrice * float64(stock.Volume)
	}

	if stocksLen > 0 {
		for i := range stocks {
			stocks[i].PercentageInPort = (stocks[i].Total / totalSum) * 100
			stocks[i].DivPercentPort = (stocks[i].ExpectedDivReturn / totalSum) * 100
		}
	}

	for _, stock := range req.Stock {
		divPerShare, err := strconv.ParseFloat(stock.DivPerShare, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dividend per share value"})
			return
		}

		averagePrice, err := strconv.ParseFloat(stock.AveragePrice, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid average price value"})
			return
		}
		total := float64(stock.Volume) * averagePrice
		divInPercent := (divPerShare / averagePrice) * 100
		expectedDivReturn := divPerShare * float64(stock.Volume)
		stockInfo := models.PortStock{
			OwnerId:           user.ID,
			UserPortID:        port.ID,
			StockSymbol:       stock.Symbol,
			DivPerShare:       divPerShare,
			Total:             total,
			ExpectedDivReturn: expectedDivReturn,
			PercentageInPort:  (total / totalSum) * 100,
			DivPercentPort:    (expectedDivReturn / totalSum) * 100,
			AveragePrice:      averagePrice,
			StockType:         stock.Type,
			DivInPercent:      divInPercent,
			Volume:            uint(stock.Volume),

			// Add any additional fields as necessary
		}

		if err := db.DB.Create(&stockInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if stocksLen > 0 {
		if err := db.DB.Save(stocks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	response := map[string]interface{}{
		"ID":       port.ID,
		"portName": port.PortName,
		"userID":   user.ID,
		"stock":    req.Stock,
	}
	c.JSON(http.StatusCreated, response)

}

func getNameFromToken(c *gin.Context) (string, error) {
	token := c.GetHeader("token")
	session := sessions.Default(c)
	usernameInterface := session.Get(token)

	// Check if the token exists in the session
	if usernameInterface == nil {
		return "", errors.New("token not found in session")
	}

	// Perform type assertion to retrieve the username as a string
	username, ok := usernameInterface.(string)
	if !ok {
		return "", errors.New("invalid username type")
	}

	return username, nil
}

func CheckToken(c *gin.Context) {

	token := c.GetHeader("token")
	fmt.Println(token)
	session := sessions.Default(c)

	session.Set("token", "hehehe")

	username := session.Get(token)
	test := session.Get("token")
	fmt.Println(username)
	fmt.Println(test)

	c.JSON(http.StatusCreated, username)

}

func getUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserPortByIDAndName(userID uuid.UUID, portName string) (*models.UserPort, error) {
	var port models.UserPort
	err := db.DB.Where("user_id = ? AND port_name = ?", userID, portName).First(&port).Error
	if err != nil {
		return nil, err
	}
	return &port, nil
}

func getStocksByPortID(portID uuid.UUID) ([]models.PortStock, error) {
	var stocks []models.PortStock
	err := db.DB.Where("user_port_id = ?", portID).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
