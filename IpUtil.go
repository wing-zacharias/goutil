package util

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	RegexIp = "(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})(\\.(2(5[0-5]{1}|[0-4]\\d{1})|[0-1]?\\d{1,2})){3}"
)

// GetHostname
/**
* @author: wing
* @time: 2020/9/4 9:52
* @param:
* @return:
* @comment: get local hostname
**/
func GetHostname() string {
	hostName, _ := os.Hostname()
	return hostName
}

// Ipv4Check
/**
* @author: wing
* @time: 2020/9/4 12:29
* @param:
* @return:
* @comment: ipv4 check
**/
func Ipv4Check(ipv4Str string) bool {
	ip := net.ParseIP(ipv4Str)
	if ip != nil {
		return true
	}
	return false
}

// ParseCidr
/**
* @author: wing
* @time: 2020/9/4 12:31
* @param:
* @return: ip,length of mask
* @comment: cidr address disassemble
**/
func ParseCidr(cidrAddress string) (string, int, error) {
	cidr, i, err := net.ParseCIDR(cidrAddress)
	if err != nil {
		return "", 0, err
	}
	mask, _ := i.Mask.Size()
	return cidr.String(), mask, nil

}

// GetNetIpList
/**
* @author: wing
* @time: 2020/9/4 9:56
* @param:
* @return:
* @comment: get all ips in special net
**/
func GetNetIpList(ip string, mask string) ([]net.IP, net.IP, net.IP) {
	var res []net.IP
	iMask := IpStringToUint(mask)
	iIp := IpStringToUint(ip)
	start := iIp & iMask
	for i := uint32(1); i < ^iMask; i++ {
		res = append(res, net.ParseIP(IpUintToString(start+i)))
	}
	return res, net.ParseIP(IpUintToString(start)), net.ParseIP(IpUintToString(iIp | ^iMask))
}

// IpUintToString
/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment:
**/
func IpUintToString(ip uint32) string {
	p0 := uint8(ip >> 24)
	p1 := uint8(ip >> 16)
	p2 := uint8(ip >> 8)
	p3 := uint8(ip & uint32(255))
	return fmt.Sprintf("%d.%d.%d.%d", p0, p1, p2, p3)
}

// IpStringToUint
/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment:
**/
func IpStringToUint(ip string) uint32 {
	if net.ParseIP(ip) != nil {
		pip := strings.Split(ip, ".")
		if len(pip) == 4 {
			p0, _ := strconv.Atoi(pip[0])
			p1, _ := strconv.Atoi(pip[1])
			p2, _ := strconv.Atoi(pip[2])
			p3, _ := strconv.Atoi(pip[3])
			return uint32(p0)<<24 + uint32(p1)<<16 + uint32(p2)<<8 + uint32(p3)
		}
	}
	return 0
}

// MaskCidrToString
/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment:
**/
func MaskCidrToString(mask int) string {
	ip := uint32(math.Pow(2, 32) - math.Pow(2, float64(32-mask)))
	p0 := uint8(ip >> 24)
	p1 := uint8(ip >> 16)
	p2 := uint8(ip >> 8)
	p3 := uint8(ip & uint32(255))
	return fmt.Sprintf("%d.%d.%d.%d", p0, p1, p2, p3)
}

// MaskStringToCidr
/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment:
**/
func MaskStringToCidr(mask string) int {
	if net.ParseIP(mask) != nil {
		pip := strings.Split(mask, ".")
		if len(pip) == 4 {
			p0, _ := strconv.Atoi(pip[0])
			p1, _ := strconv.Atoi(pip[1])
			p2, _ := strconv.Atoi(pip[2])
			p3, _ := strconv.Atoi(pip[3])
			return 32 - int(math.Log2(float64(^(uint32(p0)<<24+uint32(p1)<<16+uint32(p2)<<8+uint32(p3))+1)))
		}
	}
	return 0
}

// MaskCidrToUint
/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment:
**/
func MaskCidrToUint(mask int) uint32 {
	return uint32(math.Pow(2, 32) - math.Pow(2, float64(32-mask)))
}

// MaskUintToCidr
/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment:
**/
func MaskUintToCidr(mask uint32) int {
	return int(32 - math.Log2(float64(^mask+1)))
}
