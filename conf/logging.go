package conf

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LoggingConfig struct {
	Level            string                 `mapstructure:"log_level" json:"log_level"`
	File             string                 `mapstructure:"log_file" json:"log_file"`
	DisableColors    bool                   `mapstructure:"disable_colors" split_words:"true" json:"disable_colors"`
	QuoteEmptyFields bool                   `mapstructure:"quote_empty_fields" split_words:"true" json:"quote_empty_fields"`
	TSFormat         string                 `mapstructure:"ts_format" json:"ts_format"`
	Fields           map[string]interface{} `mapstructure:"fields" json:"fields"`
	UseNewLogger     bool                   `mapstructure:"use_new_logger",split_words:"true"`

	BugSnag *BugSnagConfig
}

func ConfigureLogging(config *LoggingConfig) (*logrus.Entry, error) {
	logger := logrus.New()

	tsFormat := time.RFC3339Nano
	if config.TSFormat != "" {
		tsFormat = config.TSFormat
	}
	// always use the full timestamp
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  tsFormat,
		DisableColors:    config.DisableColors,
		QuoteEmptyFields: config.QuoteEmptyFields,
	})

	// use a file if you want
	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
		if errOpen != nil {
			return nil, errOpen
		}
		logger.SetOutput(f)
		logger.Infof("Set output file to %s", config.File)
	}

	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(level)
		logger.Debug("Set log level to: " + logger.GetLevel().String())
	}

	if err := AddBugSnagHook(config.BugSnag); err != nil {
		return nil, errors.Wrap(err, "Failed to configure bugsnag")
	}

	f := logrus.Fields{}
	for k, v := range config.Fields {
		f[k] = v
	}

	return logger.WithFields(f), nil
}