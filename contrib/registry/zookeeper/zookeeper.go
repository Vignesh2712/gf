// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"github.com/gogf/gf/v2/net/gsvc"
	"golang.org/x/sync/singleflight"
)

var _ gsvc.Registry = &Registry{}

type Content struct {
	Key   string
	Value string
}

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	namespace string
	user      string
	password  string
}

// WithRootPath with registry root path.
func WithRootPath(path string) Option {
	return func(o *options) { o.namespace = path }
}

// WithDigestACL with registry password.
func WithDigestACL(user string, password string) Option {
	return func(o *options) {
		o.user = user
		o.password = password
	}
}

// Registry is consul registry
type Registry struct {
	opts  *options
	conn  *zk.Conn
	group singleflight.Group
}

func New(conn *zk.Conn, opts ...Option) *Registry {
	options := &options{
		namespace: "/microservices",
	}
	for _, o := range opts {
		o(options)
	}
	return &Registry{
		opts: options,
		conn: conn,
	}
}
