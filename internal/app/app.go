package app

import (
	"github.com/eachinchung/component-base/app"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/app/config"
	"github.com/eachinchung/e-service/internal/app/options"
)

const commandDesc = `Eachin Service 包含了 Eachin 提供的所有云服务。`

func NewApplication(basename string) *app.Application {
	opts := options.NewOptions()
	application := app.NewApplication("Eachin Cloud Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.LogOptions)
		defer log.Flush()

		cfg := config.GetConfigIns(opts)
		return Run(cfg)
	}
}
