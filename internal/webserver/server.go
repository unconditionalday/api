package webserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	api "github.com/unconditionalday/server/api"
	"github.com/unconditionalday/server/internal/app"
)

type Server struct {
	config Config
	repo   app.FeedRepository
	logger *zap.Logger
	client *echo.Echo
}

type Config struct {
	Port           int
	Address        string
	AllowedOrigins []string
}

func NewServer(config Config, repo app.FeedRepository, logger *zap.Logger) *Server {
	return &Server{
		client: echo.New(),
		config: config,
		repo:   repo,
		logger: logger,
	}
}

func (s *Server) Start() error {
	api.RegisterHandlers(s.client, s)
	s.client.Use(
		middleware.RequestLoggerWithConfig(
			middleware.RequestLoggerConfig{
				LogURI:    true,
				LogStatus: true,
				LogError:  true,
				LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
					s.logger.Info("request",
						zap.String("URI", v.URI),
						zap.Int("status", v.Status),
					)

					return nil
				},
			}))
	s.client.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.config.AllowedOrigins,
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
