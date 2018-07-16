package store

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

func getsize(file_bytes int64) string {
	var i     int
	var units = [6]string{"B", "K", "M", "G", "T", "P"}
	i = 0
	for {
		if file_bytes < 1024 {
			return fmt.Sprintf("%d", file_bytes) + units[i]
		}
		file_bytes = file_bytes >> 10
		i++
	}
}

func IsFile(path string) bool {
	var n int = 0
	if stmtQueryIsFile.QueryRow(PathHash(path)).Scan(&n) != nil {
		return false
	}
	return n!=0
}

func PathHash(path string) string{
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(path))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}