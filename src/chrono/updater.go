/*

MIT License

Copyright (c) 2019 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package chrono

import (
	"fmt"
	"runtime"
	"time"
)

type Updater struct {
	onUpdate        func() bool
	updateFrequency int
	quit            chan bool
}

func NewUpdater(updateFrequency int, onUpdate func() bool) (*Updater, error) {
	if updateFrequency == 0 {
		return nil, fmt.Errorf("illegal update frequency: %d", updateFrequency)
	}
	u := &Updater{
		onUpdate:        onUpdate,
		updateFrequency: updateFrequency,
		quit:            make(chan bool),
	}
	return u, nil
}

func (u *Updater) startUpdater() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	waitTimeInMs := 1000 / time.Duration(u.updateFrequency)
	tickInterval := time.Millisecond * waitTimeInMs

	for {
		time.Sleep(tickInterval)
		wantContinue := u.onUpdate()
		if !wantContinue {
			return
		}
	}

	/*
		ticker := time.NewTicker(tickInterval)

		for {
			select {
			case <-ticker.C:
				wantContinue := u.onUpdate()
				if !wantContinue {
					ticker.Stop()
					return
				}

			case <-u.quit:
				ticker.Stop()
			}
		}
	*/
}

func (u *Updater) Frequency() int {
	return u.updateFrequency
}

func (u *Updater) SetFrequency(frequency int) {
	u.updateFrequency = frequency
	u.Restart()
}

func (u *Updater) SetOnUpdate(onUpdate func() bool) {
	u.onUpdate = onUpdate
}

func (u *Updater) Start() {
	go u.startUpdater()
}

func (u *Updater) Stop() {
	u.quit <- true
}

func (u *Updater) Restart() {
	u.Stop()
	u.Start()
}
