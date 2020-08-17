package utility

func MinAndMax(num ...int) (min, max int) {
	if len(num) == 0 {
		panic("stat MinAndMax from empty")
	}

	min, max = num[0], num[0]
	for i := 1; i < len(num); i++ {
		if num[i] > max {
			max = num[i]
		}
		if num[i] < min {
			min = num[i]
		}
	}

	return
}
