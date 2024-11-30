package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

func LoggerMiddleware(c *gin.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path
	c.Next()

	status := c.Writer.Status()
	fmt.Printf("%s %s -> %d\n", method, path, status)
}

type SalaryRequest struct {
	Income    float64 `json:"income"`
	Expenses  float64 `json:"expenses"`
	Employees int     `json:"employees"`
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.Use(LoggerMiddleware)

	router.LoadHTMLGlob("template/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Name": "Guest",
		})
	})

	router.GET("/calculate_ip", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calculate_ip.html", nil)
	})

	router.GET("/calculate_employees", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calculate_employees.html", nil)
	})

	router.GET("/calculate_kpi", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calculate_kpi.html", nil)
	})

	router.GET("/calculate_social_tax", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calculate_social_tax.html", nil)
	})

	router.POST("/calculate_ip", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		income := req.Income
		expenses := req.Expenses

		netIncome := income - expenses
		ipTax := netIncome * 0.06 // УСН 6%
		insuranceContributions := netIncome * 0.3 // Страховые взносы 30%

		c.JSON(http.StatusOK, gin.H{
			"netIncome":              netIncome,
			"ipTax":                  ipTax,
			"insuranceContributions": insuranceContributions,
		})
	})

	router.POST("/calculate_employees", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		income := req.Income
		expenses := req.Expenses
		employees := req.Employees

		netIncome := income - expenses
		employeeSalary := netIncome / float64(employees)
		employeeTax := employeeSalary * 0.13 // НДФЛ 13%

		c.JSON(http.StatusOK, gin.H{
			"employeeSalary": employeeSalary,
			"employeeTax":    employeeTax,
		})
	})

	router.POST("/calculate_kpi", func(c *gin.Context) {
		// Пример расчета KPI
		kpi := 85.0 // Пример значения KPI

		c.JSON(http.StatusOK, gin.H{
			"kpi": kpi,
		})
	})

	router.POST("/calculate_social_tax", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		income := req.Income
		expenses := req.Expenses

		netIncome := income - expenses
		socialTax := netIncome * 0.3 // Пример расчета социального налога

		c.JSON(http.StatusOK, gin.H{
			"socialTax": socialTax,
		})
	})

	router.Run(":8080")
}