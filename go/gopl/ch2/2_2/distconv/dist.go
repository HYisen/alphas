package distconv

import "fmt"

type Feet float64
type Meter float64

const oneFeetEqualsMeter = 3.2808

func FToM(f Feet) Meter {
	return Meter(f * oneFeetEqualsMeter)
}

func MToF(m Meter) Feet {
	return Feet(m / oneFeetEqualsMeter)
}

func (f Feet) String() string {
	return fmt.Sprintf("%.3f ft", f)
}

func (m Meter) String() string {
	return fmt.Sprintf("%.3f m", m)
}
