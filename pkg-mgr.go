package main

import "github.com/masif-upgrader/common"

type pkgMgr interface {
	getName() string
	whatIfUpgradeAll(critOpRunner criticalOperationRunner) (tasks map[common.PkgMgrTask]struct{}, err error)
	whatIfUpgrade(critOpRunner criticalOperationRunner, packageName string) (tasks map[common.PkgMgrTask]struct{}, err error)
	upgrade(critOpRunner criticalOperationRunner, packageName string) error
}
