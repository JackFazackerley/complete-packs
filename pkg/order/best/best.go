package best

import (
	"container/list"
	"math"

	"github.com/JackFazackerley/complete-packs/pkg/pack"
)

// Calculate simply calculates the best sequence of packs for the given target of items
func Calculate(target int, sizes []float64) (sequence []pack.Pack) {
	// queue is used for keeping a queue of the next number, will only ever be positive numbers
	queue := list.New()

	// positive is used for positive positions
	pos := make([][]pack.Pack, target+1)

	// offset is used to find the closest result to 0
	offset := math.MinInt64

	//sequence is used to store the optimum sequence of pack sizes
	sequence = make([]pack.Pack, 0)

	// we need to fill the first positive position so we know where to start from
	pos[target] = make([]pack.Pack, len(sizes))
	for i, size := range sizes {
		pos[target][i] = pack.Pack{
			Size:  size,
			Count: 0,
		}
	}

	// push the target as the initial positive number
	queue.PushBack(target)

	// we want to continue until the queue is empty
	for queue.Len() > 0 {
		// get the front of the queue
		e := queue.Front()
		num := e.Value.(int)

		for index, size := range sizes {

			result := num - int(size)
			// if the result is greater than 0 then we know we can get it lower
			if result > 0 {
				// if the positive number already exists then we don't need to compute again
				if pos[result] == nil {
					// push the result onto the queue as this is the first time we've seen it
					queue.PushBack(result)
					pos[result] = make([]pack.Pack, len(sizes))

					// pos[num] might be empty but if it's not then we can loop over and add the values to the current
					// result
					for k, v := range pos[num] {
						pos[result][k].Count += v.Count
						pos[result][k].Size = v.Size
					}

					// keep track of the current size and add to it
					pos[result][index].Count += 1
				}
			} else {
				// we know the result is now negative so we need to check to see if the result is closer to 0 than the
				// previous
				if result > offset {
					offset = result
					// set sequence equal to the current num
					sequence = pos[num]

					// keep track of the current size and add to it
					sequence[index].Count += 1
				}
			}
		}
		// we're done with this number remove it from the queue
		queue.Remove(e)
	}

	return sequence
}
