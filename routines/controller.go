package routines

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// struct to help control routines.
// will wait till all started routines finish to call `Wait`
// Might could have use sync error group but this will wait
// till all finish, rather than on first error
type Controller struct {
	//----running------
	mainCountChan chan bool // true up false down
	backCountChan chan bool // true up false down
	//----listeners------
	doneChan   chan struct{} // used for `Done` (notify other of gracefully close)
	finishChan chan struct{} // used for `Wait` (notify main of finished)
	//-----Pass errors-----
	errorChan chan error
	errors    []string
}

type Runner interface {
	Run(rc *Controller) error
}

// Controller can run two types of jobs:
// - Runner: When these finish `Done()` will be called
// - BackgroundRunner: These should listen to `IsDone()` and gracefully exit
func NewController() *Controller {
	c := &Controller{
		mainCountChan: make(chan bool),
		backCountChan: make(chan bool),
		doneChan:      make(chan struct{}),
		finishChan:    make(chan struct{}),
		errorChan:     make(chan error),
		errors:        make([]string, 0),
	}

	go c.listenForCtrlC()
	go c.runMain()
	go c.runErr()

	return c
}

// ----------------Handle close-----------------
func (c *Controller) listenForCtrlC() {
	ch := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	log.Info("Gracefully shutting down...")
	c.Done()

	<-ch
	log.Info("Killing!")
	c.Finish()
}

// Call when you want to gracefully close everything
func (c *Controller) Done() {
	log.Debug("Done")
	select {
	case <-c.doneChan:
		log.Debug("Already shuting down...")
	default:
		close(c.doneChan)
	}
}
func (c *Controller) IsDone() chan struct{} {
	return c.doneChan
}

// Call when you want to exit
func (c *Controller) Finish() {
	log.Debug("Finish")
	select {
	case <-c.finishChan:
		log.Debug("Already finished...")
	default:
		close(c.finishChan)
	}
}
func (c *Controller) IsFinished() error {
	c.mainCountChan <- false // so that we dont close until something is waiting
	<-c.finishChan
	if len(c.errors) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(c.errors, "\n"))
}

//----------------Handle close-----------------

func (c *Controller) runMain() {
	mainCount := 1 // start with 1 so dont try closing until wait is called
	backgroundCount := 0

	for {
		select {
		case mc := <-c.mainCountChan:
			if mc {
				mainCount += 1
			} else {
				mainCount -= 1
			}

			if mainCount == 0 {
				c.Done()
			}
		case bc := <-c.backCountChan:
			if bc {
				backgroundCount += 1
			} else {
				backgroundCount -= 1
			}
		}

		log.Debugf("Main %v, Background %v", mainCount, backgroundCount)
		if mainCount+backgroundCount == 0 {
			select {
			case <-c.doneChan:
				log.Debug("Closing Error Chan")
				close(c.errorChan)
			default:
				log.Debug("Not done?")
			}
		}
	}
}

func (c *Controller) runErr() {
	for {
		newError, ok := <-c.errorChan
		if !ok {
			break
		}
		c.errors = append(c.errors, newError.Error())
	}
	c.Finish()
}

// will run all of these until none left
// if the runner returns error it will add to chan and print at end
func (c *Controller) Go(runner Runner) {
	select {
	case <-c.doneChan:
		log.Debug("Not running job because shuting down")
	default:
		c.mainCountChan <- true
		go func() {
			err := runner.Run(c)
			if err != nil {
				c.errorChan <- err
			}
			c.mainCountChan <- false
		}()
	}

}

// Running in background will get stopped when all `Go` created ones finish
// if the bgRunner returns error it will gracefully shutdown everything else
func (c *Controller) GoBackground(bgRunner Runner) {
	select {
	case <-c.doneChan:
		log.Debug("Not running background job because shuting down")
	default:
		c.backCountChan <- true
		go func() {
			err := bgRunner.Run(c)
			if err != nil {
				c.errorChan <- err
				c.Done()
			}
			c.backCountChan <- false
		}()
	}
}
