// Package bootstrap
package bootstrap

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
)

const (
	redisInitializeNil         = `redis cannot connect, please check your config or network`
	redisPingError             = `redis cannot connect, error: %v`
	logFieldHost               = "host"
	logFieldDB                 = "db"
	logFieldMaxRedirect        = "max_redirect"
	logFieldReadOnly           = "read_only"
	logFieldTlsEnable          = "tls_enable"
	logFieldInsecureSKipVerify = "insecure_skip_verify"
	logFieldRouteRandomly      = "route_randomly"
	logFieldRouteByLatency     = "route_by_latency"
)

// RegistryRedisNative initiate redis session
func RegistryRedisNative(conf *appctx.Config) redis.Cmdable {

	lf := []logger.Field{
		logger.Any(logFieldHost, conf.Redis.Hosts),
		logger.Any(logFieldDB, conf.Redis.DB),
		logger.Any(logFieldReadOnly, conf.Redis.ReadOnly),
		logger.Any(logFieldRouteByLatency, conf.Redis.RouteByLatency),
		logger.Any(logFieldRouteRandomly, conf.Redis.RouteRandomly),
		logger.Any(logFieldMaxRedirect, conf.Redis.MaxRedirect),
		logger.Any(logFieldTlsEnable, conf.Redis.TLSEnable),
		logger.Any(logFieldInsecureSKipVerify, conf.Redis.InsecureSkipVerify),
	}

	if conf.Redis.ClusterMode {
		return registryRedisCluster(conf, lf)
	}

	return registryRedisUniversal(conf, lf)
}

// registryRedisUniversal initiate redis session
func registryRedisUniversal(conf *appctx.Config, lf []logger.Field) redis.Cmdable {
	cfg := redis.UniversalOptions{
		Addrs: strings.Split(conf.Redis.Hosts, ","),
		DB:    conf.Redis.DB,

		ClientName: conf.Redis.ClientName,

		DialTimeout:  conf.Redis.DialTimeout,
		ReadTimeout:  conf.Redis.ReadTimeout,
		WriteTimeout: conf.Redis.WriteTimeout,

		PoolFIFO:        conf.Redis.PoolFIFO,
		PoolSize:        conf.Redis.PoolSize,
		PoolTimeout:     conf.Redis.PoolTimeout,
		MinIdleConns:    conf.Redis.MinIdleConn,
		MaxIdleConns:    conf.Redis.MaxIdleConn,
		MaxActiveConns:  conf.Redis.MaxActiveConn,
		ConnMaxIdleTime: conf.Redis.ConnIdleTime,
		ConnMaxLifetime: conf.Redis.ConnIdleTime,

		ReadOnly:       conf.Redis.ReadOnly,
		RouteByLatency: conf.Redis.RouteByLatency,
		MaxRedirects:   conf.Redis.MaxRedirect,
		RouteRandomly:  conf.Redis.RouteRandomly,

		Protocol: conf.Redis.Protocol,
		Username: conf.Redis.Username,
		Password: conf.Redis.Password,
		//SentinelUsername: conf.Redis.SentinelUsername,
		//SentinelPassword: conf.Redis.SentinelPassword,
		//
		//MasterName: conf.Redis.MasterName,
	}

	if conf.Redis.TLSEnable {
		cfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: conf.Redis.InsecureSkipVerify,
		}
	}

	r := redis.NewUniversalClient(&cfg)
	//r.AddHook(cachex.NewRedisHook(cfg.Addrs, 0))

	if r == nil {
		logger.Fatal(redisInitializeNil, lf...)
	}

	c := r.Ping(context.Background())

	if err := c.Err(); err != nil {
		logger.Fatal(fmt.Sprintf(redisPingError, err), lf...)
	}

	return r
}

// registryRedisCluster initiate redis session
func registryRedisCluster(conf *appctx.Config, lf []logger.Field) redis.Cmdable {

	cfg := &redis.ClusterOptions{
		Addrs: strings.Split(conf.Redis.Hosts, ","),

		ClientName: conf.Redis.ClientName,

		DialTimeout:  conf.Redis.DialTimeout,
		ReadTimeout:  conf.Redis.ReadTimeout,
		WriteTimeout: conf.Redis.WriteTimeout,

		PoolFIFO:        conf.Redis.PoolFIFO,
		PoolSize:        conf.Redis.PoolSize,
		PoolTimeout:     conf.Redis.PoolTimeout,
		MinIdleConns:    conf.Redis.MinIdleConn,
		MaxIdleConns:    conf.Redis.MaxIdleConn,
		MaxActiveConns:  conf.Redis.MaxActiveConn,
		ConnMaxIdleTime: conf.Redis.ConnIdleTime,
		ConnMaxLifetime: conf.Redis.ConnIdleTime,

		ReadOnly:       conf.Redis.ReadOnly,
		RouteByLatency: conf.Redis.RouteByLatency,
		MaxRedirects:   conf.Redis.MaxRedirect,
		RouteRandomly:  conf.Redis.RouteRandomly,

		Protocol: conf.Redis.Protocol,
		Username: conf.Redis.Username,
		Password: conf.Redis.Password,
	}

	if conf.Redis.TLSEnable {
		cfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: conf.Redis.InsecureSkipVerify,
		}
	}

	r := redis.NewClusterClient(cfg)
	//r.AddHook(cachex.NewRedisHook(cfg.Addrs, 0))

	if r == nil {
		logger.Fatal(redisInitializeNil, lf...)
	}

	c := r.Ping(context.Background())

	if err := c.Err(); err != nil {
		logger.Fatal(fmt.Sprintf(redisPingError, err), lf...)
	}

	return r
}

