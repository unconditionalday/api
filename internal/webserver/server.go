package webserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luigibarbato/isolated-think-source/api"
	"github.com/luigibarbato/isolated-think-source/internal/repository"
	"go.uber.org/zap"
)

type Server struct {
	config ServerConfig
	repo   repository.Repository
	client *echo.Echo
}

type ServerConfig struct {
	Port    int
	Address string
}

func NewServer(config ServerConfig, repo repository.Repository) *Server {
	return &Server{
		client: echo.New(),
		config: config,
		repo:   repo,
	}
}

func (s *Server) Start() error {
	api.RegisterHandlers(s.client, s)
	logger, _ := zap.NewProduction()
	s.client.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	s.client.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	return s.client.Start(fmt.Sprintf("%s:%d", s.config.Address, s.config.Port))
}

// (GET /v1/search/feed/{query})
func (s *Server) GetV1SearchFeedQuery(ctx echo.Context, query string) error {
	feeds, err := s.repo.Find(query)
	if err != nil {
		e := api.Error{
			Code:    500,
			Message: "Internal Server Error",
		}

		return ctx.JSON(500, e)
	}

	fi := make([]api.FeedItem, len(feeds))
	for i, f := range feeds {
		fi[i] = api.FeedItem{
			Source:   f.Source,
			Date:     f.Date,
			Language: f.Language,
			Link:     f.Link,
			Summary:  f.Summary,
			Title:    f.Title,
		}

		if f.Image != nil {
			fi[i].Image = &api.FeedImage{
				Title: f.Image.Title,
				Url:   f.Image.URL,
			}
		}
	}

	return ctx.JSON(200, fi)
}
