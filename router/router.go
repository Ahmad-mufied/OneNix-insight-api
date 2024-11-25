package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google-custom-search/handler"
)

func RegisterRoutes(e *echo.Echo, newsHandler *handler.NewsHandler) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/news", newsHandler.GetLatestNews)
}
