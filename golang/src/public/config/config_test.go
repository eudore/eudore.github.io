package config;

import (
	"fmt"
	"testing"
	"public/config"
)

func TestConfig(t *testing.T) {
	fmt.Println(config.Instance())
	dir := config.Getconst("file_oss_upload")

	*dir = "ssssssss-------"
	fmt.Println(config.Instance())
	fmt.Println(*config.Getconst("file_oss_upload"))
}