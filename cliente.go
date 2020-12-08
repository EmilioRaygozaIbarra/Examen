package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"time"
)

type Mensaje struct {
	Men     string
	Usuario string
}

var CListaMensaje = []Mensaje{}
var copiaCListaMensaje = []Mensaje{}
var cuser string
var userID int64
var bandera bool

func mostrarMensaje(c *rpc.Client) {
	bandera = false
	for {
		error := c.Call("Server.EnviarMensaje", bandera, &CListaMensaje)
		if error != nil {
			fmt.Println(error)
		}
		l := len(CListaMensaje)
		if l > 0 {
			fmt.Println("--------")
			for i := 0; i < l; i++ {
				fmt.Print(CListaMensaje[i].Usuario, ": ")
				fmt.Println(CListaMensaje[i].Men)
				if i == l-1 {
					m := CListaMensaje[l-1]
					copiaCListaMensaje = append(copiaCListaMensaje, m)
				}

			}
			fmt.Println("--------")
		}

	}
}

func client() {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	var opc int64
	var m string
	b := true
	scanner := bufio.NewScanner(os.Stdin)
	for b {
		fmt.Println("Nombre del Usuario: ")
		scanner.Scan()
		cuser = scanner.Text()
		error := c.Call("Server.AgregarUsuario", cuser, &userID)
		if error != nil {
			fmt.Println(error)
		} else {
			b = false
		}
	}
	go mostrarMensaje(c)
	for opc != 4 {
		/*cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()*/
		fmt.Println("1.- Enviar Mensaje")
		fmt.Println("2.- Enviar Archivo")
		fmt.Println("3.- Mostrar Chat")
		fmt.Println("4.- Salir")
		fmt.Scanln(&opc)
		switch opc {
		case 1:
			fmt.Print("Mensaje: ")
			scanner.Scan()
			m = scanner.Text()
			men := Mensaje{Men: m, Usuario: cuser}
			error := c.Call("Server.RecibirMensaje", men, &CListaMensaje)
			if error != nil {
				fmt.Println(error)
			}
			time.Sleep(time.Millisecond * 500)
			bandera = false
		case 2:
			reply := cuser
			archivo, error := os.Open("clienteTXT.txt")
			if error != nil {
				fmt.Println(error)
			}
			defer func() {
				if error = archivo.Close(); error != nil {
					fmt.Println(error)
				}
			}()
			mandar, error := ioutil.ReadAll(archivo)
			error = c.Call("Server.RecibirUsuarioArchivo", reply, &reply)
			error = c.Call("Server.RecibirArchivo", mandar, &reply)
			if error != nil {
				fmt.Println(error)
			}
			fmt.Println(reply)
		case 3:
			l := len(copiaCListaMensaje)
			fmt.Println(l)
			if l > 0 {
				fmt.Println("--------")
				fmt.Println("Mensajes Enviados: ")
				for i := 0; i < l; i++ {
					fmt.Print(copiaCListaMensaje[i].Usuario, ": ")
					fmt.Println(copiaCListaMensaje[i].Men)
				}
				fmt.Println("--------")
			}
		}
	}
}

func main() {
	client()
}
