package app

import (
	"github.com/eachinchung/component-base/shutdown"
	"github.com/eachinchung/component-base/shutdown/managers"
	"github.com/eachinchung/e-service/internal/app/config"
	"github.com/eachinchung/e-service/internal/app/storage"
	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/casbin"
	"github.com/eachinchung/e-service/internal/app/store/postgres"
	"github.com/eachinchung/e-service/internal/app/validator"
	"github.com/eachinchung/e-service/internal/pkg/server"
	"github.com/eachinchung/log"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *server.APIServer
}

type preparedAPIServer struct {
	*apiServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(managers.NewPosixSignalManager())

	if err := validator.InitValidator(); err != nil {
		return nil, err
	}

	storeIns, err := postgres.GetPostgresFactoryOr(cfg.PostgresOptions)
	if err != nil {
		log.Fatalf("获取 postgres 工厂失败, error: %v", err)
	}
	store.SetClient(storeIns)

	_, err = storage.GetRedisClientOr(cfg.RedisOptions)
	if err != nil {
		log.Fatalf("获取 redis 客户端失败, error: %v", err)
	}

	if _, err := casbin.GetEnforcerOr(cfg.CasbinOptions); err != nil {
		log.Fatalf("获取 casbin 失败, error: %v", err)
	}

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	s := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
	}

	return s, nil
}

func buildGenericConfig(cfg *config.Config) (*server.Config, error) {
	genericConfig := server.NewConfig()
	if err := cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return nil, err
	}

	return genericConfig, nil
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	initRouter(s.genericAPIServer.Engine)

	s.gs.AddShutdownCallback(shutdown.Func(func(string) error {
		if mysqlStore, err := postgres.GetPostgresFactoryOr(nil); err == nil {
			_ = mysqlStore.Close()
		}

		s.genericAPIServer.Close()
		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// 启动优雅关机管理器
	if err := s.gs.Start(); err != nil {
		log.Fatalf("优雅关机管理器启动失败: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}
