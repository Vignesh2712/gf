// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localGauge is the local implements for interface Gauge.
type localGauge struct {
	Metric
	MetricConfig
	GaugePerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localGauge)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localGauge)(nil)
)

// NewGauge creates and returns a new Gauge.
func NewGauge(config MetricConfig) (Gauge, error) {
	baseMetric, err := newMetric(MetricTypeGauge, config)
	if err != nil {
		return nil, err
	}
	m := &localGauge{
		Metric:         baseMetric,
		MetricConfig:   config,
		GaugePerformer: newNoopGaugePerformer(),
	}
	if globalProvider != nil {
		if err = m.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, m)
	return m, nil
}

// MustNewGauge creates and returns a new Gauge.
// It panics if any error occurs.
func MustNewGauge(config MetricConfig) Gauge {
	m, err := NewGauge(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localGauge) Init(provider Provider) (err error) {
	if _, ok := l.GaugePerformer.(noopGaugePerformer); !ok {
		// already initialized.
		return
	}
	l.GaugePerformer, err = provider.Performer().Gauge(l.MetricConfig)
	return err
}

// Performer exports internal Performer.
func (l *localGauge) Performer() any {
	return l.GaugePerformer
}

func (*localGauge) observable() {}
