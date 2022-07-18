package app

import (
	"github.com/eachinchung/component-base/shutdown"
	"github.com/eachinchung/component-base/shutdown/managers"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/config"
	"github.com/eachinchung/e-service/internal/app/store"
	"github.com/eachinchung/e-service/internal/app/store/mysql"
	"github.com/eachinchung/e-service/internal/pkg/server"
	"github.com/eachinchung/e-service/internal/pkg/validator"
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

	storeIns, err := mysql.GetMySQLFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatalf("获取 mysql 工厂失败, error: %v", err)
	}
	store.SetClient(storeIns)

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
		if mysqlStore, err := mysql.GetMySQLFactoryOr(nil); err == nil {
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
