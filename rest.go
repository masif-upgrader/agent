package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	v1 "github.com/masif-upgrader/agent/v1"
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
	context.JSON(&v1.Load{
		queryStats.queryLoad(),
		actionsStats[common.PkgMgrInstall].queryLoad(),
		actionsStats[common.PkgMgrUpdate].queryLoad(),
		actionsStats[common.PkgMgrConfigure].queryLoad(),
		actionsStats[common.PkgMgrRemove].queryLoad(),
		actionsStats[common.PkgMgrPurge].queryLoad(),
		errorStats.queryLoad(),
	})
}
