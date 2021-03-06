package dns

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/net/common"
)

type DOH struct {
	DNS
	Server string
	Subnet *net.IPNet
	Proxy  func(domain string) (net.Conn, error)
	cache  *common.CacheExtend

	httpClient *http.Client
}

func NewDOH(host string) DNS {
	_, subnet, _ := net.ParseCIDR("0.0.0.0/0")
	dns := &DOH{
		Server: host,
		Subnet: subnet,
		Proxy: func(domain string) (net.Conn, error) {
			return net.DialTimeout("tcp", domain, 5*time.Second)
		},
		cache: common.NewCacheExtend(time.Minute * 20),
	}
	dns.SetProxy(dns.Proxy)
	return dns
}

// DOH DNS over HTTPS
// https://tools.ietf.org/html/rfc8484
func (d *DOH) Search(domain string) (DNS []net.IP, err error) {
	if x, _ := d.cache.Get(domain); x != nil {
		return x.([]net.IP), nil
	}
	DNS, err = dnsCommon(domain, d.Subnet, func(data []byte) ([]byte, error) { return d.post(data) })
	if err != nil || len(DNS) <= 0 {
		return nil, fmt.Errorf("DNS over HTTPS Search -> %v", err)
	}
	d.cache.Add(domain, DNS)
	return
}

func (d *DOH) SetSubnet(ip *net.IPNet) {
	if ip == nil {
		_, d.Subnet, _ = net.ParseCIDR("0.0.0.0/0")
		return
	}
	d.Subnet = ip
}

func (d *DOH) GetSubnet() *net.IPNet {
	return d.Subnet
}

func (d *DOH) SetServer(host string) {
	d.Server = host
}

func (d *DOH) GetServer() string {
	return d.Server
}

func (d *DOH) SetProxy(proxy func(addr string) (net.Conn, error)) {
	if proxy == nil {
		return
	}
	d.Proxy = proxy
	d.httpClient = &http.Client{
		Transport: &http.Transport{
			//Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				switch network {
				case "tcp":
					return d.Proxy(addr)
				default:
					return net.Dial(network, addr)
				}
			},
			DisableKeepAlives: false,
		},
		Timeout: 10 * time.Second,
	}
}

func (d *DOH) get(dReq []byte) (body []byte, err error) {
	query := strings.Replace(base64.URLEncoding.EncodeToString(dReq), "=", "", -1)
	urls := "https://" + d.Server + "/dns-query?dns=" + query
	res, err := d.httpClient.Get(urls)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

// https://www.cnblogs.com/mafeng/p/7068837.html
func (d *DOH) post(dReq []byte) (body []byte, err error) {
	resp, err := d.httpClient.Post(fmt.Sprintf("https://%s/dns-query", d.Server), "application/dns-message", bytes.NewReader(dReq))
	if err != nil {
		return nil, fmt.Errorf("DOH:post() req -> %v", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("DOH:post() readBody -> %v", err)
	}
	return
}
