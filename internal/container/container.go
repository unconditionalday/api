package container

import (
	"net/http"

	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/client/github"
	"github.com/unconditionalday/server/internal/parser"
	"github.com/unconditionalday/server/internal/repository/bleve"
	"github.com/unconditionalday/server/internal/version"
	"github.com/unconditionalday/server/internal/webserver"
	blevex "github.com/unconditionalday/server/internal/x/bleve"
	calverx "github.com/unconditionalday/server/internal/x/calver"
	netx "github.com/unconditionalday/server/internal/x/net"
	"go.uber.org/zap"
)

func NewDefaultParameters() Parameters {
	return Parameters{
		ServerAddress:        "0.0.0.0",
		ServerPort:           8080,
		ServerAllowedOrigins: []string{"*"},
		SourceRepository:     "source",
		SourceClientKey:      "secret",
		FeedIndex:            "feed.index",
	}
}

func NewParameters(serverAddress, feedIndex, sourceRepository,sourceClientKey string, serverPort int, serverAllowedOrigins []string) Parameters {
	return Parameters{
		ServerAddress:        serverAddress,
		ServerPort:           serverPort,
		ServerAllowedOrigins: serverAllowedOrigins,

		SourceRepository: sourceRepository,
		SourceClientKey: sourceClientKey,

		FeedIndex: feedIndex,
	}
}

type Parameters struct {
	ServerAddress        string
	ServerPort           int
	ServerAllowedOrigins []string

	SourceRepository string
	SourceClientKey string

	FeedIndex string
}

type Services struct {
	apiServer      *webserver.Server
	feedRepository *bleve.FeedRepository
	sourceClient   *github.Client
	httpClient     *netx.HttpClient
	logger         *zap.Logger
	parser         *parser.Parser
	versioning     *calverx.CalVer
}

func NewContainer(p Parameters) (*Container, error) {
	return &Container{
		Parameters: p,
	}, nil
}

func (c *Container) GetAPIServer() *webserver.Server {
	if c.apiServer != nil {
		return c.apiServer
	}

	config := webserver.Config{
		Address:        c.Parameters.ServerAddress,
		Port:           c.Parameters.ServerPort,
		AllowedOrigins: c.Parameters.ServerAllowedOrigins,
	}

	c.apiServer = webserver.NewServer(config, c.GetFeedRepository(), c.GetLogger())

	return c.apiServer
}

func (c *Container) GetFeedRepository() app.FeedRepository {
	if c.feedRepository != nil {
		return c.feedRepository
	}

	b, err := blevex.New(c.FeedIndex)
	if err != nil {
		panic(err)
	}

	c.feedRepository = bleve.NewFeedRepository(b)

	return c.feedRepository
}

func (c *Container) GetSourceClient() app.SourceClient {
	if c.sourceClient != nil {
		return c.sourceClient
	}

	c.sourceClient = github.New(c.SourceRepository, "unconditionalday", c.SourceClientKey, http.DefaultClient)

	return c.sourceClient
}

func (c *Container) GetVersioning() version.Versioning {
	if c.versioning != nil {
		return c.versioning
	}

	c.versioning = calverx.New()

	return c.versioning
}

func (c *Container) GetHTTPClient() netx.Client {
	if c.httpClient != nil {
		return c.httpClient
	}

	c.httpClient = netx.NewHttpClient()

	return c.httpClient
}

// TODO: Needs to export a Logger interface
func (c *Container) GetLogger() *zap.Logger {
	if c.logger != nil {
		return c.logger
	}

	var err error
	c.logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return c.logger
}

func (c *Container) GetParser() *parser.Parser {
	if c.parser != nil {
		return c.parser
	}

	c.parser = parser.New()

	return c.parser
}

type Container struct {
	Parameters
	Services
}
