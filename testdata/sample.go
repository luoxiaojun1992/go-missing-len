package main

func main() {
	s1 := []int{1, 2, 3}
	s := make([]int, 0)
	for _, num := range s1 {
		s = append(s, num)
	}

	s2 := []string{"foo", "bar"}
	m := make(map[string]string, 0)
	for _, str := range s2 {
		m[str] = str
	}
}
