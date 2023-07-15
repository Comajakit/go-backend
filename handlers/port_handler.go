package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "go-backend/database"
	"go-backend/database/models"
	dao "go-backend/database/utils"
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

	_, user, port, err := getPreRequire(c, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	_, user, port, err := getPreRequire(c, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	stocks, err := dao.GetStocksByPortID(port.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can't find stocks"})
		return
	}

	stocksLen := len(stocks)
	totalSum := 0.0
	totalSum, err = getTotalStockSumWithRequest(stocks, req.Stock)

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
			DivInPercent:      divInPercent,
			Volume:            uint(stock.Volume),

			// Add any additional fields as necessary
		}
		if stock.Type != nil {
			stockInfo.StockType = *stock.Type
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

func UpdateStock(c *gin.Context) {
	var req itf.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	_, user, port, err := getPreRequire(c, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	for _, stock := range req.Stock {
		currentStock, err := dao.GetStockByPortIDAndSymbol(port.ID, stock.Symbol)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Can't find stock"})
			return
		}
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
		currentStock.Volume = uint(stock.Volume)
		currentStock.DivPerShare = divPerShare
		currentStock.AveragePrice = averagePrice
		if stock.Type != nil {
			currentStock.StockType = *stock.Type
		}

		if err := db.DB.Save(currentStock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}
	err = calibratePortByPortID(port.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	response := map[string]interface{}{
		"ID":           port.ID,
		"portName":     port.PortName,
		"userID":       user.ID,
		"stockUpdated": req.Stock,
	}
	c.JSON(http.StatusOK, response)

}

func DeleteStock(c *gin.Context) {
	var req itf.DeleteStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	_, user, port, err := getPreRequire(c, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var deletedSymbols []string
	for _, stock := range req.Stock {
		currentStock, err := dao.GetStockByPortIDAndSymbol(port.ID, stock)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Can't find stock"})
			return
		}

		if err := db.DB.Delete(currentStock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		deletedSymbols = append(deletedSymbols, stock)

	}
	err = calibratePortByPortID(port.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	response := map[string]interface{}{
		"ID":           port.ID,
		"portName":     port.PortName,
		"userID":       user.ID,
		"stockDeleted": deletedSymbols,
	}
	c.JSON(http.StatusOK, response)

}

func SummaryPort(c *gin.Context) {
	var req itf.CheckStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}
	_, _, port, err := getPreRequire(c, req.PortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	stockTypes := make(map[string]float64) // Map to store stock types and their corresponding sum of PercentageInPort
	var stockPorts []models.PortStock
	err = db.DB.Where("user_port_id = ?", port.ID).Find(&stockPorts).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	for _, stockPort := range stockPorts {
		stockTypes[stockPort.StockType] += stockPort.PercentageInPort // Accumulate the sum of PercentageInPort for each stock type

	}

	portStrategy, err := dao.GetPortStrategyByPortID(port.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	pair, err := dao.GetStrategyPairPercentage(portStrategy.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Print the stock types and their corresponding sum of PercentageInPort
	themePercentageMap := make(map[string]float64)
	for _, themePair := range pair {
		themePercentageMap[themePair.Theme] = themePair.Percentage * 100
	}

	summary := make(map[string][]types.StockSummary)
	for stockType, sum := range stockTypes {
		var stockSymbols []string

		// Check if stockType exists in themePercentageMap
		if themePercentage, ok := themePercentageMap[stockType]; ok {
			// Compare stockType sum with corresponding theme percentage
			if sum > themePercentage {
				// Add stock symbols with stockType as "Sell"
				for _, stockPort := range stockPorts {
					if stockPort.StockType == stockType {
						stockSymbols = append(stockSymbols, stockPort.StockSymbol)
					}
				}
				summary["sell"] = append(summary["sell"], types.StockSummary{
					Theme:   stockType,
					Symbol:  stockSymbols,
					Target:  fmt.Sprintf("%.2f%%", themePercentage),
					Current: fmt.Sprintf("%.2f%%", sum),
				})
			} else if sum < themePercentage {
				// Add stock symbols with stockType as "Buy"
				for _, stockPort := range stockPorts {
					if stockPort.StockType == stockType {
						stockSymbols = append(stockSymbols, stockPort.StockSymbol)
					}
				}
				summary["buy"] = append(summary["buy"], types.StockSummary{
					Theme:   stockType,
					Symbol:  stockSymbols,
					Target:  fmt.Sprintf("%.2f%%", themePercentage),
					Current: fmt.Sprintf("%.2f%%", sum),
				})
			}
		} else {
			// No matching stockType in themePercentageMap, assume "SELL"
			for _, stockPort := range stockPorts {
				if stockPort.StockType == stockType {
					stockSymbols = append(stockSymbols, stockPort.StockSymbol)
				}
			}
			summary["sell"] = append(summary["sell"], types.StockSummary{
				Theme:   stockType,
				Symbol:  stockSymbols,
				Target:  "",
				Current: fmt.Sprintf("%.2f%%", sum),
			})
		}
	}

	response := types.SummaryResponse{
		PortName:     port.PortName,
		StrategyName: portStrategy.StrategyName,
		Summary:      summary,
	}

	c.JSON(http.StatusOK, response)
}

func getTotalStockSumWithRequest(currentStock []models.PortStock, reqStock []itf.StockInfo) (float64, error) {
	stocksLen := len(currentStock)
	totalSum := 0.0

	if stocksLen > 0 {
		for _, stock := range currentStock {
			totalSum += stock.Total
		}
	}

	for _, stock := range reqStock {
		averagePrice, err := strconv.ParseFloat(stock.AveragePrice, 64)
		if err != nil {
			return 0.0, err
		}
		totalSum += averagePrice * float64(stock.Volume)
	}

	return totalSum, nil
}

func getTotalStockSum(currentStock []models.PortStock) (float64, error) {
	stocksLen := len(currentStock)
	totalSum := 0.0

	if stocksLen > 0 {
		for _, stock := range currentStock {
			totalSum += stock.Total
		}
	}

	return totalSum, nil
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

func getPreRequire(c *gin.Context, portName string) (string, *models.User, *models.UserPort, error) {
	username, err := getNameFromToken(c)
	if err != nil {
		return "", nil, nil, err
	}

	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return "", nil, nil, err
	}

	port, err := dao.GetUserPortByIDAndName(user.ID, portName)
	if err != nil {
		return "", nil, nil, err
	}

	return username, user, port, nil

}

func calibratePortByPortID(portID uuid.UUID) error {
	var stocks []models.PortStock
	err := db.DB.Where("user_port_id = ?", portID).Find(&stocks).Error
	if err != nil {
		return err
	}
	for i := range stocks {
		stock := &stocks[i]

		// Calculate the derived stock values
		stock.Total = float64(stock.Volume) * stock.AveragePrice
		stock.ExpectedDivReturn = float64(stock.Volume) * stock.DivPerShare
		stock.DivInPercent = (stock.DivPerShare / stock.AveragePrice) * 100

		// Save the updated stock to the database
		if err := db.DB.Save(stock).Error; err != nil {
			return err
		}
	}

	// Retrieve the updated stocks
	err = db.DB.Where("user_port_id = ?", portID).Find(&stocks).Error
	if err != nil {
		return err
	}

	// Calculate the total stock sum
	totalSum, err := getTotalStockSum(stocks)
	if err != nil {
		return err
	}

	// Update the derived percentage values and save them to the database
	for i := range stocks {
		stock := &stocks[i]
		stock.PercentageInPort = (stock.Total / totalSum) * 100
		stock.DivPercentPort = (stock.ExpectedDivReturn / totalSum) * 100

		// Save the updated stock to the database
		if err := db.DB.Save(stock).Error; err != nil {
			return err
		}
	}

	return nil
}
