package logger

type Config struct {
	Level      string `yaml:"level" json:"level"`   // "trace", "debug", "info", "warn", "error"
	Format     string `yaml:"format" json:"format"` // "json", "text"
	OutputFile string `yaml:"output_file" json:"output_file"`
}

// func (l Config) Validate() error {
// 	if l.Level == "" {
// 		return errors.New("empty 'log' param: 'level'")
// 	}
// 	return nil
// }
