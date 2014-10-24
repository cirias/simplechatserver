package main

import (
  "fmt"
  "net"
  "time"
  "strings"
  "strconv"
  )

const layout = "15:04:05"

type Server struct {
  clientMap map[Client]string
  listener *net.TCPListener
}

func NewServer(port string) (Server, error) {
  var s Server
  addr, err := net.ResolveTCPAddr("tcp", ":" + port)
  if err != nil {
    return s, err
  }
  server, err := net.ListenTCP("tcp", addr)
  if err != nil {
    return s, err
  }
  s = Server{make(map[Client]string), server}
  fmt.Println("Server listen on", port)
  return s, nil
}

func (s *Server) Close() {
  fmt.Println("Server stoping")
  s.listener.Close()
}

func (s *Server) Listen() {
  for {
    c, err := CreateClient(s.listener)
    if err != nil {
      fmt.Println(err)
      c.In <- "\nClose connection: " + err.Error()
      c.Close()
      continue
    }
    defer c.Close()
    names := []string {}
    for key := range s.clientMap {
      names = append(names, key.name)
    }
    message := "\n" + strconv.Itoa(len(s.clientMap)) + " clients online:\n"
    message += strings.Join(names, "\n")
    message += "\nEnjoy chating!\n\n> "
    c.In <- message

    s.BoardCast("\r" + time.Now().Format(layout) + " " + c.name + " entered the room.\n> ")
    s.clientMap[c] = c.conn.RemoteAddr().String()

    go s.ListenClient(c)
  }
}

func (s *Server) BoardCast(message string) {
  for key := range s.clientMap {
    key.In <- message
  }
}

func (s *Server) ListenClient(c Client) {
  for {
    select {
      case err := <-c.Err:
        fmt.Println(err, "from", c.name)
        delete(s.clientMap, c)
        s.BoardCast("\r" + time.Now().Format(layout) + " " + c.name + " left the room.\n> ")
        return
      case out := <-c.Out:
        fmt.Print("Message from " + c.name + ": ", out)
        for key := range s.clientMap {
          if key != c {
            key.In <- "\r" + time.Now().Format(layout) + " " + c.name + "> " + out + "> "
          } else {
            if (len(s.clientMap) == 1) {
              warn := "Poor child. Here's only you on the server.\n"
              warn += "You can call some friends.\n"
              warn += "Or just open another terminal and talk to yourself =)\n"
              warn += "> "
              key.In <- warn
            } else {
              key.In <- "> "
            }
          }
        }
    }
  }
}
