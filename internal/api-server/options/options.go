// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"

	"github.com/lunarianss/Luna/infrastructure/app"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
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
	EmailOptions            *options.EmailOptions           `json:"email" mapstructure:"email"`
	JwtOptions              *options.JwtOptions             `json:"jwt" mapstructure:"jwt"`
	SystemOptions           *options.SystemOptions          `json:"system" mapstructure:"system"`
	MQOptions               *options.RocketMQOptions        `json:"rocket-mq" mapstructure:"rocket-mq"`
	WeaviateOptions         *options.WeaviateOptions        `json:"weaviate" mapstructure:"weaviate"`
	MinioOptions            *options.MinioOptions           `json:"minio" mapstructure:"minio"`
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
		EmailOptions:            options.NewEmailOptions(),
		JwtOptions:              options.NewJwtOptions(),
		SystemOptions:           options.NewSystemOptions(),
		MQOptions:               options.NewRocketMQOptions(),
		WeaviateOptions:         options.NewWeaviateOptions(),
		MinioOptions:            options.NewMinioOptions(),
	}

	return &o
}

// Flags returns flags for a specific api-server by section name.
func (o *Options) Flags() (fss app.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.EmailOptions.AddFlags(fss.FlagSet("email"))
	o.JwtOptions.AddFlags(fss.FlagSet("jwt"))

	return fss
}

func (o *Options) String() string {
	//nolint: errchkjson
	data, _ := json.Marshal(o)

	return string(data)
}
