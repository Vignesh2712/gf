// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	ctx = context.TODO()
)

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(ctx, func(ctx context.Context) {
			panic(s)
		}), s)
	})
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(ctx, func(ctx context.Context) {
			panic(gerror.New(s))
		}), s)
	})
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(ctx, func(ctx context.Context) {
			panic("gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(ctx, func(ctx context.Context) {
			panic("gutil TryCatch test")

		}, func(ctx context.Context, err error) {
			t.Assert(err, "gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(ctx, func(ctx context.Context) {
			panic(gerror.New("gutil TryCatch test"))

		}, func(ctx context.Context, err error) {
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Keys(new(T))
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// *struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// **struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(&pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", values)
		t.AssertIN("2", values)
	})

	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, g.Slice{"1", 2})
	})
}
