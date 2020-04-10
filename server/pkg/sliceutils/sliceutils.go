package sliceutils

type ArrayElement interface{}

func IndexOf(element *ArrayElement, data []*ArrayElement) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
