package cmpopts

import (
	"strings"

	"github.com/google/go-cmp/cmp"
)

var IgnoreInternalProtbufFieldsOption = cmp.FilterPath(func(p cmp.Path) bool {
	return strings.HasPrefix(p.Last().String(), ".XXX")
}, cmp.Ignore())
