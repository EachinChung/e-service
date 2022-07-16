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
	Log                     *log.Options              `json:"log"      mapstructure:"log"`
}

func (o Options) Flags() (fss flag.NamedFlagSets) {
	//o.GenericServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("rides"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	//o.Secret.AddFlags(fss.FlagSet("secret"))
	return fss
}

func (o *Options) Validate() []error {
	var errs []error

	//errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	//errs = append(errs, o.Secret.Validate()...)

	return errs
}

func NewOptions() *Options {
	return &Options{
		//GenericServerRunOptions: options.NewServerRunOptions(),
		MySQLOptions: baseoptions.NewMySQLOptions(),
		RedisOptions: baseoptions.NewRedisOptions(),
		Log:          log.NewOptions(),
		//Secret:                  options.NewSecretOptions(),
	}
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
