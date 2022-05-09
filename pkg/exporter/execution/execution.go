package execution

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type Node interface {
	Name() string
	URL() string
	Bootstrapped() bool
	Bootstrap(ctx context.Context) error
	SyncStatus(ctx context.Context) (*SyncStatus, error)
	NetworkID(ctx context.Context) (int64, error)
}

type node struct {
	name    string
	url     string
	client  *ethclient.Client
	log     logrus.FieldLogger
	metrics Metrics
}

func NewExecutionNode(ctx context.Context, log logrus.FieldLogger, name string, url string, metrics Metrics) (*node, error) {
	return &node{
		name:    name,
		url:     url,
		log:     log,
		metrics: metrics,
	}, nil
}

func (e *node) Name() string {
	return e.name
}

func (e *node) URL() string {
	return e.url
}

func (e *node) Bootstrapped() bool {
	return e.client != nil
}

func (e *node) Bootstrap(ctx context.Context) error {
	client, err := ethclient.Dial(e.url)
	if err != nil {
		return err
	}

	e.client = client
	return nil
}

func (e *node) SyncStatus(ctx context.Context) (*SyncStatus, error) {
	status, err := e.client.SyncProgress(ctx)
	if err != nil {
		return nil, err
	}

	if status == nil && err == nil {
		// Not syncing
		syncStatus := &SyncStatus{}
		syncStatus.IsSyncing = false
		e.metrics.ObserveSyncStatus(*syncStatus)
		return syncStatus, nil
	}

	syncStatus := &SyncStatus{
		IsSyncing:     true,
		CurrentBlock:  status.CurrentBlock,
		HighestBlock:  status.HighestBlock,
		StartingBlock: status.StartingBlock,
	}

	e.metrics.ObserveSyncStatus(*syncStatus)

	return syncStatus, nil
}

func (e *node) NetworkID(ctx context.Context) (int64, error) {
	id, err := e.client.NetworkID(ctx)
	if err != nil {
		return 0, err
	}

	e.metrics.ObserveNetworkID(id.Int64())

	return id.Int64(), nil
}