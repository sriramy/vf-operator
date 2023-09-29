package utils

import "github.com/k8snetworkplumbingwg/sriovnet"

func IsSriovVF(pciAddress *string) bool {
	_, err := sriovnet.GetPfPciFromVfPci(*pciAddress)
	return (err == nil)
}
