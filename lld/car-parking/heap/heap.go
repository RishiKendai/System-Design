package heap

import "github.com/lld/car-parking/slot"

type SpotMinHeap []*slot.ParkingSlot

func (h SpotMinHeap) Less(i, j int) bool {
	if h[i].FloorNumber == h[j].FloorNumber {
		return h[i].Id < h[j].Id
	}
	return h[i].FloorNumber < h[j].FloorNumber
}

func (h *SpotMinHeap) Push(x any) {
	*h = append(*h, x.(*slot.ParkingSlot))
}

func (h *SpotMinHeap) Pop() any {
	old := *h
	n := len(old)
	spot := old[n-1]
	*h = old[:n-1]
	return spot
}

func (h SpotMinHeap) Len() int {
	return len(h)
}

func (h SpotMinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
