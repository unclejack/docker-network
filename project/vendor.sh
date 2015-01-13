#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."

# Downloads dependencies into vendor/ directory
mkdir -p vendor
cd vendor

clone() {
	vcs=$1
	pkg=$2
	rev=$3
	
	pkg_url=https://$pkg
	target_dir=src/$pkg
	
	echo -n "$pkg @ $rev: "
	
	if [ -d $target_dir ]; then
		echo -n 'rm old, '
		rm -fr $target_dir
	fi
	
	echo -n 'clone, '
	case $vcs in
		git)
			git clone --quiet --no-checkout $pkg_url $target_dir
			( cd $target_dir && git reset --quiet --hard $rev )
			;;
		hg)
			hg clone --quiet --updaterev $rev $pkg_url $target_dir
			;;
	esac
	
	echo -n 'rm VCS, '
	( cd $target_dir && rm -rf .{git,hg} )
	
	echo done
}

clone git github.com/codegangsta/cli
clone git go.googlesource.com/sys
clone git github.com/vishvananda/netlink
clone git github.com/Sirupsen/logrus v0.6.0
if [ "$1" == '--docker' ]; then
	clone git github.com/docker/docker 0da9540eda55caaa8ed486b6814afc253a0e181d
fi
clone git github.com/docker/libpack
clone git github.com/docker/libcontainer 6460fd79667466d2d9ec03f77f319a241c58d40b
clone git github.com/erikh/ping
clone git github.com/syndtr/gocapability

mkdir -p src/golang.org/x/
mv src/go.googlesource.com/sys src/golang.org/x/
rm -rf src/go.googlesource.com

