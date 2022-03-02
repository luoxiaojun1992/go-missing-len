package main

func main() {
	s := make([]int, 0)
	s1 := []int{1, 2, 3}
	for _, num := range s1 {
		s = append(s, num)
	}

	m := make(map[string]string, 0)
	m["foo"] = "bar"
}
