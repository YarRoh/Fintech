package main

import (
	"fmt"
	"math"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/message"
    "golang.org/x/text/language"
)

type ResidencyStatus string

const (
	MRP                         = 3692
	Resident    ResidencyStatus = "Резидент"
	NonResident ResidencyStatus = "Нерезидент"
)
func formatNumber(num float64) string {
	p := message.NewPrinter(language.Russian)
	return p.Sprintf("%.2f", num)
}

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
		vosms := 5950.0

		socialTax := income * 0.035
		if socialTax < 2975 {
			socialTax = 2975
		}

		//Обязательные пенсионные платежи работодателя
		pensionTax := income * 0.015
		formattedopv := formatNumber(opv)
		formattedvosms := formatNumber(vosms)
		formatedsocialTax := formatNumber(socialTax)
		formatedpensionTax := formatNumber(pensionTax)
		
		c.JSON(http.StatusOK, gin.H{
			"opv":        formattedopv,
			"vosms":      formattedvosms,
			"socialTax":  formatedsocialTax,
			"pensionTax": formatedpensionTax,
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
		
		formattedopv := formatNumber(opv)
		formattedvosms := formatNumber(vosms)
		formattedipn := formatNumber(ipn)
		formattedopvr := formatNumber(opvr)
		formattedso := formatNumber(so)
		formattedoosms := formatNumber(oosms)
		formattednetIncome := formatNumber(netIncome)
		
		c.JSON(http.StatusOK, gin.H{
			"opv":       formattedopv,
			"ipn":       formattedipn,
			"vosms":     formattedvosms,
			"opvr":      formattedopvr,
			"so":        formattedso,
			"oosms":     formattedoosms,
			"netIncome": formattednetIncome,
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
			fmt.Println("Ошибка привязки JSON:", err) // Логирование ошибки
			return
		}

		turnover := req.Turnover
		socialTax := (turnover * 0.015) - req.SocialTax
		if socialTax < 0 {
			socialTax = 0
		}
		
		c.JSON(http.StatusOK, gin.H{
			"socialTax": formatNumber(socialTax), // Форматируем ответ
		})
})


	router.Run(":8080")
}


