package disc

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	errMCInterfaceNotAvail = errors.New("multicast interface not available")
)

type vFlowServer struct {
	timestamp int64
}

type DiscoveryConfig struct {
	DiscoveryStrategy string
	Params            map[string]string `yaml:"params,omitempty"`
	Logger            *log.Logger
}

func (c *DiscoveryConfig) LoadConfig(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	yaml.Unmarshal(b, &c.Params)
}

func (c *DiscoveryConfig) GetConfigItem(key string, defaultValue string) string {
	var v, found = c.Params[key]
	if !found {
		return defaultValue
	}
	return v
}

// Discovery represents Discovery interface
type Discovery interface {
	Setup(config *DiscoveryConfig) error
	GetvFlowServers() map[string]vFlowServer
	GetRPCServers() []string
}

func NewDiscovery(config *DiscoveryConfig) Discovery {
	d := &_Discovery{}
	d.vflowServers = make(map[string]vFlowServer, 10)
	return d
}

type _Discovery struct {
	vflowServers map[string]vFlowServer
}

func (d *_Discovery) GetvFlowServers() map[string]vFlowServer {
	return d.vflowServers
}

func (d *_Discovery) GetRPCServers() []string {
	return BuildRpcServersList(d.GetvFlowServers())
}

func (d *_Discovery) Setup(config *DiscoveryConfig) error {
	d.vflowServers = make(map[string]vFlowServer, 10)
	return nil
}

// Utility method to manage vflow servers list
func BuildRpcServersList(vFlowServers map[string]vFlowServer) []string {
	var servers []string

	now := time.Now().Unix()

	// Add locks

	for ip, server := range vFlowServers {
		if now-server.timestamp < 300 {
			servers = append(servers, ip)
		} else {
			delete(vFlowServers, ip)
		}
	}

	return servers
}

func getLocalIPs() (map[string]struct{}, error) {
	ips := make(map[string]struct{})

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifs {
		addrs, err := i.Addrs()
		if err != nil || i.Flags != 19 {
			continue
		}
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			ips[ip.String()] = struct{}{}
		}
	}

	return ips, nil
}

// Simple factory to initialize discovery based on configuration
func BuildDiscovery(config *DiscoveryConfig) (Discovery, error) {

	var discRegistered = map[string]Discovery{
		"vFlowDiscovery":    new(MulticastDiscovery),
		"k8sDiscovery.rest": new(K8SDiscovery),
		//"k8sDiscovery.rest": DNSDiscovery
	}

	disc, ok := discRegistered[config.DiscoveryStrategy]
	if !ok {
		return nil, errors.New("Discovery strategy not found")
	}
	setup_err := disc.Setup(config)
	if setup_err != nil {
		return nil, setup_err
	}
	return disc, nil
}
