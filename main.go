package main

import (
  "fmt"
  )

func main() {
  s, err := NewServer("25252")
  if err != nil {
		fmt.Println(err)
	}
  defer s.Close()
  s.Listen()
}
