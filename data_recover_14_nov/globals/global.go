package globals

import (
	"data_recover_14_nov/model"
	"sync"
)

var (
	ApplicationConfig *model.Config
	DataMap           *sync.Map
)
