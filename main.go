package main

import (
  "net"
  "bufio"
  )

var sockets = make([]net.Conn, 0)

func main() {
  ln, err := net.Listen("tcp", ":25252")
  if err != nil {
    return
  }
  defer ln.Close()

  for {
    conn, err := ln.Accept()
    if err != nil {
      continue
    }
    defer conn.Close()

    sockets = append(sockets, conn)
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  reader := bufio.NewReader(conn)
  for {
    input, err := reader.ReadString('\n')
    if err != nil {
      for i, v := range sockets {
        if v == conn {
          sockets = append(sockets[:i], sockets[i + 1:]...)
        }
      }
      continue
    }
    for _, v := range sockets {
      if v == conn {
        continue
      }
      v.Write([]byte(input))
    }
  }
}
