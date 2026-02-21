package logger

import "go.uber.org/zap/zapcore"

type Config struct {
	Level      string `yaml:"level" json:"level"`
	Format     string `yaml:"format" json:"format"`
	OutputFile string `yaml:"output_file" json:"output_file"`
}

func (c Config) Validate() []error {
	var errs []error
	_, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}
