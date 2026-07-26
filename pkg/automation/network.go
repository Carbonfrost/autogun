// Copyright 2026 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package automation

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Carbonfrost/autogun/pkg/config"
	cdpnetwork "github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

type networkState struct {
	eg *errgroup.Group
	mu *sync.Mutex
}

func newNetworkState() *networkState {
	eg := new(errgroup.Group)
	eg.SetLimit(4)
	var mu sync.Mutex

	return &networkState{
		eg: eg,
		mu: &mu,
	}
}

func bindNetworkFilterStart(cfg *config.Network) Task {
	return TaskFunc(func(ctx context.Context) error {
		ar := mustAutomationResult(ctx)
		state := ar.networkState
		eg := state.eg
		mu := state.mu

		ctx, _ = context.WithTimeout(ctx, 2*60*time.Second)
		chromedp.ListenTarget(ctx, func(ev any) {
			e, ok := ev.(*cdpnetwork.EventResponseReceived)
			if !ok || e.Response == nil {
				return
			}

			for _, filter := range cfg.Filters {
				if !matchesNetworkFilter(e, filter) {
					continue
				}
				shouldLog := filter.Log != nil || filter.Download == nil
				if filter.Download != nil {
					requestID := e.RequestID
					responseURL := e.Response.URL

					mu.Lock()
					key := allocateDownloadKey(ar.OutputFiles, responseURL)
					ar.OutputFiles[key] = nil
					mu.Unlock()

					fmt.Printf("Download starting: %s\n", responseURL)
					eg.Go(func() error {
						body, err := cdpnetwork.GetResponseBody(requestID).Do(ctx)
						if err != nil {
							fmt.Printf("Download error: %s: %v\n", responseURL, err)
							mu.Lock()
							delete(ar.OutputFiles, key)
							mu.Unlock()
							return nil
						}
						mu.Lock()
						ar.OutputFiles[key] = &body
						mu.Unlock()
						fmt.Printf("Download complete: %s\n", responseURL)
						return nil
					})
				}
				if shouldLog {
					fmt.Printf("Network filter matched: %s\n", e.Response.URL)
				}
			}
		})
		return nil
	})
}

func bindNetworkFilterStop() Task {
	return TaskFunc(func(ctx context.Context) error {
		ar := mustAutomationResult(ctx)
		state := ar.networkState
		eg := state.eg
		mu := state.mu

		if err := eg.Wait(); err != nil {
			return err
		}
		res := mustAutomationResult(ctx)
		mu.Lock()
		defer mu.Unlock()
		for k, v := range ar.OutputFiles {
			if v != nil {
				data := v
				res.OutputFiles[k] = data
			}
		}
		return nil
	})
}

func matchesNetworkFilter(e *cdpnetwork.EventResponseReceived, filter *config.NetworkFilter) bool {
	if filter.Type != "" && string(e.Type) != filter.Type {
		return false
	}
	if filter.Response != nil && len(filter.Response.URLs) > 0 {
		matched := false
		for _, pattern := range filter.Response.URLs {
			if matchURLPattern(e.Response.URL, pattern) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

func matchURLPattern(responseURL, pattern string) bool {
	u, err := url.Parse(responseURL)
	if err != nil {
		return false
	}
	hostPath := u.Host + u.Path
	matched, _ := path.Match(pattern, hostPath)
	return matched
}

func allocateDownloadKey(downloads map[string]*[]byte, responseURL string) string {
	base := "download"
	if u, err := url.Parse(responseURL); err == nil {
		if b := path.Base(u.Path); b != "" && b != "." && b != "/" {
			base = b
		}
	}
	if _, exists := downloads[base]; !exists {
		return base
	}
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)
	for i := 1; ; i++ {
		key := fmt.Sprintf("%s-%d%s", name, i, ext)
		if _, exists := downloads[key]; !exists {
			return key
		}
	}
}
