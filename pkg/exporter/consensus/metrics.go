package consensus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samcm/ethereum-metrics-exporter/pkg/exporter/consensus/jobs"
)

type Metrics interface {
	ObserveSyncStatus(status SyncStatus)
	ObserveNodeVersion(version string)
	ObserveSpec(spec map[string]interface{})
}

type metrics struct {
	nodeVersion *prometheus.GaugeVec

	syncMetrics jobs.SyncStatus
	specMetrics jobs.Spec
}

func NewMetrics(nodeName, namespace string) Metrics {
	constLabels := make(prometheus.Labels)
	constLabels["ethereum_role"] = "consensus"
	constLabels["node_name"] = nodeName

	m := &metrics{

		nodeVersion: prometheus.NewGaugeVec(
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
		specMetrics: jobs.NewSpec(namespace, constLabels),
		syncMetrics: jobs.NewSyncStatus(namespace, constLabels),
	}

	prometheus.MustRegister(m.syncMetrics.Percentage)
	prometheus.MustRegister(m.syncMetrics.EstimatedHighestSlot)
	prometheus.MustRegister(m.syncMetrics.HeadSlot)
	prometheus.MustRegister(m.syncMetrics.Distance)
	prometheus.MustRegister(m.syncMetrics.IsSyncing)
	prometheus.MustRegister(m.nodeVersion)

	prometheus.MustRegister(m.specMetrics.SafeSlotsToUpdateJustified)
	prometheus.MustRegister(m.specMetrics.DepositChainID)
	prometheus.MustRegister(m.specMetrics.ConfigName)
	prometheus.MustRegister(m.specMetrics.MaxValidatorsPerCommittee)
	prometheus.MustRegister(m.specMetrics.SecondsPerEth1Block)
	prometheus.MustRegister(m.specMetrics.BaseRewardFactor)
	prometheus.MustRegister(m.specMetrics.EpochsPerSyncCommitteePeriod)
	prometheus.MustRegister(m.specMetrics.EffectiveBalanceIncrement)
	prometheus.MustRegister(m.specMetrics.MaxAttestations)
	prometheus.MustRegister(m.specMetrics.MinSyncCommitteeParticipants)
	prometheus.MustRegister(m.specMetrics.GenesisDelay)
	prometheus.MustRegister(m.specMetrics.SecondsPerSlot)
	prometheus.MustRegister(m.specMetrics.MaxEffectiveBalance)
	prometheus.MustRegister(m.specMetrics.TerminalTotalDifficulty)
	prometheus.MustRegister(m.specMetrics.MaxDeposits)
	prometheus.MustRegister(m.specMetrics.MinGenesisActiveValidatorCount)
	prometheus.MustRegister(m.specMetrics.TargetCommitteeSize)
	prometheus.MustRegister(m.specMetrics.SyncCommitteeSize)
	prometheus.MustRegister(m.specMetrics.Eth1FollowDistance)
	prometheus.MustRegister(m.specMetrics.TerminalBlockHashActivationEpoch)
	prometheus.MustRegister(m.specMetrics.MinDepositAmount)
	prometheus.MustRegister(m.specMetrics.SlotsPerEpoch)
	prometheus.MustRegister(m.specMetrics.PresetBase)

	return m
}

func (m *metrics) ObserveNodeVersion(version string) {
	m.nodeVersion.WithLabelValues(version).Set(float64(1))
}

func (m *metrics) ObserveSpec(spec map[string]interface{}) {
	m.specMetrics.Update(spec)
}

func (m *metrics) ObserveSyncStatus(status SyncStatus) {
	m.syncMetrics.ObserveSyncDistance(status.SyncDistance)
	m.syncMetrics.ObserveSyncEstimatedHighestSlot(status.EstimatedHeadSlot)
	m.syncMetrics.ObserveSyncHeadSlot(status.HeadSlot)
	m.syncMetrics.ObserveSyncIsSyncing(status.IsSyncing)
	m.syncMetrics.ObserveSyncPercentage(status.Percent())
}