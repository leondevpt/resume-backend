package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	userpb "github.com/leondevpt/resume-backend/apigen/go/user/v1"
	"github.com/leondevpt/resume-backend/internal/config"
	"github.com/leondevpt/resume-backend/pkg/logging"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"math"
	"time"
)

// GrpcServer implements a gRPC Server for the user service
type GrpcServer struct {
	Cfg    *config.Config
	Server *grpc.Server
	health *health.Server
}

func NewGrpcServer(c context.Context, cfg *config.Config, userService userpb.UserServiceServer) *GrpcServer {
	// 设置一元拦截器
	var interceptors []grpc.UnaryServerInterceptor

	interceptors = append(interceptors,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_recovery.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logging.FromContext(c).Desugar()),
	)
	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.MaxRecvMsgSize(math.MaxInt64),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				/*
					When the connection reaches its max-age, it will be closed and will trigger a re-resolve from the client.
					If new instances were added in the meantime, the client will see them now and open connections to them as well.
				*/
				MaxConnectionAge: time.Minute * 5,
			},
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptors...,
		)),
	}
	// 初始化grpc对象
	server := grpc.NewServer(opts...)

	// 注册服务
	userpb.RegisterUserServiceServer(server, userService)
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(server, healthcheck)
	reflection.Register(server)
	grpc_prometheus.Register(server)

	return &GrpcServer{Cfg: cfg, Server: server, health: healthcheck}
}

// Stop stops the gRPC server
func (g *GrpcServer) Stop() {
	g.health.Shutdown()
	g.Server.GracefulStop()
	return
}
