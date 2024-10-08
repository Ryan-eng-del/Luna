package options

import (
	"encoding/json"
	"hurricane/internal/pkg/options"
	"hurricane/pkg/app"
	"hurricane/pkg/log"
)

type Options struct {
	GenericServerRunOptions *options.ServerRunOptions       `json:"server"   mapstructure:"server"`
	GRPCOptions             *options.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	InsecureServing         *options.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *options.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	MySQLOptions            *options.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	RedisOptions            *options.RedisOptions           `json:"redis"    mapstructure:"redis"`
	EtcdOptions             *options.EtcdOptions            `json:"etcd"     mapstructure:"etcd"`
	Log                     *log.Options                    `json:"log"      mapstructure:"log"`
	FeatureOptions          *options.FeatureOptions         `json:"feature"  mapstructure:"feature"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		GRPCOptions:             options.NewGRPCOptions(),
		InsecureServing:         options.NewInsecureServingOptions(),
		SecureServing:           options.NewSecureServingOptions(),
		MySQLOptions:            options.NewMySQLOptions(),
		RedisOptions:            options.NewRedisOptions(),
		Log:                     log.NewOption(),
		FeatureOptions:          options.NewFeatureOptions(),
	}
	return &o
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss app.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
