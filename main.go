package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")

	r.GET("/", func(c *gin.Context) {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		res, err := getCalendar(ctx)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title":  "ASX Calendar Crawler",
				"header": res.Header,
				"body":   res.Body,
			})
		}
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/calendar", func(c *gin.Context) {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		res, err := getCalendar(ctx)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": res,
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

type asxCalendar struct {
	Header []string
	Body   [][]string
}

func getCalendar(ctx context.Context) (asxCalendar, error) {
	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(`https://www.asx.com.au/markets/market-resources/trading-hours-calendar/cash-market-trading-hours/trading-calendar`)); err != nil {
		return asxCalendar{}, fmt.Errorf("could not navigate to github: %v", err)
	}

	// Đếm số lượng col
	var colCount int
	chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelectorAll(".scroll-container > .table-asx > thead > tr > td").length`, &colCount))

	// Đếm số lượng row
	var rowCount int
	chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelectorAll(".scroll-container > .table-asx > tbody > tr").length`, &rowCount))

	var cellHeaders []string
	// Lấy giá trị của các headers từ bảng
	for i := 1; i <= colCount; i++ {
		var cellHeader string
		chromedp.Run(ctx,
			chromedp.Text(fmt.Sprintf(`#multi-column-0 > div > div.table-component > div > div > div > div > div.right-container > div.scroll-container > table > thead > tr > td:nth-child(%d)`, i), &cellHeader, chromedp.NodeVisible),
		)

		cellHeaders = append(cellHeaders, strings.TrimSpace(cellHeader))
	}

	// Lấy giá trị từ bảng
	var cellBodies [][]string
	for i := 1; i <= rowCount; i++ {
		var cellRow []string
		for j := 1; j <= colCount; j++ {
			var cellValue string
			queryString := fmt.Sprintf(`#multi-column-0 > div > div.table-component > div > div > div > div > div.right-container > div.scroll-container > table > tbody > tr:nth-child(%d) > td:nth-child(%d)`, i, j)
			chromedp.Run(ctx,
				chromedp.Text(queryString, &cellValue, chromedp.NodeVisible),
			)
			cellRow = append(cellRow, strings.TrimSpace(cellValue))
		}
		cellBodies = append(cellBodies, cellRow)
	}

	res := asxCalendar{
		Header: cellHeaders,
		Body:   cellBodies,
	}

	return res, nil
}
