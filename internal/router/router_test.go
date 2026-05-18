package router_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/router"
)

func makeEntries() []parser.Entry {
	now := time.Now()
	return []parser.Entry{
		{Timestamp: now, Level: parser.LevelError, Message: "disk full"},
		{Timestamp: now, Level: parser.LevelWarn, Message: "high memory"},
		{Timestamp: now, Level: parser.LevelInfo, Message: "started"},
		{Timestamp: now, Level: parser.LevelError, Message: "timeout"},
	}
}

func TestRouter_Dispatch_FiltersPerRoute(t *testing.T) {
	errorBuf := &bytes.Buffer{}
	infoBuf := &bytes.Buffer{}

	routes := []router.Route{
		{Name: "errors", Options: filter.NewOptionsBuilder().WithLevel("ERROR").Build(), Format: "json"},
		{Name: "info", Options: filter.NewOptionsBuilder().WithLevel("INFO").Build(), Format: "json"},
	}

	r := router.New(routes)
	factory := router.MapWriterFactory(map[string]interface{ Write([]byte) (int, error) }{
		"errors": errorBuf,
		"info":   infoBuf,
	})

	if err := r.Dispatch(makeEntries(), factory); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(errorBuf.String(), "disk full") {
		t.Error("errors route should contain 'disk full'")
	}
	if strings.Contains(errorBuf.String(), "started") {
		t.Error("errors route should not contain 'started'")
	}
	if !strings.Contains(infoBuf.String(), "started") {
		t.Error("info route should contain 'started'")
	}
}

func TestRouter_Dispatch_InvalidFormat(t *testing.T) {
	routes := []router.Route{
		{Name: "bad", Options: filter.Options{}, Format: "xml"},
	}
	r := router.New(routes)
	factory := router.MapWriterFactory(map[string]interface{ Write([]byte) (int, error) }{})

	if err := r.Dispatch(makeEntries(), factory); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestRouter_Summary_CountsPerRoute(t *testing.T) {
	routes := []router.Route{
		{Name: "errors", Options: filter.NewOptionsBuilder().WithLevel("ERROR").Build(), Format: "json"},
		{Name: "all", Options: filter.Options{}, Format: "json"},
	}
	r := router.New(routes)
	summary := r.Summary(makeEntries())

	if summary["errors"].Total != 2 {
		t.Errorf("expected 2 errors, got %d", summary["errors"].Total)
	}
	if summary["all"].Total != 4 {
		t.Errorf("expected 4 total, got %d", summary["all"].Total)
	}
}

func TestRouter_Dispatch_EmptyEntries(t *testing.T) {
	buf := &bytes.Buffer{}
	routes := []router.Route{
		{Name: "errors", Options: filter.NewOptionsBuilder().WithLevel("ERROR").Build(), Format: "json"},
	}
	r := router.New(routes)
	factory := router.MapWriterFactory(map[string]interface{ Write([]byte) (int, error) }{
		"errors": buf,
	})
	if err := r.Dispatch([]parser.Entry{}, factory); err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
}
