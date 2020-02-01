package gctscript

import (
	"github.com/idoall/gocryptotrader/gctscript/modules"
	"github.com/idoall/gocryptotrader/gctscript/wrappers/gct"
)

// Setup configures the wrapper interface to use
func Setup() {
	modules.SetModuleWrapper(gct.Setup())
}
