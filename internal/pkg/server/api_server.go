package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/eachinchung/component-base/core"
	"github.com/eachinchung/component-base/middleware"
	"github.com/eachinchung/errors"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/pkg/code"
)

// APIServer 常用的api服务器
type APIServer struct {
	middlewares    []string
	trustedProxies []string
	mode           string

	ServingInfo *ServingInfo

	// ShutdownTimeout 是用于服务器关闭的超时时间。
	// 这指定了服务器正常关闭返回之前的超时时间。
	ShutdownTimeout time.Duration

	*gin.Engine
	healthz bool

	server *http.Server
}

// InstallAPIs 安装默认 APIs
func (s *APIServer) InstallAPIs() {
	// 健康检查
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}
}

// Setup 一些关于 gin 的安装工作
func (s *APIServer) Setup() {
	gin.SetMode(s.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	trustedProxies := s.trustedProxies
	if len(trustedProxies) == 0 {
		trustedProxies = nil
	}

	err := s.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatalf("设置可信代理错误, err: %+v", err)
		return
	}
}

// InstallMiddlewares 安装中间件
func (s *APIServer) InstallMiddlewares() {
	// 必要的中间件
	s.Use(middleware.RecoveryWithHandle(func(c *gin.Context, err interface{}) {
		core.WriteResponse(c, errors.Code(code.ErrUnknown, "服务器内部错误"), nil)
		c.Abort()
	}))
	s.Use(middleware.RequestID())

	// 安装自定义中间件
	for _, m := range s.middlewares {
		mw, ok := middleware.Store[m]
		if !ok {
			log.Warnf("找不到中间件: %s", m)
			continue
		}

		log.Infof("安装中间件: %s", m)
		s.Use(mw)
	}
}

func initGenericAPIServer(s *APIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// Run 生成http服务器。仅当最初无法监听端口时，它才返回错误。
func (s *APIServer) Run() error {
	s.server = &http.Server{
		Addr:    s.ServingInfo.Address(),
		Handler: s,
	}

	var eg errgroup.Group

	// 在 goroutine 中初始化服务器，这样它就不会阻止下面的优雅关机处理
	eg.Go(func() error {
		log.Infof("开始在监听传入请求 (http: %s)", s.ServingInfo.Address())

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
			return err
		}

		log.Infof("服务已经被停止 (http: %s)", s.ServingInfo.Address())
		return nil
	})

	// Ping 服务器以确保服务正常工作。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

// Close 优雅关停 api server.
func (s *APIServer) Close() {
	// 上下文用于通知服务器它有 10 秒的时间来完成它当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Warnf("关闭服务器失败: %s", err.Error())
	}
}

// ping http服务器以确保路由器正常工作。
func (s *APIServer) ping(ctx context.Context) error {
	//goland:noinspection HttpUrlsUsage
	url := fmt.Sprintf("http://%s/healthz", s.ServingInfo.Address())
	if strings.Contains(s.ServingInfo.Address(), "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%d/healthz", s.ServingInfo.BindPort)
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		// 通过向 `/healthz` 发送 GET 请求来 Ping 服务器。
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("服务已部署成功。")

			//goland:noinspection GoUnhandledErrorResult
			resp.Body.Close()

			return nil
		}

		// 休眠一秒钟以继续下一次 ping。
		log.Info("等待服务器，1秒后重试。")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("无法在指定的时间间隔内 ping 通服务器。")
		default:
		}
	}
}
