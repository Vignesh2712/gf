// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	defaultSeparator = "/"
	delimiter        = ","
)

// Separator is the default defaultSeparator for path.
func Separator() string {
	return defaultSeparator
}

// NewServiceWithName creates and returns service from `name`.
func NewServiceWithName(name string) (s *Service) {
	s = &Service{
		Name:     name,
		Metadata: make(Metadata),
	}
	s.autoFillDefaultAttributes()
	return
}

// NewServiceWithKV creates and returns service from `key` and `value`.
func NewServiceWithKV(key, value []byte) (s *Service, err error) {
	array := gstr.Split(gstr.Trim(string(key), defaultSeparator), defaultSeparator)
	if len(array) < 6 {
		err = gerror.NewCodef(gcode.CodeInvalidParameter, `invalid service key "%s"`, key)

		return
	}
	s = &Service{
		Prefix:     array[0],
		Deployment: array[1],
		Namespace:  array[2],
		Name:       array[3],
		Version:    array[4],
		Endpoints:  gstr.Split(array[5], delimiter),
		Metadata:   make(Metadata),
		Separator:  defaultSeparator,
	}
	s.autoFillDefaultAttributes()
	if len(value) > 0 {
		if err = gjson.Unmarshal(value, &s.Metadata); err != nil {
			err = gerror.WrapCodef(gcode.CodeInvalidParameter, err, `invalid service value "%s"`, value)
			return nil, err
		}
	}
	return s, nil
}

// Key formats the service information and returns the Service as registering key.
func (s *Service) Key() string {
	if s.Separator == "" {
		s.Separator = defaultSeparator
	}
	serviceNameUnique := s.KeyWithoutEndpoints()
	serviceNameUnique += s.Separator + gstr.Join(s.Endpoints, ",")
	return serviceNameUnique
}

// KeyWithSchema formats the service information and returns the Service as dialing target key.
func (s *Service) KeyWithSchema() string {
	return fmt.Sprintf(`%s://%s`, Schema, s.Key())
}

// KeyWithoutEndpoints formats the service information and returns a string as a unique name of service.
func (s *Service) KeyWithoutEndpoints() string {
	s.autoFillDefaultAttributes()
	if s.Separator == "" {
		s.Separator = defaultSeparator
	}
	if s.Separator != defaultSeparator {
		return gstr.Join([]string{s.Prefix, s.Deployment, s.Namespace, s.Name, s.Version}, s.Separator)
	}
	return s.Separator + gstr.Join([]string{s.Prefix, s.Deployment, s.Namespace, s.Name, s.Version}, s.Separator)
}

// Value formats the service information and returns the Service as registering value.
func (s *Service) Value() string {
	b, err := gjson.Marshal(s.Metadata)
	if err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	return string(b)
}

// Address returns the first endpoint of Service.
// Eg: 192.168.1.12:9000.
func (s *Service) Address() string {
	if len(s.Endpoints) == 0 {
		return ""
	}
	return s.Endpoints[0]
}

func (s *Service) autoFillDefaultAttributes() {
	if s.Prefix == "" {
		s.Prefix = gcmd.GetOptWithEnv(EnvPrefix, DefaultPrefix).String()
	}
	if s.Deployment == "" {
		s.Deployment = gcmd.GetOptWithEnv(EnvDeployment, DefaultDeployment).String()
	}
	if s.Namespace == "" {
		s.Namespace = gcmd.GetOptWithEnv(EnvNamespace, DefaultNamespace).String()
	}
	if s.Version == "" {
		s.Version = gcmd.GetOptWithEnv(EnvVersion, DefaultVersion).String()
	}
	if s.Name == "" {
		s.Name = gcmd.GetOptWithEnv(EnvName).String()
	}
}
