//go:generate go run vendor/github.com/Al2Klimov/go-gen-source-repos/main.go github.com/masif-upgrader/agent

package main

import (
	"errors"
	"flag"
	"fmt"
	_ "github.com/Al2Klimov/go-gen-source-repos"
	pp "github.com/Al2Klimov/go-pretty-print"
	"github.com/go-ini/ini"
	"github.com/masif-upgrader/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
	"time"
)

type settings struct {
	interval struct {
		check, report, retry int64
	}
	master struct {
		host string
	}
	tls struct {
		cert, key, ca string
	}
	log struct {
		level log.Level
	}
}

var zeroTime = time.Duration(0)
var retryInterval time.Duration
var logLevels = map[string]log.Level{
	"error":   log.ErrorLevel,
	"err":     log.ErrorLevel,
	"warning": log.WarnLevel,
	"warn":    log.WarnLevel,
	"info":    log.InfoLevel,
	"debug":   log.DebugLevel,
}

func main() {
	if len(os.Args) == 1 && terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf(
			"For the terms of use, the source code and the authors\n"+
				"see the projects this program is assembled from:\n\n  %s\n",
			strings.Join(GithubcomAl2klimovGo_gen_source_repos, "\n  "),
		)
		os.Exit(1)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if err := runAgent(); err != nil {
		log.Fatal(err)
	}
}

func runAgent() error {
	cfg, errLC := loadCfg()
	if errLC != nil {
		return errLC
	}

	log.SetLevel(cfg.log.level)

	sigListener := &signalListener{}
	sigListener.onSignals(func(sig os.Signal) {
		log.WithFields(log.Fields{"signal": lazyLogString{sig.String}}).Warn("Caught signal, exiting")
		os.Exit(0)
	}, syscall.SIGTERM, syscall.SIGINT)

	master, errNA := newApi(cfg.master.host, cfg.tls)
	if errNA != nil {
		return errNA
	}

	log.Debug("Auto-detecting package manager")

	ourPkgMgr, errPM := newApt()
	if errPM != nil {
		return errPM
	}

	if ourPkgMgr == nil {
		return errors.New("package manager not available or not supported")
	}

	log.WithFields(log.Fields{"package_manager": ourPkgMgr.getName()}).Info("Auto-detected package manager")

	var tasks map[common.PkgMgrTask]struct{} = nil
	retryInterval = time.Duration(cfg.interval.retry) * time.Second
	approvedTasks := map[common.PkgMgrTask]struct{}{}
	ctxtWIU := "querying package manager"
	interval := struct{ check, report time.Duration }{
		check:  time.Duration(cfg.interval.check) * time.Second,
		report: time.Duration(cfg.interval.report) * time.Second,
	}
	whatIfUpgradeAll := func() (err error) {
		tasks, err = ourPkgMgr.whatIfUpgradeAll(sigListener)
		return
	}

	for {
		if tasks == nil {
			if errWIUA := retryOp(whatIfUpgradeAll, ctxtWIU); errWIUA != nil {
				return errWIUA
			}
		}

		if len(tasks) > 0 {
			notApprovedTasks := map[common.PkgMgrTask]struct{}{}

			for task := range tasks {
				if _, isApproved := approvedTasks[task]; !isApproved {
					notApprovedTasks[task] = struct{}{}
				}
			}

			log.WithFields(log.Fields{"pending": len(notApprovedTasks), "total": len(tasks)}).Info("Something to do")

			if len(notApprovedTasks) > 0 {
				var freshApprovedTasks map[common.PkgMgrTask]struct{}
				tasks = nil

				errRT := retryOp(func() (err error) {
					freshApprovedTasks, err = master.reportTasks(notApprovedTasks)
					return
				}, "reporting to master")
				if errRT != nil {
					return errRT
				}

				log.WithFields(log.Fields{"approved": len(freshApprovedTasks)}).Info("Got new approvals")

				for task := range freshApprovedTasks {
					approvedTasks[task] = struct{}{}
				}
			}

			for {
				if tasks == nil {
					if errWIUA := retryOp(whatIfUpgradeAll, ctxtWIU); errWIUA != nil {
						return errWIUA
					}
				}

				nextPackage := ""
				actionsNeeded := ^uint64(0)

			PossibleActions:
				for task := range tasks {
					if _, isApproved := approvedTasks[task]; isApproved && task.Action == common.PkgMgrUpdate {
						tasks = nil

						var tasksOnUpgrade map[common.PkgMgrTask]struct{}

						errWIU := retryOp(func() (err error) {
							tasksOnUpgrade, err = ourPkgMgr.whatIfUpgrade(sigListener, task.PackageName)
							return
						}, ctxtWIU)
						if errWIU != nil {
							return errWIU
						}

						for taskOnUpgrade := range tasksOnUpgrade {
							if _, approved := approvedTasks[taskOnUpgrade]; !approved {
								log.WithFields(log.Fields{
									"package": task.PackageName,
									"dependency": map[string]interface{}{
										"name":   taskOnUpgrade.PackageName,
										"action": taskOnUpgrade.Action,
									},
								}).Debug("Package can't be upgraded as not all required actions have been approved")

								continue PossibleActions
							}
						}

						actionsNeededForUpgrade := uint64(len(tasksOnUpgrade))
						if actionsNeededForUpgrade < actionsNeeded {
							actionsNeeded = actionsNeededForUpgrade
							nextPackage = task.PackageName
						}
					}
				}

				if nextPackage == "" {
					break
				}

				tasks = nil

				log.WithFields(log.Fields{"package": nextPackage}).Info("Upgrading")

				if errU := ourPkgMgr.upgrade(sigListener, nextPackage); errU != nil {
					if retryInterval == zeroTime {
						return errU
					}

					log.WithFields(log.Fields{
						"package": nextPackage,
						"error":   lazyLogString{errU.Error},
					}).Error("Upgrade failed")

					sleep(retryInterval)
				}
			}

			if tasks == nil {
				if errWIUA := retryOp(whatIfUpgradeAll, ctxtWIU); errWIUA != nil {
					return errWIUA
				}
			}

			if len(tasks) > 0 {
				tasks = nil
				sleep(interval.report)
			} else {
				approvedTasks = map[common.PkgMgrTask]struct{}{}
			}
		} else {
			log.Info("Nothing to do")

			approvedTasks = map[common.PkgMgrTask]struct{}{}
			tasks = nil
			sleep(interval.check)
		}
	}
}

func loadCfg() (config *settings, err error) {
	cfgFile := flag.String("config", "", "config file")
	flag.Parse()

	if *cfgFile == "" {
		return nil, errors.New("config file missing")
	}

	log.WithFields(log.Fields{"file": *cfgFile}).Debug("Loading config")

	cfg, errLI := ini.Load(*cfgFile)
	if errLI != nil {
		return nil, errLI
	}

	cfgInterval := cfg.Section("interval")
	cfgTls := cfg.Section("tls")
	result := &settings{
		interval: struct{ check, report, retry int64 }{
			check:  cfgInterval.Key("check").MustInt64(),
			report: cfgInterval.Key("report").MustInt64(),
			retry:  cfgInterval.Key("retry").MustInt64(),
		},
		master: struct{ host string }{
			host: cfg.Section("master").Key("host").String(),
		},
		tls: struct{ cert, key, ca string }{
			cert: cfgTls.Key("cert").String(),
			key:  cfgTls.Key("key").String(),
			ca:   cfgTls.Key("ca").String(),
		},
	}

	if result.interval.check <= 0 {
		return nil, errors.New("config: interval.check missing")
	}

	if result.interval.report <= 0 {
		return nil, errors.New("config: interval.report missing")
	}

	if result.interval.retry <= 0 {
		result.interval.retry = 0
	}

	if result.master.host == "" {
		return nil, errors.New("config: master.host missing")
	}

	if result.tls.cert == "" {
		return nil, errors.New("config: tls.cert missing")
	}

	if result.tls.key == "" {
		return nil, errors.New("config: tls.key missing")
	}

	if result.tls.ca == "" {
		return nil, errors.New("config: tls.ca missing")
	}

	if rawLogLvl := cfg.Section("log").Key("level").String(); rawLogLvl == "" {
		result.log.level = log.InfoLevel
	} else if logLvl, logLvlValid := logLevels[rawLogLvl]; logLvlValid {
		result.log.level = logLvl
	} else {
		return nil, errors.New("config: bad log.level")
	}

	return result, nil
}

func retryOp(op func() error, desc string) (err error) {
	for try := uint64(1); ; try++ {
		log.WithFields(log.Fields{
			"operation": desc,
			"try":       try,
		}).Info("Trying")

		if err = op(); err == nil || retryInterval == zeroTime {
			if err == nil && try > 1 {
				log.WithFields(log.Fields{
					"operation": desc,
					"try":       try,
				}).Info("Recovered")
			}

			return
		}

		log.WithFields(log.Fields{
			"operation": desc,
			"try":       try,
			"error":     lazyLogString{err.Error},
		}).Error("Failed")

		sleep(retryInterval)
	}
}

func sleep(dur time.Duration) {
	log.WithFields(log.Fields{"duration": pp.Duration(dur)}).Debug("Sleeping")
	time.Sleep(dur)
}
