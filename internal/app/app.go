package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	userpb "github.com/leondevpt/resume-backend/apigen/go/user/v1"
	"github.com/leondevpt/resume-backend/internal/config"
	"github.com/leondevpt/resume-backend/internal/server"
	"github.com/leondevpt/resume-backend/pkg/logging"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
)

type App struct {
	Cfg     *config.Config
	GRPC    *server.GrpcServer
	ErrChan chan error
}

func NewApp(c *config.Config, s *server.GrpcServer) *App {
	return &App{Cfg: c, GRPC: s, ErrChan: make(chan error, 1)}
}

func (a *App) Start(ctx context.Context) {
	logger := logging.FromContext(ctx)
	port := a.Cfg.App.GrpcPort
	// creating a listener for server
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Infow("app start", "Listen on port", a.Cfg.App.GrpcPort)
	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux()

	// setting up a dail up for gRPC service by specifying endpoint/target url
	err = userpb.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%d", port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		logger.Fatalln(err)
	}

	m := cmux.New(l)

	// Create a grpc listener first
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

	// All the rest is assumed to be HTTP
	httpListener := m.Match(cmux.Any())
	httpServer := &http.Server{}

	// start server
	go func() {
		a.ErrChan <- httpServer.Serve(httpListener)
	}()
	go func() {
		a.ErrChan <- a.GRPC.Server.Serve(grpcListener)
	}()

	go func() {
		a.ErrChan <- m.Serve()
	}()
}

func (a *App) Stop() {
	a.GRPC.Stop()
}

func Run(c context.Context, app *App) error {
	app.Start(c)
	defer app.Stop()
	select {
	case grpcErr := <-app.ErrChan:
		return grpcErr
	case <-c.Done():
		return nil
	}
}
