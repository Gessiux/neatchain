package consensus

import (
	"time"

	"github.com/Gessiux/neatchain/chain/log"

	. "github.com/Gessiux/go-common"
)

var (
	tickTockBufferSize = 10
)

// TimeoutTicker is a timer that schedules timeouts
// conditional on the height/round/step in the timeoutInfo.
// The timeoutInfo.Duration may be non-positive.
type TimeoutTicker interface {
	Start() (bool, error)
	Stop() bool
	Chan() <-chan timeoutInfo       // on which to receive a timeout
	ScheduleTimeout(ti timeoutInfo) // reset the timer
}

// timeoutTicker wraps time.Timer,
// scheduling timeouts only for greater height/round/step
// than what it's already seen.
// Timeouts are scheduled along the tickChan,
// and fired on the tockChan.
type timeoutTicker struct {
	BaseService

	timer    *time.Timer
	tickChan chan timeoutInfo
	tockChan chan timeoutInfo

	logger log.Logger
}

func NewTimeoutTicker(logger log.Logger) TimeoutTicker {
	tt := &timeoutTicker{
		timer:    time.NewTimer(0),
		tickChan: make(chan timeoutInfo, tickTockBufferSize),
		tockChan: make(chan timeoutInfo, tickTockBufferSize),
		logger:   logger,
	}
	tt.stopTimer() // don't want to fire until the first scheduled timeout
	tt.BaseService = *NewBaseService(logger, "TimeoutTicker", tt)
	return tt
}

func (t *timeoutTicker) OnStart() error {

	go t.timeoutRoutine()

	return nil
}

func (t *timeoutTicker) OnStop() {
	t.BaseService.OnStop()
	t.stopTimer()
}

func (t *timeoutTicker) Chan() <-chan timeoutInfo {
	return t.tockChan
}

// The timeoutRoutine is alwaya available to read from tickChan (it won't block).
// The scheduling may fail if the timeoutRoutine has already scheduled a timeout for a later height/round/step.
func (t *timeoutTicker) ScheduleTimeout(ti timeoutInfo) {
	t.tickChan <- ti
}

//-------------------------------------------------------------

// stop the timer and drain if necessary
func (t *timeoutTicker) stopTimer() {
	// Stop() returns false if it was already fired or was stopped
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
		default:
			t.logger.Debug("Timer already stopped")
		}
	}
}

// send on tickChan to start a new timer.
// timers are interupted and replaced by new ticks from later steps
// timeouts of 0 on the tickChan will be immediately relayed to the tockChan
func (t *timeoutTicker) timeoutRoutine() {
	t.logger.Info("Starting timeout routine")
	var ti timeoutInfo
	for {
		select {
		case newti := <-t.tickChan:
			t.logger.Infof("Received tick. old_ti: %v, new_ti: %v", ti, newti)

			// ignore tickers for old height/round/step
			/*
				if newti.Height < ti.Height {
					continue
				} else if newti.Height == ti.Height {
					if newti.Round < ti.Round {
						continue
					} else if newti.Round == ti.Round {
						if ti.Step > 0 && newti.Step <= ti.Step {
							continue
						}
					}
				}
			*/
			// stop the last timer
			t.stopTimer()

			// update timeoutInfo and reset timer
			// NOTE time.Timer allows duration to be non-positive
			ti = newti
			t.timer.Reset(ti.Duration)
			t.logger.Infof("Scheduled timeout. dur: %v, height: %v, round: %v, step: %v", ti.Duration, ti.Height, ti.Round, ti.Step)
		case <-t.timer.C:
			t.logger.Infof("Timed out. dur: %v, height: %v, round: %v, step: %v", ti.Duration, ti.Height, ti.Round, ti.Step)
			// go routine here gaurantees timeoutRoutine doesn't block.
			// Determinism comes from playback in the receiveRoutine.
			// We can eliminate it by merging the timeoutRoutine into receiveRoutine
			//  and managing the timeouts ourselves with a millisecond ticker
			go func(toi timeoutInfo) { t.tockChan <- toi }(ti)
		case <-t.Quit:
			return
		}
	}
}
