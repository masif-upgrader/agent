package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
)

type criticalOperationRunner interface {
	runCritical(op func())
}

type signalListener struct {
	critOpMutex sync.RWMutex
}

func (self *signalListener) runCritical(op func()) {
	self.critOpMutex.RLock()
	op()
	self.critOpMutex.RUnlock()
}

type signals struct {
	signals []os.Signal
}

func (s signals) Format(f fmt.State, c rune) {
	result := make([]interface{}, len(s.signals))

	for i, sig := range s.signals {
		result[i] = sig.String()
	}

	fmt.Fprint(f, result)
}

func (s signals) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, len(s.signals))

	for i, sig := range s.signals {
		result[i] = sig.String()
	}

	return json.Marshal(result)
}

func (self *signalListener) onSignals(handler func(sig os.Signal), sigs ...os.Signal) {
	log.WithFields(log.Fields{"signals": signals{sigs}}).Debug("Listening for signals")

	ch := make(chan os.Signal, 64)
	signal.Notify(ch, sigs...)

	go func() {
		for {
			if sig, hasSig := <-ch; hasSig {
				self.critOpMutex.Lock()
				handler(sig)
				self.critOpMutex.Unlock()
			} else {
				break
			}
		}
	}()
}
