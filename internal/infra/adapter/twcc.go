package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TwccAdapter struct {
	client     http.Client
	twccAPIKey string
}

func NewTwccAdapter(twccAPIKey string) *TwccAdapter {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &TwccAdapter{
		client:     client,
		twccAPIKey: twccAPIKey,
	}
}

func (r *TwccAdapter) newClient(method string, requestURL string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "TWCC-CLI")
	req.Header.Set("X-API-HOST", "k8s-D-twcc")
	req.Header.Set("x-API-KEY", r.twccAPIKey)

	return req
}

func (r *TwccAdapter) RunTwccJob(twccJobId string) error {
	requestURL := fmt.Sprintf("https://apigateway.twcc.ai/api/v3/k8s-D-twcc/jobs/%s/submit/", twccJobId)
	req := r.newClient("POST", requestURL, nil)
	resp, err := r.client.Do(req)

	if err != nil {
		return fmt.Errorf("TwccAdapter - RunTwccJob - client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("TwccAdapter - RunTwccJob - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusCreated || string(body) != "" {
		return fmt.Errorf("TwccAdapter - RunTwccJob - resp.StatusCode: %w - Create fail ", err)
	}

	return nil
}

type TwccJobResponse struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Project  int    `json:"project"`
	Status   string `json:"status"`
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	Callback string `json:"callback"`
	// ... include other fields if necessary
}

func (r *TwccAdapter) GetTwccJobStatus(twccJobId string) (string, error) {
	requestURL := fmt.Sprintf("https://apigateway.twcc.ai/api/v3/k8s-D-twcc/jobs/%s/", twccJobId)
	req := r.newClient("GET", requestURL, nil)
	resp, err := r.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("TwccAdapter - RunTwccJob - client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("TwccAdapter - RunTwccJob - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK || string(body) == "" {
		return "", fmt.Errorf("TwccAdapter - RunTwccJob - resp.StatusCode: %w", err)
	}

	var job TwccJobResponse
	if err := json.Unmarshal(body, &job); err != nil {
		return "", err
	}

	// Return the status field
	return job.Status, nil
}

type CreateTwccCCSResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	// ... include other fields if necessary
}

func (r *TwccAdapter) CreateTwccCCS() (string, error) {
	requestURL := "https://apigateway.twcc.ai/api/v2/k8s-D-twcc/sites/"
	body := bytes.NewBuffer([]byte(`{
		"name": "inference-service",
		"desc": "inference-service created GPU container",
		"project": 65662,
		"solution": 4
	}`))
	req := r.newClient("POST", requestURL, body)
	req.Header.Set("x-extra-property-flavor", "1 GPU + 04 cores + 090GB memory")
	req.Header.Set("x-extra-property-image", "tensorflow-23.08-tf2-py3:latest")
	req.Header.Set("x-extra-property-replica", "1")
	req.Header.Set("x-extra-property-gpfs02-mount-path", "/home/yjack0000")
	req.Header.Set("x-extra-property-gpfs01-mount-path", "/work/yjack0000")
	resp, err := r.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("TwccAdapter - CreateTwccCCS - client.Do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("TwccAdapter - CreateTwccCCS - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusCreated || string(respBody) == "" {
		return "", fmt.Errorf("TwccAdapter - CreateTwccCCS - resp.StatusCode: %d", resp.StatusCode)
	}

	var ccs CreateTwccCCSResponse
	if err := json.Unmarshal(respBody, &ccs); err != nil {
		return "", err
	}

	return fmt.Sprint(ccs.ID), nil
}

type GetTwccCCSStatusResponse struct {
	Service []struct {
		Name      string   `json:"name"`
		NetType   string   `json:"net_type"`
		ClusterIP string   `json:"cluster_ip"`
		PublicIP  []string `json:"public_ip"`
		// ... include other fields if necessary
		Ports []struct {
			Name       string `json:"name"`
			TargetPort int    `json:"target_port"`
			Port       int    `json:"port"`
			Protocol   string `json:"protocol"`
			NodePort   int    `json:"node_port"`
		} `json:"ports"`
	} `json:"Service"`
	Pod []struct {
		Name string `json:"name"`
		// ... include other fields if necessary
	} `json:"Pod"`
	/// ... include other fields if necessary
}

func (r *TwccAdapter) GetTwccCCSEntryPoint(twccCCSId string) (string, error) {
	requestURL := fmt.Sprintf("https://apigateway.twcc.ai/api/v3/k8s-D-twcc/sites/%s/container/", twccCCSId)
	req := r.newClient("GET", requestURL, nil)
	resp, err := r.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSStatus - client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSStatus - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK || string(body) == "" {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSStatus - resp.StatusCode: %w", err)
	}

	var ccs GetTwccCCSStatusResponse
	if err := json.Unmarshal(body, &ccs); err != nil {
		return "", err
	}

	if len(ccs.Service) == 0 {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSEntryPoint - len(ccs.Service) == 0")
	}

	// Return Public IP and Port concated
	// Ports[2] Refer to the port of 5000
	entryPoint := ccs.Service[0].PublicIP[0] + ":" + fmt.Sprint(ccs.Service[0].Ports[2].Port)

	return entryPoint, nil
}

func (r *TwccAdapter) getTwccCCSPodName(twccCCSId string) (string, error) {
	requestURL := fmt.Sprintf("https://apigateway.twcc.ai/api/v3/k8s-D-twcc/sites/%s/container/", twccCCSId)
	req := r.newClient("GET", requestURL, nil)
	resp, err := r.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSStatus - client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSStatus - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK || string(body) == "" {
		return "", fmt.Errorf("TwccAdapter - GetTwccCCSStatus - resp.StatusCode: %w", err)
	}

	var ccs GetTwccCCSStatusResponse
	if err := json.Unmarshal(body, &ccs); err != nil {
		return "", err
	}

	return ccs.Pod[0].Name, nil
}

func (r *TwccAdapter) TwccCCSAssociateIP(twccCCSId string) error {
	podName, err := r.getTwccCCSPodName(twccCCSId)
	if err != nil {
		return fmt.Errorf("TwccAdapter - TwccCCSAssociateIP - r.getTwccCCSPodName: %w", err)
	}

	requestURL := fmt.Sprintf("https://apigateway.twcc.ai/api/v3/k8s-D-twcc/sites/%s/container/action/", twccCCSId)
	body := bytes.NewBuffer([]byte(`{
		"pod_name": "` + podName + `",
		"action": "associateIP",
		"ports": [
			{
			"targetPort": 5000
			}
		]}`))
	req := r.newClient("PUT", requestURL, body)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("TwccAdapter - TwccCCSAssociateIP - client.Do: err: %w", err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return nil
}

func (r *TwccAdapter) DeleteTwccCCS(twccCSSId string) error {
	requestURL := fmt.Sprintf("https://apigateway.twcc.ai/api/v3/k8s-D-twcc/sites/%s/", twccCSSId)
	fmt.Println(requestURL)
	req := r.newClient("DELETE", requestURL, nil)
	resp, err := r.client.Do(req)

	if err != nil {
		return fmt.Errorf("TwccAdapter - DeleteTwccCCS - client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("TwccAdapter - DeleteTwccCCS - ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent || string(body) != "" {
		return fmt.Errorf("TwccAdapter - DeleteTwccCCS - resp.StatusCode: %w", err)
	}

	return nil
}
