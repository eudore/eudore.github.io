package config;

import (
	"fmt"
	"os"
	"testing"
	"public/config"
	"public/log"
)

func TestConfig(t *testing.T) {
	dir := config.Getconst("file_oss_upload")
	*dir = "new file_oss_upload"
	fmt.Println(*config.Getconst("file_oss_upload"))

	os.Setenv("ENV_LISTEN_HTTPS","")
	os.Setenv("ENV_LISTEN_Certfile","/data/")
	config.Reload()
	log.Json(config.Instance())
}