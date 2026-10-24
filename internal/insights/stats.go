// Package insights covers the survey (scope 8.x), visit tracking and admin reports (scope 9.x).
package insights

import "math"

// MeanSD returns the mean and sample standard deviation (n-1), as reported in Thai research practice.
func MeanSD(xs []float64) (mean, sd float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	for _, x := range xs {
		mean += x
	}
	mean /= float64(len(xs))
	if len(xs) < 2 {
		return mean, 0
	}
	var ss float64
	for _, x := range xs {
		ss += (x - mean) * (x - mean)
	}
	return mean, math.Sqrt(ss / float64(len(xs)-1))
}

// Level interprets a 5-point Likert mean with the common Best (1977) intervals.
func Level(mean float64) string {
	switch {
	case mean >= 4.51:
		return "highest"
	case mean >= 3.51:
		return "high"
	case mean >= 2.51:
		return "moderate"
	case mean >= 1.51:
		return "low"
	}
	return "lowest"
}
