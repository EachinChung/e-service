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
	GenericServerRunOptions *options.ServerRunOptions `json:"server"   mapstructure:"server"`
	MySQLOptions            *baseoptions.MySQLOptions `json:"mysql"    mapstructure:"mysql"`
	RedisOptions            *baseoptions.RedisOptions `json:"redis"    mapstructure:"redis"`
	JWTOptions              *baseoptions.JWTOptions   `json:"jwt"      mapstructure:"jwt"`
	LogOptions              *log.Options              `json:"log"      mapstructure:"log"`
}

func (o Options) Flags() (fss flag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("rides"))
	o.JWTOptions.AddFlags(fss.FlagSet("jwt"))
	o.LogOptions.AddFlags(fss.FlagSet("logs"))

	return fss
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.JWTOptions.Validate()...)
	errs = append(errs, o.LogOptions.Validate()...)

	return errs
}

func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		MySQLOptions:            baseoptions.NewMySQLOptions(),
		RedisOptions:            baseoptions.NewRedisOptions(),
		JWTOptions:              baseoptions.NewJwtOptions(),
		LogOptions:              log.NewOptions(),
	}
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
