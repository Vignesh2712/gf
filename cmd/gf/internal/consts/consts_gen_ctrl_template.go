// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consts

const TemplateGenCtrlControllerEmpty = `
package {Module}

`

const TemplateGenCtrlControllerNewEmpty = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package {Module}

import (
	{ImportPath}
)

`

const TemplateGenCtrlControllerNewFunc = `
type {CtrlName} struct{}

func {NewFuncName}() {InterfaceName} {
	return &{CtrlName}{}
}

`

const TemplateGenCtrlControllerMethodFunc = `
package {Module}

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"{ImportPath}"
)

func (c *{CtrlName}) {MethodName}(ctx context.Context, req *{Version}.{MethodName}Req) (res *{Version}.{MethodName}Res, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
`

const TemplateGenCtrlApiInterface = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package {Module}

import (
{ImportPaths}
)

{Interfaces}
`
