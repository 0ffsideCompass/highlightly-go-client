package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRequestScopedContextDoesNotMutateSharedClient(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if _, err := w.Write([]byte(`[]`)); err != nil {
			t.Error(err)
		}
	}))
	defer srv.Close()
	c, err := NewWithDomain("test-key", srv.URL, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.WithContext(ctx).get(srv.URL); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
	if calls.Load() != 0 {
		t.Fatal("canceled work reached provider")
	}
	if _, err := c.get(srv.URL); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatal("original client context was changed")
	}
}

func TestInFlightProviderRequestUsesCallerDeadline(t *testing.T) {
	started := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
	}))
	defer srv.Close()
	c, err := NewWithDomain("test-key", srv.URL, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err = c.WithContext(ctx).get(srv.URL)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v", err)
	}
	select {
	case <-started:
	default:
		t.Fatal("request never started")
	}
}
