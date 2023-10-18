# Xcluster/ovl - conman

Test environment for vf-operator using [xcluster](https://github.com/Nordix/xcluster)

## Known limitations
 * xcluster 8.0.0 doesn't have hd.img, use 7.9.0 instead
 * podman ovl is cached in this test directory, later xcluster versions will have this already included

## Start test environment
```
./conman.sh test --xterm
```

## Test podman commands
```
podman network ls
podman run --name testvf0 --rm -dt --net podman --net eth2vf0 library/alpine:latest
podman run --name testvf1 --rm -dt --net podman --net eth2vf1 library/alpine:latest
podman run --name testvf2 --rm -dt --net podman --net eth2vf2 library/alpine:latest
podman run --name testvf3 --rm -dt --net podman --net eth2vf3 library/alpine:latest

podman run --name testcdi --rm -dt --net podman --device vfoperator/vhost library/alpine:latest
podman run --name testvf --rm -dt --net podman --net llscu:interface_name=llscu --net f1-u:interface_name=f1-u --device vfoperator/vhost library/alpine:latest
```
