package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	_ "github.com/go-sql-driver/mysql"
)
// MssqlConnect is type  
type MssqlConnect struct {
	Database string `json:"database"`
	User struct {
		ID  string `json:"id"`
		Pwd string `json:"pwd"`
	} `json:"user"`
	Host struct {
		Address string `json:"address"`
		Port    string    `json:"port"`
	} `json:"host"`
}

// Park datastruct
type Park struct {
	Name string `json:"name"`
	Total int `json:"total"`
	Empty int `json:"empty"`
}

// Point type of x, y
type Point struct{
	X float64 `json:"x"`
	Y float64 `json:"y"`
}


func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	db *sql.DB
)


func main() {
	config , err := LoadConfig()
	errHandler(err)

	// sever open
	l, err := net.Listen("tcp",config.Host.Port)
	errHandler(err)
	defer l.Close()

	dbname := config.User.ID+":"+config.User.Pwd+"@/"+"park"
	db, err = sql.Open("mysql",dbname) 
	errHandler(err)
	var name string
	err = db.QueryRow("select juso from parkID limit 1").Scan(&name)
	errHandler(err)
	fmt.Println(name)

	defer db.Close()

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
			recv := recvBuf[:n]
			var data Point
			err := json.Unmarshal(recv,&data)
			errHandler(err)
			fmt.Println("receve = ",data)
			onepark := Park{}
			err = db.QueryRow("select juso,total,empty from parkID where x = $1 and y = $2",data.X,data.Y).Scan(&onepark.Name,&onepark.Total,&onepark.Empty)
			errHandler(err)
			sendbyte, _ := json.Marshal(onepark)
			conn.Write(sendbyte)
			fmt.Printf("%s total = %d, empty = %d",onepark.Name,onepark.Total,onepark.Empty)
		}
	}
}

// LoadConfig is retturn config json
func LoadConfig() (MssqlConnect, error){
	var config MssqlConnect
	file , err := os.Open("config.json")
	defer file.Close()
	errHandler(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	errHandler(err)
	return config, err
}
