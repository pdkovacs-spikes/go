package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"otel-grafana-stack/internal/metadata"
	"runtime/metrics"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	// Wait for signal to close the app
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)

	endpointUrl := "http://alloy:4318"
	endpoint, err := url.Parse(endpointUrl)
	if err != nil {
		panic(fmt.Sprintf("parsing endpoint url: %v", err))
	}
	insecure := endpoint.Scheme == "http"
	// protocol := "http/protobuf"

	var exporter sdkmetric.Exporter

	if insecure {
		exporter, err = otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(endpoint.Host), otlpmetrichttp.WithInsecure())
	} else {
		exporter, err = otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(endpoint.Host))
	}
	if err != nil {
		panic(fmt.Sprintf("parsing endpoint url: %v", err))
	}

	serviceComponent := "main"
	servcieName := "bitkit/otel-grafana-stack"
	serviceNamespace := "bitkit/otel-grafana-stack"
	serviceInstanceID := "local"

	// See also https://github.com/pdkovacs/forked-quickpizza/commit/a5835b3b84d4ae995b8b886a6982a59f3997af2e
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(servcieName),
			attribute.KeyValue{Key: "service.component", Value: attribute.StringValue(serviceComponent)},
			attribute.KeyValue{Key: "service.namespace", Value: attribute.StringValue(serviceNamespace)},
			attribute.KeyValue{Key: "service.instance.id", Value: attribute.StringValue(serviceInstanceID)},
		),
	)
	metricReader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second))

	// Register the prometheus Collector to receive metrics from the go runtime package metrics
	provider := metric.NewMeterProvider(
		sdkmetric.WithReader(metricReader),
		sdkmetric.WithResource(res),
	)
	addBuiltInGoMetricsToOTEL(provider)

	apiCounter := getExampleCustomMetric(provider)
	outcomeKey := attribute.Key("outcome")

	go func() {
		i := 0
		for {
			i++
			outcome := "success"
			if i%3 == 0 {
				outcome = "failure"
			}
			log.Printf("incrementing api.counter by 1\n")
			apiCounter.Add(ctx, 1, api.WithAttributes(outcomeKey.String(outcome)))
			<-time.After(time.Second)
		}
	}()

	<-ctx.Done()
}

// addMetricsToPrometheusRegistry function to add metrics to prometheus registry
func addBuiltInGoMetricsToOTEL(provider *metric.MeterProvider) {

	meter := provider.Meter(metadata.OtelScope)

	// Get descriptions for all supported metrics.
	metricsMeta := metrics.All()

	// Register metrics and retrieve the values in prometheus client
	for i := range metricsMeta {
		// Get metric options
		meta := metricsMeta[i]
		opt := getMetricsOptions(metricsMeta[i])
		name := normalizeOtelName(meta.Name)

		// Register metrics per type of metric
		if meta.Cumulative {
			// Register as a counter
			counter, err := meter.Float64ObservableCounter(name, api.WithDescription(meta.Description))
			if err != nil {
				log.Fatal(err)
			}
			_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
				o.ObserveFloat64(counter, metadata.GetSingleMetricFloat(meta.Name), opt)
				return nil
			}, counter)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Register as a gauge
			gauge, err := meter.Float64ObservableGauge(name, api.WithDescription(meta.Description))
			if err != nil {
				log.Fatal(err)
			}
			_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
				o.ObserveFloat64(gauge, metadata.GetSingleMetricFloat(meta.Name), opt)
				return nil
			}, gauge)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// getMetricsOptions function to get metric labels
func getMetricsOptions(metric metrics.Description) api.MeasurementOption {
	tokens := strings.Split(metric.Name, "/")
	if len(tokens) < 2 {
		return nil
	}

	nameTokens := strings.Split(tokens[len(tokens)-1], ":")
	subsystem := metadata.GetMetricSubsystemName(metric)

	// create a unique name for metric, that will be its primary key on the registry
	opt := api.WithAttributes(
		attribute.Key("Namespace").String(tokens[1]),
		attribute.Key("Subsystem").String(subsystem),
		attribute.Key("Units").String(nameTokens[1]),
	)
	return opt
}

// normalizePrometheusName function to normalize prometheus metric name
func normalizeOtelName(name string) string {
	normalizedName := strings.Replace(name, "/", "", 1)
	normalizedName = strings.Replace(normalizedName, ":", "_", -1)
	normalizedName = strings.TrimSpace(strings.ReplaceAll(normalizedName, "/", "_"))
	return normalizedName
}

func getExampleCustomMetric(provider *metric.MeterProvider) api.Int64Counter {
	meter := provider.Meter(metadata.OtelScope)

	apiCounter, regErr := meter.Int64Counter(
		"api.counter",
		api.WithDescription("Number of API calls."),
		api.WithUnit("{call}"),
	)

	if regErr != nil {
		panic(regErr)
	}

	return apiCounter
}
