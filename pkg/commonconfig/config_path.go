package commonconfig

import "flag"

func GetConfigPath(defaultConfigPath string) string {
	path := flag.String("config", defaultConfigPath, "")
	flag.Parse()
	return *path
}
