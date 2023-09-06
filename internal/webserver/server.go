package webserver

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	api "github.com/unconditionalday/server/api"
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/version"
)

type Server struct {
	config       Config
	feedRepo     app.FeedRepository
	source       *app.SourceRelease
	buildVersion version.Build
	logger       *zap.Logger
	client       *echo.Echo
}

type Config struct {
	Port           int
	Address        string
	AllowedOrigins []string
}

func NewServer(config Config, repo app.FeedRepository, source *app.SourceRelease, version version.Build, logger *zap.Logger) *Server {
	return &Server{
		config:       config,
		feedRepo:     repo,
		source:       source,
		buildVersion: version,
		logger:       logger,
		client:       echo.New(),
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
	feeds, err := s.feedRepo.Find(query)
	if err != nil {
		e := api.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}

		return ctx.JSON(http.StatusInternalServerError, e)
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

	return ctx.JSON(http.StatusOK, fi)
}

func (s *Server) GetV1Version(ctx echo.Context) error {
	v := api.ServerVersion{
		Build: api.ServerBuildVersion{
			Version: s.buildVersion.Version,
			Commit:  s.buildVersion.Commit,
		},
		Source: api.SourceReleaseVersion{
			Version:       s.source.Version,
			LastUpdatedAt: s.source.LastUpdateAt,
		},
	}

	return ctx.JSON(http.StatusOK, v)
}

func (s *Server) GetV1InformerWikiQuery(ctx echo.Context, query string) error {
	// TODO: add language support
	wikiRes, err := s.wiki.Search(query, "en")
	if err != nil {
		e := api.Error{
			Code:    500,
			Message: "Internal Server Error",
		}

		s.logger.Error("wiki search", zap.Error(err))

		return ctx.JSON(500, e)
	}

	res := api.WikiResult{
		Language:  wikiRes.Language,
		Link:      wikiRes.Link,
		Summary:   wikiRes.Summary,
		Thumbnail: wikiRes.Thumbnail,
		Title:     wikiRes.Title,
	}

	return ctx.JSON(200, res)
}
