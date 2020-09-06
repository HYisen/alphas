package utility

import (
	"io"
	"log"
	"time"
)

func MustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

// ExecLater execute action function later, the first input from hint would setup delay time and start the timer,
// moreover time.Duration from hint would reset the timer by its value, only if the action has not been executed.
func ExecLater(action func(), hint <-chan time.Duration) {
	delay := <-hint
	ticker := time.NewTicker(delay)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			action()
			return
		case delay = <-hint:
			ticker.Reset(delay)
		}
	}
}

// SummonDestroyer would close target after delay,
// preAction would be run just before the close.
// Everytime signal from respawn would reset the timer to close.
func SummonDestroyer(target io.Closer, delay time.Duration, preAction func(), respawn <-chan struct{}) {
	hint := make(chan time.Duration)
	go ExecLater(func() {
		preAction()
		CloseAndLogError(target)
	}, hint)
	hint <- delay
	for range respawn {
		hint <- delay
	}
}
