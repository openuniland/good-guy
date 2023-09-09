package configs

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type MongoDB struct {
	MongoDBUsername   string `mapstructure:"MONGODB_USERNAME" json:"MONGODB_USERNAME"`
	MongoDBPassword   string `mapstructure:"MONGODB_PASSWORD" json:"MONGODB_PASSWORD"`
	MongoDBHost       string `mapstructure:"MONGODB_HOST" json:"MONGODB_HOST"`
	MongoDBReplicaSet string `mapstructure:"MONGODB_REPLICA_SET" json:"MONGODB_REPLICA_SET"`
	MongoDBName       string `mapstructure:"MONGODB_NAME" json:"MONGODB_NAME"`
	MongoDBProtocol   string `mapstructure:"MONGODB_PROTOCOL" json:"MONGODB_PROTOCOL"`
}

type ServerConfig struct {
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	Env               string `mapstructure:"ENV"`
	AuthKey           string `mapstructure:"AUTH_KEY"`
	Host              string `mapstructure:"HOST"`
	Tz                string `mapstructure:"TZ"`
}

type UrlCrawlerList struct {
	FithouUrl           string `mapstructure:"FITHOU_URL"`
	CtmsUrl             string `mapstructure:"CTMS_URL"`
	FithouCategoriesUrl string `mapstructure:"FITHOU_CATEGORIES_URL"`
	ExamScheduleUrl     string `mapstructure:"EXAM_SCHEDULE_URL"`
}

type FBConfig struct {
	FBVerifyToken string `mapstructure:"FB_VERIFY_TOKEN"`
	AppCode       string `mapstructure:"APP_CODE"`
}

type Jobs struct {
	SyncArticlesFromFithou string `mapstructure:"SYNC_ARTICLES_FROM_FITHOU"`
}

type MappingConfigs struct {
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	Env               string `mapstructure:"ENV"`
	AuthKey           string `mapstructure:"AUTH_KEY"`
	Host              string `mapstructure:"HOST"`
	Tz                string `mapstructure:"TZ"`

	MongoDBUsername   string `mapstructure:"MONGODB_USERNAME"`
	MongoDBPassword   string `mapstructure:"MONGODB_PASSWORD"`
	MongoDBHost       string `mapstructure:"MONGODB_HOST"`
	MongoDBReplicaSet string `mapstructure:"MONGODB_REPLICA_SET"`
	MongoDBName       string `mapstructure:"MONGODB_NAME"`
	MongoDBProtocol   string `mapstructure:"MONGODB_PROTOCOL"`

	FithouUrl           string `mapstructure:"FITHOU_URL"`
	CtmsUrl             string `mapstructure:"CTMS_URL"`
	FithouCategoriesUrl string `mapstructure:"FITHOU_CATEGORIES_URL"`
	ExamScheduleUrl     string `mapstructure:"EXAM_SCHEDULE_URL"`

	FBVerifyToken string `mapstructure:"FB_VERIFY_TOKEN"`
	AppCode       string `mapstructure:"APP_CODE"`

	SyncArticlesFromFithou string `mapstructure:"SYNC_ARTICLES_FROM_FITHOU"`
}

type Configs struct {
	Server         ServerConfig   `json:"server"`
	MongoDB        MongoDB        `json:"mongodb"`
	UrlCrawlerList UrlCrawlerList `json:"url_crawler_list"`
	FBConfig       FBConfig       `json:"fb_config"`
	Jobs           Jobs           `json:"jobs"`
}

func LoadConfigs(path string) (configs *Configs, err error) {
	if os.Getenv("ENV") == ".env" {
		configs = &Configs{
			Server: ServerConfig{
				HttpServerAddress: os.Getenv("HTTP_SERVER_ADDRESS"),
				Env:               os.Getenv("ENV"),
				AuthKey:           os.Getenv("AUTH_KEY"),
				Host:              os.Getenv("HOST"),
				Tz:                os.Getenv("TZ"),
			},
			MongoDB: MongoDB{
				MongoDBUsername:   os.Getenv("MONGODB_USERNAME"),
				MongoDBPassword:   os.Getenv("MONGODB_PASSWORD"),
				MongoDBHost:       os.Getenv("MONGODB_HOST"),
				MongoDBReplicaSet: os.Getenv("MONGODB_REPLICA_SET"),
				MongoDBName:       os.Getenv("MONGODB_NAME"),
				MongoDBProtocol:   os.Getenv("MONGODB_PROTOCOL"),
			},
			UrlCrawlerList: UrlCrawlerList{
				FithouUrl:           os.Getenv("FITHOU_URL"),
				CtmsUrl:             os.Getenv("CTMS_URL"),
				FithouCategoriesUrl: os.Getenv("FITHOU_CATEGORIES_URL"),
				ExamScheduleUrl:     os.Getenv("EXAM_SCHEDULE_URL"),
			},
			FBConfig: FBConfig{
				FBVerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
				AppCode:       os.Getenv("APP_CODE"),
			},
			Jobs: Jobs{
				SyncArticlesFromFithou: os.Getenv("SYNC_ARTICLES_FROM_FITHOU"),
			},
		}

		return configs, nil
	}

	var mapping *MappingConfigs
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msgf("[ERROR]:[LoadConfigs]:[viper.ReadInConfig()]:[ERROR_INFO=%v]", err)
		return configs, err
	}

	err = viper.Unmarshal(&mapping)
	if err != nil {
		log.Fatal().Err(err).Msgf("[ERROR]:[LoadConfigs]:[viper.Unmarshal()]:[ERROR_INFO=%v]", err)
		return configs, err
	}

	configs = &Configs{
		Server: ServerConfig{
			HttpServerAddress: mapping.HttpServerAddress,
			Env:               mapping.Env,
			AuthKey:           mapping.AuthKey,
			Host:              mapping.Host,
			Tz:                mapping.Tz,
		},
		MongoDB: MongoDB{
			MongoDBUsername:   mapping.MongoDBUsername,
			MongoDBPassword:   mapping.MongoDBPassword,
			MongoDBHost:       mapping.MongoDBHost,
			MongoDBReplicaSet: mapping.MongoDBReplicaSet,
			MongoDBName:       mapping.MongoDBName,
			MongoDBProtocol:   mapping.MongoDBProtocol,
		},
		UrlCrawlerList: UrlCrawlerList{
			FithouUrl:           mapping.FithouUrl,
			CtmsUrl:             mapping.CtmsUrl,
			FithouCategoriesUrl: mapping.FithouCategoriesUrl,
			ExamScheduleUrl:     mapping.ExamScheduleUrl,
		},
		FBConfig: FBConfig{
			FBVerifyToken: mapping.FBVerifyToken,
			AppCode:       mapping.AppCode,
		},
		Jobs: Jobs{
			SyncArticlesFromFithou: mapping.SyncArticlesFromFithou,
		},
	}

	return configs, nil
}
