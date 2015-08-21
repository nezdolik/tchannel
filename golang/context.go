// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tchannel

import (
	"time"

	"golang.org/x/net/context"
)

type contextKey int

const (
	contextKeyUnknown contextKey = iota
	contextKeyTracing
	contextKeyCall
)

// IncomingCall exposes properties for incoming calls through the context.
type IncomingCall interface {
	// CallerName returns the caller name from the CallerName transport header.
	CallerName() string
}

// NewContext returns a new root context used to make TChannel requests.
func NewContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	tctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx := context.WithValue(tctx, contextKeyTracing, NewRootSpan())
	return ctx, cancel
}

// WrapContextForTest returns a copy of the given Context that is associated with the call.
// This should be used in units test only.
func WrapContextForTest(ctx context.Context, call IncomingCall) context.Context {
	return context.WithValue(ctx, contextKeyCall, call)
}

// newIncomingContext creates a new context for an incoming call with the given span.
func newIncomingContext(call IncomingCall, timeout time.Duration, span *Span) (context.Context, context.CancelFunc) {
	tctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx := context.WithValue(tctx, contextKeyTracing, span)
	ctx = context.WithValue(ctx, contextKeyCall, call)
	return ctx, cancel
}

// CurrentCall returns the current incoming call, or nil if this is not an incoming call context.
func CurrentCall(ctx context.Context) IncomingCall {
	if v := ctx.Value(contextKeyCall); v != nil {
		return v.(IncomingCall)
	}
	return nil
}
