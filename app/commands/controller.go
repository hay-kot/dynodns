// Package commands contains the CLI commands for the application
package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hay-kot/dynodns/app/lib/httpclient"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	Interval int    // Interval in seconds
	PingURL  string // URL to ping to indicate health

	client *httpclient.Client

	PorkBunKey       string
	PorkBunSecret    string
	PorkBunSubDomain string
	PorkBunDomain    string
	PorkBunEndpoint  string
}

// New creates a new Controller with default values
func New() *Controller {
	client := httpclient.New(&http.Client{Transport: httpclient.NewTransport(log.Logger)}, "")
	return &Controller{client: client}
}

// GetPublicIP returns the public IP address of the client using
// https://ifconfig.co/ip to get the IP Address
func (c *Controller) GetPublicIP(ctx context.Context) (string, error) {
	resp, err := c.client.GetCtx(ctx, "https://ifconfig.co/ip")
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	byteIP, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(byteIP))

	if !strings.Contains(ip, ".") {
		return "", fmt.Errorf("invalid ip address: %s", ip)
	}

	return ip, nil
}
