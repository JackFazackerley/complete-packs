package fast

import (
	"math"

	"github.com/JackFazackerley/complete-packs/pkg/pack"
)

// Calculate is used to find the best sequence of packs for the given target
func Calculate(target float64, sizes []float64) (sequence []pack.Pack) {
	// store the target as the remainder so we can pass the remainder to find the best pack
	remainder := target

	// we loop twice as the first iteration may find a solution but it might not be the best
	for j := 0; j < 2; j++ {
		// we use calculatedTotal to find the closest number to the target
		calculatedTotal := 0.
		sequence = make([]pack.Pack, 0, len(sizes)-1)

		// loop through the length of sizes
		for i := 0; i < len(sizes); i++ {
			// we pass everything after the current index so we don't do the same number twice
			newRemainder, pack := bestPack(sizes[i:], remainder)

			// store the new remainder
			remainder = newRemainder

			// store the pack to the sequence
			if pack.Count > 0 {
				sequence = append(sequence, pack)
			}

			//keep track of the total Calculated for the second iteration
			calculatedTotal += pack.Size * pack.Count
		}
		// remainder is now set to the total calculated on the next iteration we find the best solution
		remainder = calculatedTotal
	}

	//sequence = sequence[:len(sequence)-1]

	return
}

// bestPack is used to find the best single pack for the given target. We check that current target is greater than the
// next size so we know which size has the best count of packs for the target. Sizes must end with 0 because each index
// checks the next
func bestPack(sizes []float64, target float64) (remainder float64, cp pack.Pack) {
	bestDiv := 0.
	bestCount := math.MaxFloat64

	for i := 0; i < len(sizes)-1; i++ {
		div := target / sizes[i]

		// rather than using math.Floor for floating point precision we convert the float to
		// an int and then back to a float - done entirely for speed
		count := float64(int(div))

		// we only want to ceil the last index, excluding 0
		if len(sizes)-1 == 1 {
			count = math.Ceil(div)
		}

		// if the target is greater than the next size
		better := false
		if target > sizes[i+1] {
			better = true
		}

		// by default we get the difference by using mod. Math.Mod wasn't used because of float precision. If target is
		// better than next size then we want to get the max difference by multiplying size by count
		thisRemainder := float64(int(target) % int(sizes[i]))
		if better {
			thisRemainder = math.Max(0, target-(count*sizes[i]))
		}

		// check if div is greater than bestDiv. This is is for the amount of times target goes into size
		if div > bestDiv {
			// if the count is lesser than best count we know that we've got a greater amount of items for less parcels
			if count < bestCount {
				bestDiv = div
				bestCount = count
				remainder = thisRemainder
				cp = pack.Pack{
					Size:  sizes[i],
					Count: count,
				}
			}
		}
	}

	return
}
