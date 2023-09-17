package main

func inArray[K comparable](arr []K, el K) int {
	for i, c := range arr {
		if el == c {
			return i
		}
	}
	return -1
}

func arrayRemove[K comparable](arr []K, i int) []K {
	if i == -1 {
		return arr
	} else if i == len(arr)-1 {
		return arr[:len(arr)-1]
	} else {
		copy(arr[i:], arr[i+1:])
		arr[len(arr)-1] = arr[0]
		arr = arr[:len(arr)-1]
		return arr
	}
}