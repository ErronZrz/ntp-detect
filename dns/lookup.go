package dns

import (
	"active/output"
	"active/parser"
	"active/tcp"
	"active/udpdetect"
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type extraWork func(string, string) error

func OutputDNS(src, dst string) error {
	return commonDNS(src, dst, nil)
}

func DetectAfterDNS(src, dst string) error {
	visited := make(map[string]struct{})
	detectWork := func(domain, ip string) error {
		netAddr := net24(ip)
		if _, ok := visited[netAddr]; ok {
			return nil
		}
		visited[netAddr] = struct{}{}
		return detect(domain, ip)
	}
	return commonDNS(src, dst, detectWork)
}

func AsyncDetectAfterDNS(src, dst string) error {
	visited := make(map[string]struct{})
	var mu sync.RWMutex
	asyncDetectWork := func(domain, ip string) error {
		netAddr := net24(ip)
		mu.RLock()
		_, ok := visited[netAddr]
		mu.RUnlock()
		if ok {
			return nil
		}
		mu.Lock()
		visited[netAddr] = struct{}{}
		mu.Unlock()

		return asyncDetect(domain, ip)
	}
	return commonDNS(src, dst, asyncDetectWork)
}

func TLSAfterDNS(src, dst string) error {
	return commonDNS(src, dst, checkTLS)
}

func commonDNS(src, dst string, work extraWork) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening srcFile %s: %v", src, err)
	}
	defer closeFunc(srcFile, src)

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating dstFile %s: %v", dst, err)
	}
	defer closeFunc(dstFile, dst)

	scanner := bufio.NewScanner(srcFile)
	writer := bufio.NewWriter(dstFile)

	for scanner.Scan() {
		domain := scanner.Text()
		if len(domain) == 0 {
			_ = writer.WriteByte('\n')
			continue
		}
		if domain[0] == '#' {
			_, err = writer.WriteString(domain + "\n")
			if err != nil {
				return fmt.Errorf("error writing comment %s: %v", domain, err)
			}
			continue
		}
		fmt.Println(domain)
		_, err = writer.WriteString(domain + "\n")
		if err != nil {
			return fmt.Errorf("error writing domain %s: %v", domain, err)
		}

		ips, err := net.LookupIP(domain)
		if err != nil {
			_, err = writer.WriteString(fmt.Sprintf("    %v\n\n", err))
			if err != nil {
				return fmt.Errorf("error writing error: %v", err)
			}
			continue
		}

		if len(ips) == 0 {
			_, err = writer.WriteString("    no IP address found\n\n")
			if err != nil {
				return fmt.Errorf("error writing empty result: %v", err)
			}
			continue
		}

		for _, ip := range ips {
			ipStr := ip.String()
			if work != nil {
				err = work(domain, ipStr)
				if err != nil {
					return fmt.Errorf("error handling IP %s: %v", ipStr, err)
				}
			}
			_, err = writer.WriteString(fmt.Sprintf("    %s\n", ipStr))
			if err != nil {
				return fmt.Errorf("error writing IP %s: %v", ipStr, err)
			}
		}
		_ = writer.WriteByte('\n')
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}

func asyncDetect(domain, ip string) error {
	go func() {
		err := detect(domain, ip)
		if err != nil {
			fmt.Printf("error during detection: %v", err)
		}
	}()
	return nil
}

func detect(domain, ip string) error {
	cidr := ip + "/24"
	dataCh := udpdetect.DialNetworkNTP(cidr)
	if dataCh == nil {
		return errors.New("dataCh is nil")
	}

	seqNum := 0
	now := time.Now()
	for p, ok := <-dataCh; ok; p, ok = <-dataCh {
		err := p.Err
		if err != nil {
			return err
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			return err
		}
		seqNum++
		output.WriteToFile(p.Lines(), header.Lines(), domain+"_"+cidr, seqNum, p.RcvTime, now)
	}
	return nil
}

func checkTLS(domain, ip string) error {
	result := "x"
	if tcp.IsTLSEnabled(ip, 4460, "") {
		result = "Support"
	}
	fmt.Printf("%-30s%-20s%s\n", domain, ip, result)
	return nil
}

func closeFunc(f *os.File, path string) {
	err := f.Close()
	if err != nil {
		fmt.Printf("error closing file %s: %v", path, err)
	}
}

func net24(ip string) string {
	nums := strings.Split(ip, ".")
	return fmt.Sprintf("%s.%s.%s", nums[0], nums[1], nums[2])
}
