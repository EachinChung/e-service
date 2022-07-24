package options

//goland:noinspection SpellCheckingInspection
import (
	"encoding/json"

	"github.com/eachinchung/component-base/cli/flag"
	baseoptions "github.com/eachinchung/component-base/options"
	"github.com/eachinchung/log"

	"github.com/eachinchung/e-service/internal/pkg/options"
)

type Options struct {
	GenericServerRunOptions *options.ServerRunOptions    `json:"server"      mapstructure:"server"`
	PostgresOptions         *baseoptions.PostgresOptions `json:"postgres"    mapstructure:"postgres"`
	RedisOptions            *baseoptions.RedisOptions    `json:"redis"       mapstructure:"redis"`
	JWTOptions              *baseoptions.JWTOptions      `json:"jwt"         mapstructure:"jwt"`
	CasbinOptions           *baseoptions.CasbinOptions   `json:"casbin"      mapstructure:"casbin"`
	LogOptions              *log.Options                 `json:"log"         mapstructure:"log"`
}

func (o Options) Flags() (fss flag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.PostgresOptions.AddFlags(fss.FlagSet("postgres"))
	o.RedisOptions.AddFlags(fss.FlagSet("rides"))
	o.JWTOptions.AddFlags(fss.FlagSet("jwt"))
	o.CasbinOptions.AddFlags(fss.FlagSet("casbin"))
	o.LogOptions.AddFlags(fss.FlagSet("logs"))

	return fss
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.PostgresOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.JWTOptions.Validate()...)
	errs = append(errs, o.CasbinOptions.Validate()...)
	errs = append(errs, o.LogOptions.Validate()...)

	return errs
}

func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		PostgresOptions:         baseoptions.NewPostgresOptions(),
		RedisOptions:            baseoptions.NewRedisOptions(),
		JWTOptions:              baseoptions.NewJWTOptions(),
		CasbinOptions:           baseoptions.NewCasbinOptions(),
		LogOptions:              log.NewOptions(),
	}
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
