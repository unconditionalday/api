package container

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/typesense/typesense-go/typesense"
	"github.com/unconditionalday/server/internal/app"
	"github.com/unconditionalday/server/internal/client/github"
	"github.com/unconditionalday/server/internal/client/wikipedia"
	"github.com/unconditionalday/server/internal/parser"
	typesenseRepo "github.com/unconditionalday/server/internal/repository/typesense"
	"github.com/unconditionalday/server/internal/search"
	"github.com/unconditionalday/server/internal/version"
	"github.com/unconditionalday/server/internal/webserver"
	calverx "github.com/unconditionalday/server/internal/x/calver"
	netx "github.com/unconditionalday/server/internal/x/net"
)

func NewDefaultParameters() Parameters {
	return Parameters{
		ServerAddress:        "0.0.0.0",
		ServerPort:           8080,
		ServerAllowedOrigins: []string{"*"},
		SourceRepository:     "source",
		SourceClientKey:      "secret",
		FeedRepositoryIndex:  "feed.index",
		LogEnv:               "dev",
	}
}

func NewParameters(serverAddress, feedIndex, feedRepoHost, feedRepoKey, sourceRepository, sourceClientKey, logEnv string, serverPort int, serverAllowedOrigins []string, buildVersion version.Build) Parameters {
	return Parameters{
		ServerAddress:        serverAddress,
		ServerPort:           serverPort,
		ServerAllowedOrigins: serverAllowedOrigins,

		SourceRepository: sourceRepository,
		SourceClientKey:  sourceClientKey,

		BuildVersion: buildVersion,

		FeedRepositoryIndex: feedIndex,
		FeedRepositoryHost:  feedRepoHost,
		FeedRepositoryKey:   feedRepoKey,

		LogEnv: logEnv,
	}
}

type Parameters struct {
	ServerAddress        string
	ServerPort           int
	ServerAllowedOrigins []string

	SourceRepository string
	SourceClientKey  string
	SourceRelease    *app.SourceRelease

	BuildVersion version.Build

	FeedRepositoryIndex string
	FeedRepositoryHost  string
	FeedRepositoryKey   string

	LogEnv string
}

type Services struct {
	apiServer       *webserver.Server
	feedRepository  *typesenseRepo.FeedRepository
	sourceClient    *github.Client
	searchClient    *wikipedia.Client
	httpClient      *netx.HttpClient
	typesenseClient *typesense.Client
	logger          *zap.Logger
	parser          *parser.Parser
	versioning      *calverx.CalVer
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

	c.apiServer = webserver.NewServer(config, c.GetFeedRepository(), c.SourceRelease, c.GetSearchClient(), c.BuildVersion, c.GetLogger())

	return c.apiServer
}

func (c *Container) GetFeedRepository() app.FeedRepository {
	if c.feedRepository != nil {
		return c.feedRepository
	}

	client := c.GetTypesenseClient()

	c.feedRepository = typesenseRepo.NewFeedRepository(client)

	return c.feedRepository
}

func (c *Container) GetTypesenseClient() *typesense.Client {
	if c.typesenseClient != nil {
		return c.typesenseClient
	}

	typesenseConnTimeout := 1 * time.Hour

	// TODO: Export it in a x/typesense pkg as wrapper with check config logic
	client := typesense.NewClient(
		typesense.WithServer(c.FeedRepositoryHost),
		typesense.WithAPIKey(c.FeedRepositoryKey),
		typesense.WithConnectionTimeout(typesenseConnTimeout),
		typesense.WithCircuitBreakerInterval(typesenseConnTimeout),
	)

	c.GetLogger().Info("Waiting typesense healthcheck...")
	if _, err := client.Health(context.Background(), typesenseConnTimeout); err != nil {
		panic(err)
	}

	c.typesenseClient = client

	return c.typesenseClient
}

func (c *Container) GetSourceClient() app.SourceClient {
	if c.sourceClient != nil {
		return c.sourceClient
	}

	c.sourceClient = github.New(c.SourceRepository, "unconditionalday", c.SourceClientKey, http.DefaultClient)

	return c.sourceClient
}

func (c *Container) GetSearchClient() search.SearchClient {
	if c.searchClient != nil {
		return c.searchClient
	}

	c.searchClient = wikipedia.NewClient()

	return c.searchClient
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

	if c.LogEnv == "prod" {
		c.logger, _ = zap.NewProduction()
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
