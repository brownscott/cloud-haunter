package operation

import (
	ctx "github.com/hortonworks/cloud-haunter/context"
	"github.com/hortonworks/cloud-haunter/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	ctx.Operations[types.Instances] = instances{}
}

type instances struct {
}

func (o instances) Execute(clouds []types.CloudType) []types.CloudItem {
	log.Debugf("[GET_INSTANCES] Collecting long instances instances on: [%s]", clouds)
	itemsChan, errChan := collectInstances(clouds)
	return wait(itemsChan, errChan, "[GET_INSTANCES] Failed to collect instances")
}
