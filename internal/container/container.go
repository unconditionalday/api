package container

import (
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/parser"
	"github.com/unconditionalday/server/internal/repository/bleve"
	"github.com/unconditionalday/server/internal/service"
	"github.com/unconditionalday/server/internal/webserver"
	netx "github.com/unconditionalday/server/internal/x/net"
	"go.uber.org/zap"
)

func NewDefaultParameters() Parameters {
	return Parameters{
		ServerAddress:        "0.0.0.0",
		ServerPort:           8080,
		ServerAllowedOrigins: []string{"*"},

		Index: "test.index",
	}
}

func NewParameters(serverAddress, index string, serverPort int, serverAllowedOrigins []string) Parameters {
	return Parameters{
		ServerAddress:        serverAddress,
		ServerPort:           serverPort,
		ServerAllowedOrigins: serverAllowedOrigins,
		Index:                index,
	}
}

type Parameters struct {
	ServerAddress        string
	ServerPort           int
	ServerAllowedOrigins []string

	Index string
}

type Services struct {
	apiServer      *webserver.Server
	feedRepository *bleve.Bleve
	sourceService  *service.Source
	httpClient     *netx.HttpClient
	logger         *zap.Logger
	parser         *parser.Parser
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

	b, err := bleve.NewBleve(c.Index)
	if err != nil {
		panic(err)
	}

	c.feedRepository = b

	return c.feedRepository
}

func (c *Container) GetSourceService() app.SourceService {
	if c.sourceService != nil {
		return c.sourceService
	}

	client := c.GetHTTPClient()

	c.sourceService = service.NewSource(client)

	return c.sourceService
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
	c.logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return c.logger
}

// TODO: Needs to export a Parser interface
func (c *Container) GetParser() *parser.Parser {
	if c.parser != nil {
		return c.parser
	}

	c.parser = parser.NewParser()

	return c.parser
}

type Container struct {
	Parameters
	Services
}
