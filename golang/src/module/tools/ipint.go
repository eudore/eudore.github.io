package tools;

import (
    "fmt"
    "strings"
    "strconv"
)


func Ipstr(ipnr int64) string {
    var bytes [4]int
    bytes[0] = int(ipnr & 0xFF)
    bytes[1] = int((ipnr >> 8) & 0xFF)
    bytes[2] = int((ipnr >> 16) & 0xFF)
    bytes[3] = int((ipnr >> 24) & 0xFF)
    return fmt.Sprintf("%d.%d.%d.%d",bytes[3],bytes[2],bytes[1],bytes[0])
}
 
func Ipint(ipnr string) int64 {     
    bits := strings.Split(ipnr, ".")
     
    b0, _ := strconv.Atoi(bits[0])
    b1, _ := strconv.Atoi(bits[1])
    b2, _ := strconv.Atoi(bits[2])
    b3, _ := strconv.Atoi(bits[3])
 
    var sum int64
     
    sum += int64(b0) << 24
    sum += int64(b1) << 16
    sum += int64(b2) << 8
    sum += int64(b3)
     
    return sum
}
