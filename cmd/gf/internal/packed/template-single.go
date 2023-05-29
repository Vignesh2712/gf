package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
