package main

import (
	"github.com/aidar-darmenov/message-delivery-client/config"
	"log"
	"net"
	"strconv"
)
import "fmt"
import "bufio"
import "os"

func main() {

	cfg := config.NewConfiguration("config/config.json")

	conn, err := net.Dial(cfg.Type, cfg.Host+":"+strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to send: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("err reading from input", err)
		return
	}
	fmt.Fprintf(conn, text+"\n")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("err receiving message from the server", err)
		return
	}
	fmt.Print("Message from server: " + message)
}
