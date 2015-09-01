package handle

import (
	"github.com/syhlion/gusher/core"
	"github.com/syhlion/gusher/model"
)

type Handler struct {
	AppData    *model.AppData
	Collection *core.Collection
}
