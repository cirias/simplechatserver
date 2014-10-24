package main

import (
  "fmt"
  "net"
  "time"
  "bufio"
  "errors"
  )

type Client struct {
  name string
  conn net.Conn
  In chan string
  Out chan string
  Err chan error
}

func CreateClient(l *net.TCPListener) (Client, error) {
  var c Client
  conn, err := l.Accept()
  if err != nil {
    return c, err
  }
  fmt.Println("Connection accepted")
  c = Client{"", conn, make(chan string), make(chan string), make(chan error)}
  go c.receiver()
  go c.sender()

  // Block accept other connection
  welcome := "***************************************\n"
  welcome += "*                                     *\n"
  welcome += "*         Welcome to chat room        *\n"
  welcome += "*                                     *\n"
  welcome += "***************************************\n"
  welcome += "\n"
  welcome += "Please input your name: "
  c.In <- welcome
  select {
    case name := <-c.Out:
      c.name = name[:len(name) - 2]
      return c, nil
    case err := <- c.Err:
      return c, errors.New(c.conn.RemoteAddr().String() + " " + err.Error())
    case <-time.After(20 * time.Second):
      return c, errors.New(c.conn.RemoteAddr().String() + " timeout")
  }
}

func (c *Client) Close() {
  c.conn.Close()
}

func (c *Client) receiver() {
  reader := bufio.NewReader(c.conn)
  for {
    data, err := reader.ReadString('\n')
    if err != nil {
      c.Err <- err
    } else {
      c.Out <- data
    }
  }
}

func (c *Client) sender() {
  for data := range c.In {
    _, err := c.conn.Write([]byte(data))
    if err != nil {
      c.Err <- err
    }
  }
}
