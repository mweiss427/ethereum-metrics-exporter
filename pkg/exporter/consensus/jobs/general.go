package jobs

import (
	"context"
	"errors"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type General struct {
	MetricExporter
	client      eth2client.Service
	log         logrus.FieldLogger
	Slots       prometheus.GaugeVec
	NodeVersion prometheus.GaugeVec
	NetworkdID  prometheus.Gauge
}

const (
	NameGeneral = "general"
)

func NewGeneralJob(client eth2client.Service, log logrus.FieldLogger, namespace string, constLabels map[string]string) General {
	constLabels["module"] = NameGeneral
	return General{
		client: client,
		log:    log,
		Slots: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "slot_number",
				Help:        "The slot number of the beacon chain.",
				ConstLabels: constLabels,
			},
			[]string{
				"identifier",
			},
		),
		NodeVersion: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "node_version",
				Help:        "The version of the running beacon node.",
				ConstLabels: constLabels,
			},
			[]string{
				"version",
			},
		),
		NetworkdID: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "network_id",
				Help:        "The network id of the node.",
				ConstLabels: constLabels,
			},
		),
	}
}

func (g *General) Name() string {
	return NameGeneral
}

func (g *General) Start(ctx context.Context) {
	g.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
			g.tick(ctx)
		}
	}
}

func (g *General) tick(ctx context.Context) {
	if err := g.GetNodeVersion(ctx); err != nil {
		g.log.WithError(err).Error("Failed to get node version")
	}

	if err := g.GetBeaconSlot(ctx, "head"); err != nil {
		g.log.WithError(err).Error("Failed to get beacon slot: head")
	}

	if err := g.GetBeaconSlot(ctx, "genesis"); err != nil {
		g.log.WithError(err).Error("Failed to get beacon slot: genesis")
	}

	if err := g.GetBeaconSlot(ctx, "finalized"); err != nil {
		g.log.WithError(err).Error("Failed to get beacon slot: finalized")
	}
}

func (g *General) GetNodeVersion(ctx context.Context) error {
	provider, isProvider := g.client.(eth2client.NodeVersionProvider)
	if !isProvider {
		return errors.New("client does not implement eth2client.NodeVersionProvider")
	}

	version, err := provider.NodeVersion(ctx)
	if err != nil {
		return err
	}

	g.NodeVersion.WithLabelValues(version).Set(1)

	return nil
}

func (g *General) GetBeaconSlot(ctx context.Context, identifier string) error {
	provider, isProvider := g.client.(eth2client.BeaconBlockHeadersProvider)
	if !isProvider {
		return errors.New("client does not implement eth2client.BeaconBlockHeadersProvider")
	}

	block, err := provider.BeaconBlockHeader(ctx, identifier)
	if err != nil {
		return err
	}

	g.ObserveSlot(identifier, uint64(block.Header.Message.Slot))

	return nil
}

func (g *General) ObserveSlot(identifier string, slot uint64) {
	g.Slots.WithLabelValues(identifier).Set(float64(slot))
}

func (g *General) ObserveNetworkID(networkID uint64) {
	g.NetworkdID.Set(float64(networkID))
}