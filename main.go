package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const brokerURL = ""

var (
	CeId          = flag.String("Ce-Id", uuid.NewString(), "ce-id")
	CeSpecversion = flag.String("Ce-Specversion", "1.0", "ce-specversion")
	ContentType   = flag.String("Content-Type", "application/json", "content-type")
)

type requestBody struct {
	Message string `json:"message"`
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/seed", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		var requestBody = new(requestBody)
		if err := c.Bind(&requestBody); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, brokerURL, getParsedRequestBody(requestBody))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		headers := map[string]string{
			"Ce-Id":          *CeId,
			"Ce-Specversion": *CeSpecversion,
			"Ce-Type":        c.Request().Header.Get("Ce-Type"),
			"Ce-Source":      c.Request().Header.Get("Ce-Source"),
			"Content-Type":   *ContentType,
		}

		for key, val := range headers {
			req.Header.Add(key, val)
		}

		if _, err := http.DefaultClient.Do(req); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusAccepted, "seed sent successfully")
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))

}

func getParsedRequestBody(req *seedRequest) io.Reader {
	reqData, _ := json.Marshal(req)
	return bytes.NewBuffer(reqData)
}
