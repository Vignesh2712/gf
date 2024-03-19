// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	ctx = context.TODO()
)

func TestNewVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).Int(), 1)
		t.Assert(g.NewVar(1, true).Int(), 1)
	})
}

func TestDump(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.Dump("GoFrame")
	})
}

func TestDumpTo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.DumpTo(os.Stdout, "GoFrame", gutil.DumpOption{})
	})
}

func TestDumpWithType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.DumpWithType("GoFrame", 123)
	})
}

func TestDumpWithOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.DumpWithOption("GoFrame", gutil.DumpOption{})
	})
}

func TestTry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.Try(ctx, func(ctx context.Context) {
			g.Dump("GoFrame")
		})
	})
}

func TestTryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.TryCatch(ctx, func(ctx context.Context) {
			g.Dump("GoFrame")
		}, func(ctx context.Context, exception error) {
			g.Dump(exception)
		})
	})
	gtest.C(t, func(t *gtest.T) {
		g.TryCatch(ctx, func(ctx context.Context) {
			g.Throw("GoFrame")
		}, func(ctx context.Context, exception error) {
			t.Assert(exception.Error(), "GoFrame")
		})
	})
}

func TestIsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.IsNil(nil), true)
		t.Assert(g.IsNil(0), false)
		t.Assert(g.IsNil("GoFrame"), false)
	})
}

func TestIsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.IsEmpty(nil), true)
		t.Assert(g.IsEmpty(0), true)
		t.Assert(g.IsEmpty("GoFrame"), false)
	})
}

func TestSetDebug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.SetDebug(true)
	})
}

func TestObject(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(g.Client(), nil)
		t.AssertNE(g.Server(), nil)
		t.AssertNE(g.TCPServer(), nil)
		t.AssertNE(g.UDPServer(), nil)
		t.AssertNE(g.View(), nil)
		t.AssertNE(g.Config(), nil)
		t.AssertNE(g.Cfg(), nil)
		t.AssertNE(g.Resource(), nil)
		t.AssertNE(g.I18n(), nil)
		t.AssertNE(g.Res(), nil)
		t.AssertNE(g.Log(), nil)
		t.AssertNE(g.Validator(), nil)
	})
}

func TestGo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			wg    = sync.WaitGroup{}
			array = garray.NewArray(true)
		)
		wg.Add(1)
		g.Go(context.Background(), func(ctx context.Context) {
			defer wg.Done()
			array.Append(1)
		}, nil)
		wg.Wait()
		t.Assert(array.Len(), 1)
	})
}
