// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"time"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
)

// localProvider implements interface gmetric.Provider.
type localProvider struct {
	*metric.MeterProvider
}

// newLocalProvider creates and returns an object that implements gmetric.Provider.
func newLocalProvider(options ...metric.Option) (gmetric.Provider, error) {
	// TODO global logger set for otel
	// otel.SetLogger()

	var (
		err          error
		metrics      = gmetric.GetAllMetrics()
		metricViews  = createViewsByMetrics(metrics)
		builtinViews = createViewsForBuiltInMetrics()
		callbacks    = gmetric.GetRegisteredCallbacks()
	)
	metricViews = append(metricViews, builtinViews...)
	options = append(options, metric.WithView(metricViews...))
	provider := &localProvider{
		MeterProvider: metric.NewMeterProvider(options...),
	}
	for _, callback := range callbacks {
		for _, m := range callback.Metrics {
			hasGlobalCallbackMetricSet.Add(m.Info().Key())
		}
	}
	initializeAllMetrics(metrics, provider)
	err = initializeCallback(callbacks, provider.MeterProvider)
	if err != nil {
		return nil, err
	}

	// builtin metrics: golang.
	err = runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
		runtime.WithMeterProvider(provider),
	)
	if err != nil {
		return nil, gerror.WrapCode(
			gcode.CodeInternalError, err, `start runtime metrics failed`,
		)
	}

	return provider, nil
}

// SetAsGlobal sets current provider as global meter provider for current process.
func (l *localProvider) SetAsGlobal() {
	gmetric.SetGlobalProvider(l)
	otel.SetMeterProvider(l)
}

// Performer creates and returns a Performer.
// A Performer can produce types of Metric performer.
func (l *localProvider) Performer() gmetric.Performer {
	return newPerformer(l.MeterProvider)
}

// createViewsForBuiltInMetrics creates and returns views for builtin metrics.
func createViewsForBuiltInMetrics() []metric.View {
	var views = make([]metric.View, 0)
	views = append(views, metric.NewView(
		metric.Instrument{
			Name: "process.runtime.go.gc.pause_ns",
			Scope: instrumentation.Scope{
				Name:    runtime.ScopeName,
				Version: runtime.Version(),
			},
		},
		metric.Stream{
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{
					500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000,
				},
			},
		},
	))
	views = append(views, metric.NewView(
		metric.Instrument{
			Name: "runtime.uptime",
			Scope: instrumentation.Scope{
				Name:    runtime.ScopeName,
				Version: runtime.Version(),
			},
		},
		metric.Stream{
			Name: "process.runtime.uptime",
		},
	))
	return views
}

// createViewsByMetrics creates and returns OpenTelemetry metric.View according metric type for all metrics,
// especially the Histogram which needs a metric.View to custom buckets.
func createViewsByMetrics(metrics []gmetric.Metric) []metric.View {
	var views = make([]metric.View, 0)
	for _, m := range metrics {
		switch m.Info().Type() {
		case gmetric.MetricTypeCounter:
		case gmetric.MetricTypeGauge:
		case gmetric.MetricTypeHistogram:
			// Custom buckets for each Histogram.
			views = append(views, metric.NewView(
				metric.Instrument{
					Name: m.Info().Name(),
					Scope: instrumentation.Scope{
						Name:    m.Info().Instrument().Name(),
						Version: m.Info().Instrument().Version(),
					},
				},
				metric.Stream{
					Aggregation: metric.AggregationExplicitBucketHistogram{
						Boundaries: m.(gmetric.Histogram).Buckets(),
					},
				},
			))
		}
	}
	return views
}

// initializeAllMetrics initializes all metrics in provider creating.
// The initialization replaces the underlying metric performer using noop-performer with truly performer
// that implements operations for types of metric.
func initializeAllMetrics(metrics []gmetric.Metric, provider gmetric.Provider) {
	for _, m := range metrics {
		if initializer, ok := m.(gmetric.MetricInitializer); ok {
			initializer.Init(provider)
		}
	}
}

func initializeCallback(callbacks []gmetric.GlobalCallbackItem, provider *metric.MeterProvider) error {
	for _, callback := range callbacks {
		// group the metric by instrument and instrument version.
		var (
			instSet  = gset.NewStrSet()
			meterMap = map[otelmetric.Meter][]otelmetric.Observable{}
		)
		for _, m := range callback.Metrics {
			var meter = provider.Meter(
				m.Info().Instrument().Name(),
				otelmetric.WithInstrumentationVersion(m.Info().Instrument().Version()),
			)
			instSet.Add(fmt.Sprintf(
				`%s@%s`,
				m.Info().Instrument().Name(),
				m.Info().Instrument().Version(),
			))
			if _, ok := meterMap[meter]; !ok {
				meterMap[meter] = make([]otelmetric.Observable, 0)
			}
			meterMap[meter] = append(meterMap[meter], metricToFloat64Observable(m))
		}
		if len(meterMap) > 1 {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`multiple instrument or instrument version metrics used in the same callback: %s`,
				instSet.Join(","),
			)
		}
		// do callback registering.
		for meter, observables := range meterMap {
			_, err := meter.RegisterCallback(
				func(ctx context.Context, observer otelmetric.Observer) error {
					return callback.Callback(ctx, newCallbackSetter(observer))
				},
				observables...,
			)
			if err != nil {
				return gerror.WrapCode(
					gcode.CodeInternalError, err,
					`RegisterCallback failed`,
				)
			}
		}
	}
	return nil
}