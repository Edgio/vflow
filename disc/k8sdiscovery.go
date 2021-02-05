package disc

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

var (
	logger *log.Logger
)

// K8SDiscovery represents k8s configuration for discovery of vflow nodes
type K8SDiscovery struct {
	discovery          Discovery
	k8sApiServer       string
	k8sCertificatePath string
	k8sNamespace       string
	k8sServiceName     string
	k8sToken           string
	pollInterval       int
	restClient         *resty.Client
	GetPodsEndpoint    string
}

func (d *K8SDiscovery) GetvFlowServers() map[string]vFlowServer {
	return d.discovery.GetvFlowServers()
}

func (d *K8SDiscovery) GetRPCServers() []string {
	return BuildRpcServersList(d.GetvFlowServers())
}

// NewK8SDiscovery constructs and initializes K8S discovery.
func (k *K8SDiscovery) Setup(config *DiscoveryConfig) error {

	logger = config.Logger

	k.discovery = NewDiscovery(config)
	k.k8sApiServer = "kubernetes.default.svc"
	k.k8sServiceName = config.Params["k8s-service-name"]
	k.k8sCertificatePath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

	b, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		logger.Printf("k8s discovery failed: %s\n", err)
		return err
	}
	k.k8sToken = string(b)

	b, err = ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		logger.Printf("k8s discovery failed: %s\n", err)
		return err
	}
	k.k8sNamespace = string(b)

	// Start with default pollinterval of 10s
	k.pollInterval = 10
	if val, ok := config.Params["k8s-discovery-poll-interval"]; ok {
		k.pollInterval, _ = strconv.Atoi(val)
	}

	go k.runDiscovery()

	return nil
}

func (d *K8SDiscovery) runDiscovery() {
	logger.Printf("Starting k8s pod discovery using REST API (%d)\n", d.pollInterval)
	tick := time.NewTicker(time.Duration(d.pollInterval) * time.Second)

	q := url.QueryEscape(fmt.Sprintf("app.kubernetes.io/instance=%s", d.k8sServiceName))

	d.GetPodsEndpoint = fmt.Sprintf("https://%s/api/v1/namespaces/%s/pods?labelSelector=%s",
		d.k8sApiServer, d.k8sNamespace, q)

	d.restClient = resty.New()
	d.restClient.SetRootCertificate(d.k8sCertificatePath)
	for {
		<-tick.C
		d.pollForServers()
	}
}

func (d *K8SDiscovery) pollForServers() {
	resp, err := d.restClient.R().
		EnableTrace().
		SetAuthToken(d.k8sToken).
		SetHeader("Accept", "application/json").
		Get(d.GetPodsEndpoint)

	laddrs, err := getLocalIPs()

	if err != nil {
		logger.Printf("Error %s\n", err)
	} else {
		srch_result := gjson.Get(resp.String(), "items.#.status.podIP")
		pod_ips := srch_result.Array()
		for _, s := range pod_ips {
			if _, ok := laddrs[s.Str]; ok {
				continue
			}
			d.GetvFlowServers()[s.Str] = vFlowServer{time.Now().Unix()}
		}
	}
	logger.Printf("Discovered servers %+v", d.GetvFlowServers())
}
