package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// sever open
	l, err := net.Listen("tcp", ":30122")
	if nil != err {
		fmt.Print("fial to bind address", err)
		return
	}
	defer l.Close()

	// loop server
	for {
		conn, err := l.Accept()
		if nil != err {
			log.Printf("fail to accept; err : %v", err)
			continue
		}
		go ConnHandler(conn)
	}
}

// ConnHandler input db receved data
func ConnHandler(conn net.Conn) {
	recvBuf := make([]byte, 256)
	ok := []byte{'o', 'k'}
	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Printf("fail to receive data; err: %v", err)
			return
		}
		if 0 < n {
			// pm10 pm25 hum tem time
			recv := recvBuf[:n]
			data := string(recv)
			fmt.Println(data)
			conn.Write(ok)
			datas := strings.Split(data, " ")
			Time := datas[4]+ " " + datas[5]
			_, err := db.Exec("insert into data(pm10,pm25,time) values (?,?,?)",datas[0], datas[1], Time)
			errHandler(err)
			_, err =db2.Exec("insert into hum(hum,tem,time) values (?,?,?)",datas[2], datas[3], Time)
			errHandler(err)
		}
	}
}
