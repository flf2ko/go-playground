package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/flf2ko/playground/go-api-sample/database"
	"github.com/flf2ko/playground/go-api-sample/handlers"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	_ "net/http/pprof"

	_ "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
)

func main() {
	// Initialize database
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize handlers
	linkHandler := handlers.NewLinkHandler(db)

	// Setup Gin router
	router := setupRouter(linkHandler)

	// Setup graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(linkHandler *handlers.LinkHandler) *gin.Engine {
	// Set Gin mode
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.ContextWithFallback = true
	pprof.Register(router)

	lp, tp, mp, _ := setupOTEL()
	defer func() {
		_ = tp.Shutdown(context.Background())
		_ = lp.Shutdown(context.Background())
		_ = mp.Shutdown(context.Background())
	}()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API server is running",
		})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/fetch-json", linkHandler.FetchJSON)
		api.GET("/records", linkHandler.GetRecords)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "go-api-sample",
			"version":     "1.0.0",
			"description": "A simple JSON fetcher API",
			"endpoints": gin.H{
				"health":     "GET /health",
				"metrics":    "GET /metrics",
				"fetch-json": "GET /api/v1/fetch-json?link=<url>",
				"records":    "GET /api/v1/records",
			},
		})
	})

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupOTEL() (*sdktrace.TracerProvider, *sdklog.LoggerProvider, *sdkmetric.MeterProvider, error) {
	tp, err := initTracerProvider()
	if err != nil {
		return nil, nil, nil, err
	}
	lp, err := initLoggerProvider()
	if err != nil {
		return nil, nil, nil, err
	}
	mp, err := initMeterProvider()
	if err != nil {
		return nil, nil, nil, err
	}
	return tp, lp, mp, nil
}

func initTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient())
	if err != nil {
		log.Fatalf("new otlp trace grpc exporter failed: %v", err)
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func initMeterProvider() (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatalf("new otlp metric grpc exporter failed: %v", err)
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource()),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

func initLoggerProvider() (*sdklog.LoggerProvider, error) {
	// tell loki to use host.name and cloud.region as labels
	lokiHint := attribute.KeyValue{
		Key: attribute.Key("loki.resource.labels"),
		Value: attribute.StringSliceValue([]string{
			string(semconv.CloudRegionKey.String("").Key),
			string(semconv.ServiceNameKey.String("").Key),
			string(semconv.HostName("").Key),
		}),
	}

	exporter, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(),
		stdoutlog.WithWriter(os.Stdout),
	)
	if err != nil {
		log.Fatalf("new stdout log exporter failed: %v", err)
		return nil, err
	}
	processor := sdklog.NewBatchProcessor(exporter)

	lp := sdklog.NewLoggerProvider(
		// sdklog.WithBatcher(exporter),
		sdklog.WithProcessor(processor),
		sdklog.WithResource(newResource(lokiHint)),
	)

	global.SetLoggerProvider(lp)

	slog.SetDefault(otelslog.NewLogger("my/pkg/default", otelslog.WithLoggerProvider(lp)))

	return lp, nil
}

func initResource() (resource *sdkresource.Resource) {
	extraResources, _ := sdkresource.New(
		context.Background(),
		sdkresource.WithOS(),
		sdkresource.WithProcess(),
		sdkresource.WithContainer(),
		sdkresource.WithHost(),
		sdkresource.WithDetectors(
			ecs.NewResourceDetector(),
		),
	)
	resource, _ = sdkresource.Merge(
		sdkresource.Default(),
		extraResources,
	)

	return
}

func newResource(extraAttrs ...attribute.KeyValue) *sdkresource.Resource {
	host, _ := os.Hostname()

	attrs := append([]attribute.KeyValue{
		semconv.ServiceNameKey.String("go-playground"),
		// semconv.CloudRegionKey.String(os.Getenv("REGION")),
		semconv.HostName(host),
	}, extraAttrs...)
	return sdkresource.NewWithAttributes(
		semconv.SchemaURL,
		// Note that ServiceNameKey attribute can include chars not allowed in Pyroscope
		// application name, therefore it should be used carefully.
		attrs...,
	)
}
