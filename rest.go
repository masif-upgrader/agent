package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/masif-upgrader/common"
	log "github.com/sirupsen/logrus"
	"net"
)

func startRestServer(sock string) (*iris.Application, error) {
	log.WithFields(log.Fields{"socket": sock}).Info("Listening on *nix socket")

	server, errLs := net.Listen("unix", sock)
	if errLs != nil {
		return nil, errLs
	}

	app := iris.New()

	app.Get("/v1/load", getV1Load)

	go app.Run(iris.Listener(server), iris.WithoutStartupLog)

	return app, nil
}

func getV1Load(context context.Context) {
	context.JSON(&struct {
		Install   [3]float64 `json:"install"`
		Update    [3]float64 `json:"update"`
		Configure [3]float64 `json:"configure"`
		Remove    [3]float64 `json:"remove"`
		Purge     [3]float64 `json:"purge"`
		Error     [3]float64 `json:"error"`
	}{
		actionsStats[common.PkgMgrInstall].queryLoad(),
		actionsStats[common.PkgMgrUpdate].queryLoad(),
		actionsStats[common.PkgMgrConfigure].queryLoad(),
		actionsStats[common.PkgMgrRemove].queryLoad(),
		actionsStats[common.PkgMgrPurge].queryLoad(),
		errorStats.queryLoad(),
	})
}
