package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hay-kot/dynodns/app/lib/httpclient"
	"github.com/hay-kot/dynodns/app/lib/providers/porkbun"
	"github.com/rs/zerolog/log"
)

func (c *Controller) Run(ctx context.Context) (chan struct{}, error) {
	pb := porkbun.New(
		httpclient.New(&http.Client{Transport: httpclient.NewTransport()}, c.PorkBunEndpoint),
		c.PorkBunKey,
		c.PorkBunSecret,
		c.PorkBunSubDomain,
		c.PorkBunDomain,
	)

	// Initialize The Client and Create a Record
	var err error
	var record *porkbun.DNSRecord
	try := 5
	for try > 0 {
		log.Info().Msg("initializing the client")
		record, err = pb.First(ctx)
		if record != nil && err == nil {
			log.Info().Str("recordID", record.ID).Msg("record found")
			break
		}

		if !errors.Is(err, porkbun.ErrNoRecordsFound) {
			log.Err(err).Msg("failed to initialize the client")
			return nil, err
		}

		myIP, err := c.GetPublicIP(ctx)
		if err != nil {
			log.Err(err).Msg("failed to get public ip")
			try--
			continue
		}

		err = pb.CreateRecord(ctx, porkbun.DNSRecord{
			Name:    c.PorkBunSubDomain,
			Content: myIP,
			Type:    "A",
			TTL:     "300",
		})

		if err == nil {
			log.Info().Msg("record created")
			continue
		}

		log.Err(err).Msg("failed to create record, retrying")

		try--
		continue
	}

	if record == nil {
		return nil, fmt.Errorf("failed to initialize the client")
	}

	do := func(ctx context.Context) chan struct{} {
		done := make(chan struct{})

		log.Info().Int("interval", c.Interval).Msg("starting client")
		go func() {
			ticker := time.NewTicker(time.Duration(c.Interval) * time.Second)

			for {
				select {
				case <-ctx.Done():
					log.Info().Msg("stopping client")
					ticker.Stop()
					close(done)
					return
				case <-ticker.C:
					log.Debug().Msg("tick")
					// Replace this with the task you want to run every 5 minutes
					ip, err := c.GetPublicIP(ctx)
					if err != nil {
						log.Err(err).Msg("failed to get public ip")
						continue
					}

					healthcheck := func() {
						if c.PingURL != "" {
							resp, err := c.client.GetCtx(ctx, c.PingURL)
							if err != nil {
								log.Err(err).Msg("failed to ping url")
							}
							_ = resp.Body.Close()
						}
					}

					log.Info().Str("ip", ip).Msg("IP address retrieved")

					if record.Content == ip {
						log.Debug().Msg("ip address has not changed")
						healthcheck()
						continue
					}

					log.Debug().Msg("ip address has changed")

					record := porkbun.DNSRecord{
						ID:      record.ID,
						Name:    c.PorkBunSubDomain,
						Content: ip,
						Type:    "A",
						TTL:     "300",
					}

					err = pb.SetRecord(ctx, record.ID, record)
					if err != nil {
						log.Err(err).Msg("failed to set porkbun dns record")
						continue
					}

					log.Info().Msg("porkbun dns record updated")
					healthcheck()
				}
			}
		}()

		return done
	}

	return do(ctx), nil
}
