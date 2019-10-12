package main

import "github.com/masif-upgrader/common"

var actionsStats = map[common.PkgMgrAction]*statsBookkeeper{
	common.PkgMgrInstall:   {},
	common.PkgMgrUpdate:    {},
	common.PkgMgrConfigure: {},
	common.PkgMgrRemove:    {},
	common.PkgMgrPurge:     {},
}

var queryStats, errorStats statsBookkeeper
