#! /bin/sh
##
## conman.sh --
##
##   Help script for the xcluster ovl/conman.
##

prg=$(basename $0)
dir=$(dirname $0); dir=$(readlink -f $dir)
me=$dir/$prg
tmp=/tmp/${prg}_$$
XCDIR=$dir/xcluster
test -n "$PREFIX" || PREFIX=1000::1

die() {
    echo "ERROR: $*" >&2
    rm -rf $tmp
    exit 1
}
help() {
    grep '^##' $0 | cut -c3-
    rm -rf $tmp
    exit 0
}
test -n "$1" || help
echo "$1" | grep -qi "^help\|-h" && help

log() {
	echo "$prg: $*" >&2
}
dbg() {
	test -n "$__verbose" && echo "$prg: $*" >&2
}

install() {
	test -d $XCDIR && return

        ver=8.0.0
	tmp=$(mktemp -d .xcluster-XXXX)
	file=xcluster-$ver.tar.xz
	url=https://github.com/Nordix/xcluster/releases/download/$ver/$file
	log "Downloading $url > $file"
        curl -o $tmp/$file -L $url > /dev/null 2>&1
	tar -xvf $tmp/$file > /dev/null 2>&1
	rm -r $tmp

	# Copy all xcluster caches
	log "Copying local xcluster cache"
	cp cache/* $XCDIR/workspace/xcluster/cache/default
}

## Commands;
##

##   env
##     Print environment.
cmd_env() {

	if test "$cmd" = "env"; then
		set | grep -E '^(__.*)='
		return 0
	fi

	test -n "$xcluster_DOMAIN" || xcluster_DOMAIN=xcluster
	test -n "$XCLUSTER" || die 'Not set [$XCLUSTER]'
	test -x "$XCLUSTER" || die "Not executable [$XCLUSTER]"
	eval $($XCLUSTER env)
}

##
## Tests;
##   test [--xterm] [--no-stop] test <test-name>  [ovls...] > $log
##   test [--xterm] [--no-stop] > $log   # default test
##     Exec tests
##
cmd_test() {
	cmd_env
	start=starts
	test "$__xterm" = "yes" && start=start
	rm -f $XCLUSTER_TMP/cdrom.iso

	if test -n "$1"; then
		local t=$1
		shift
		test_$t $@
	else
		test_start
	fi

	now=$(date +%s)
	tlog "Xcluster test ended. Total time $((now-begin)) sec"
}
##   test start_empty
##     Start empty cluster
test_start_empty() {
	cd $dir
	test -n "$__nrouters" || export __nrouters=1
	export xcluster_PREFIX=$PREFIX
	if test -n "$TOPOLOGY"; then
		tlog "Using TOPOLOGY=$TOPOLOGY"
		. $($XCLUSTER ovld network-topology)/$TOPOLOGY/Envsettings
	fi
	xcluster_start network-topology iptools lspci podman . $@
}
##   test start
##     Start cluster with ovl functions
test_start() {
	export __mem=1024
	export __nvm=1
	export __nets_vm=0,1,2,3,4
	export __net_setup=$(dirname $XCLUSTER)/config/net-setup-pci-emulator.sh
	export __kvm=$GOPATH/src/github.com/qemu/qemu/build/qemu-system-x86_64
	export __machine=q35
	test -n "$__nrouters" || export __nrouters=0

	test_start_empty $@
	otcw "vf eth2 4"

	vm=1
	port=5001
	gwPort=15001
	log "Kill all existing socat sessions"
	pkill socat
	log "Start port forwarding to vm-$vm:$port"
	socat tcp-l:$port,fork,reuseaddr tcp:192.168.0.$vm:$port &
	log "Start port forwarding to $vm:$gwPort"
	socat tcp-l:$gwPort,fork,reuseaddr tcp:192.168.0.$vm:$gwPort &
}
##   test basic (default)
##     Just start and stop
test_basic() {
	test_start $@
	xcluster_stop
}

##
install
cd $XCDIR
. ./Envsettings
export XCLUSTER_OVLPATH=$XCDIR/ovl:$dir
. $($XCLUSTER ovld test)/default/usr/lib/xctest
cd - > /dev/null
indent=''

# Get the command
cmd=$1
shift
grep -q "^cmd_$cmd()" $0 $hook || die "Invalid command [$cmd]"

while echo "$1" | grep -q '^--'; do
	if echo $1 | grep -q =; then
		o=$(echo "$1" | cut -d= -f1 | sed -e 's,-,_,g')
		v=$(echo "$1" | cut -d= -f2-)
		eval "$o=\"$v\""
	else
		o=$(echo "$1" | sed -e 's,-,_,g')
		eval "$o=yes"
	fi
	shift
done
unset o v
long_opts=`set | grep '^__' | cut -d= -f1`

# Execute command
trap "die Interrupted" INT TERM
cmd_$cmd "$@"
status=$?
rm -rf $tmp
exit $status
