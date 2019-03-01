package main

import (
	"net"
	"fmt"
	"time"
)

func socks5_dial(){
	conn,err := net.Dial("tcp","127.0.0.1:1080")
	if err != nil{
		fmt.Println(err)
		return
	}
	_,err = conn.Write([]byte{5,1,0})
	var b [1024]byte
	status,err := conn.Read(b[:])
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(b)
	fmt.Println(status)
	fmt.Println("\n")

	before := []byte{5,1,0,3,11}
	de := []byte("youtube.com")
	port := []byte{byte(8),byte(0)}
	head_temp := append(before,de...)
	head := append(head_temp,port...)
	fmt.Println(head)
	_,err = conn.Write(head)
	if err!=nil{
		fmt.Println(err)
		return
	}

	var c [1024]byte
	status_2,err := conn.Read(c[:])
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(status_2)
	fmt.Println(c)
	conn.Close()
}

func main(){

temp := time.Now()

socks5_dial()

deply := time.Since(temp)
fmt.Println(deply)
}