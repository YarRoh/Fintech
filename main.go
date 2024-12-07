package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResidencyStatus string

const (
	MRP                         = 3692
	Resident    ResidencyStatus = "Резидент"
	NonResident ResidencyStatus = "Нерезидент"
)

func LoggerMiddleware(c *gin.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path
	c.Next()

	status := c.Writer.Status()
	fmt.Printf("%s %s -> %d\n", method, path, status)
}

type SalaryRequest struct {
	Income    float64         `json:"income"`
	Expenses  float64         `json:"expenses"`
	Employees ResidencyStatus `json:"employees"`
	Turnover  float64         `json:"turnover"`
	SocialTax float64         `json:"socialTax"`
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.Use(LoggerMiddleware)

	router.LoadHTMLGlob("templates/*")

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

	router.GET("/calculate_ipn", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calculate_ipn.html", nil)
	})

	router.GET("/calculate_cn", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calculate_social_tax.html", nil)
	})

	router.POST("/calculate_ip", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Пример расчета для ИП
		income := req.Income

		opv := income * 0.1
		if opv > 425000 {
			opv = 425000
		}

		vosms := 5950
		

		// СоцОтчисления
		socialTax := income * 0.035
		if socialTax < 2975 {
			socialTax = 2975
		}

		//Обязательные пенсионные платежи работодателя
		pensionTax := income * 0.015

		c.JSON(http.StatusOK, gin.H{
			"opv":        opv,
			"vosms":      vosms,
			"socialTax":  socialTax,
			"pensionTax": pensionTax,
		})
	})

	router.POST("/calculate_employees", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Доход работника
		income := req.Income
		//OПВ
		opv := income * 0.1
		//ВОСМС
		vosms := income * 0.02
		//ИПН
		var ipn float64
		if req.Employees == Resident {
			ipn = (income - opv- vosms - MRP*14) * 0.1
			ipn = math.Round(ipn)
		} else if req.Employees == NonResident {
			ipn = income * 0.1
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid residency status"})
			return
		}

		//ОПВР
		opvr := income * 0.015
		//СО
		so := income * 0.035
		//ООСМС
		oosms := income * 0.03

		//Сумма на руки
		netIncome := income - opv - ipn - vosms
		c.JSON(http.StatusOK, gin.H{
			"opv":       opv,
			"ipn":       ipn,
			"vosms":     vosms,
			"opvr":      opvr,
			"so":        so,
			"oosms":     oosms,
			"netIncome": netIncome,
		})
	})

	router.POST("/calculate_ipn", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		turnover := req.Turnover

		ipn := turnover * 0.015
		c.JSON(http.StatusOK, gin.H{
			"ipn": ipn,
		})
	})

	router.POST("/calculate_cn", func(c *gin.Context) {
		var req SalaryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		turnover := req.Turnover
		//Социальный налог за полгода
		socialTax := (turnover * 0.015) - req.SocialTax
		if socialTax < 0 {
			socialTax = 0
		}
		c.JSON(http.StatusOK, gin.H{
			"socialTax": socialTax,
		})
	})

	router.Run(":8080")
}
