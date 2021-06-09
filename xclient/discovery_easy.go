package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type EasyRegistryDiscovery struct {
	*MultiServersDiscovery               // 嵌套了 MultiServersDiscovery，很多能力可以复用。
	registry               string        // 注册中心的地址
	timeout                time.Duration // 服务列表的过期时间
	lastUpdate             time.Time     // 是代表最后从注册中心更新服务列表的时间，默认 10s 过期，即 10s 之后，需要从注册中心更新新的列表。
}

const defaultUpdateTimeout = time.Second * 10

func NewEasyRegistryDiscovery(registerAddr string, timeout time.Duration) *EasyRegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	d := &EasyRegistryDiscovery{
		MultiServersDiscovery: NewMultiServerDiscovery(make([]string, 0)),
		registry:              registerAddr,
		timeout:               timeout,
	}
	return d
}

func (d *EasyRegistryDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *EasyRegistryDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}
	log.Println("rpc registry: refresh servers from registry", d.registry)
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return err
	}
	servers := strings.Split(resp.Header.Get("X-Easyrpc-Servers"), ",")
	d.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		if strings.TrimSpace(server) != "" {
			d.servers = append(d.servers, strings.TrimSpace(server))
		}
	}
	d.lastUpdate = time.Now()
	return nil
}

func (d *EasyRegistryDiscovery) Get(mode SelectMode) (string, error) {
	// 先调用 Refresh 确保服务列表没有过期。
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

func (d *EasyRegistryDiscovery) GetAll() ([]string, error) {
	// 先调用 Refresh 确保服务列表没有过期。
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}
