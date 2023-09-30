package devices

import (
	"github.com/sriramy/vf-operator/pkg/config"
	"github.com/sriramy/vf-operator/pkg/utils"
)

func (d *NetDevice) Configure(c *config.ResourceConfig) error {
	err := utils.SetNumVfs(d.PCIAddress, c.NumVfs)
	if err != nil {
		return err
	}
	return nil
}
