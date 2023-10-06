/*
 Copyright (c) 2023 Sriram Yagaraman

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/invoke"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/utils/buildversion"
)

var defaultConfDir = "/etc/cni/vf-operator/net.d"

func init() {
	// this ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()
}

// NetConf extends types.NetConf for vf-cni
type NetConf struct {
	types.NetConf
	ConfDir       string `json:"confDir"`
	ResourceName  string `json:"resourceName"`
	InterfaceName string `json:"interfaceName"`
	MAC           string `json:"mac"`
	Vlan          uint32 `json:"vlan"`
}

func getNetConf(args *skel.CmdArgs) (NetConf, error) {
	conf := NetConf{}

	if err := json.Unmarshal(args.StdinData, &conf); err != nil {
		return NetConf{}, fmt.Errorf("failed to load netconf: %v", err)
	}

	return conf, nil
}

func findCNI(confDir string, plugin string) (*libcni.NetworkConfig, string, error) {
	exec := &invoke.RawExec{Stderr: os.Stderr}
	paths := filepath.SplitList(os.Getenv("CNI_PATH"))

	files, err := libcni.ConfFiles(confDir, []string{".conf", ".json"})
	if err != nil {
		return nil, "", fmt.Errorf("Error finding CNI config files: %v", err)
	}

	for _, filename := range files {
		conf, err := libcni.ConfFromFile(filename)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load CNI config file %s: %v", filename, err)
		}
		if conf.Network.Type == plugin {
			path, err := exec.FindInPath(plugin, paths)
			if err == nil {
				return conf, "", fmt.Errorf("failed to find %s CNI: %s", plugin, err)
			}

			return conf, path, nil
		}
	}

	return nil, "", fmt.Errorf("failed to find %s CNI config: %v", plugin, err)
}

func cmdAdd(args *skel.CmdArgs) error {
	conf, err := getNetConf(args)
	if err != nil {
		return err
	}

	netconf, path, err := findCNI(conf.ConfDir, conf.Type)
	if err != nil {
		return err
	}

	pluginArgs := &invoke.Args{
		Command:     "ADD",
		ContainerID: args.ContainerID,
		NetNS:       args.Netns,
		IfName:      args.IfName,
		Path:        args.Path,
	}

	res, err := invoke.ExecPluginWithResult(context.TODO(), path, netconf.Bytes, pluginArgs, nil)
	if err != nil {
		return fmt.Errorf("%s plugin failed: %s", conf.Type, err)
	}

	_, err = current.NewResultFromResult(res)
	if err != nil {
		return err
	}

	return types.PrintResult(res, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	conf, err := getNetConf(args)
	if err != nil {
		return err
	}

	return nil
}

func cmdCheck(args *skel.CmdArgs) error {
	_, err := getNetConf(args)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, buildversion.BuildString("vf"))
}
