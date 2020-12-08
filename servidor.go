package main

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
)

type Usuario struct {
	Nombre string
}

type Mensajes struct {
	Men     string
	Usuario string
}

var listaUsuario = list.New()

var ListaMensaje = []Mensajes{}
var user int64
var mensajeRecibido = false
var usuarioArchivo string

type Server struct{}

func (this *Server) AgregarUsuario(usuario string, reply *int64) error {
	if usuario == "" {
		return errors.New("Favor de agregar un usuario valido")
	} else {
		for e := listaUsuario.Front(); e != nil; e = e.Next() {
			if usuario == e.Value {
				return errors.New("Este usuario ya existe")
			}
		}
		listaUsuario.PushBack(usuario)
		*reply = user
		user++
		fmt.Println("Se conecto: ", usuario)
		return nil
	}
}

func (this *Server) RecibirMensaje(mensaje Mensajes, reply *[]Mensajes) error {
	if mensaje.Men == "" {
		return errors.New("Mensaje vacio no enviado")
	} else {
		ListaMensaje = append(ListaMensaje, mensaje)
		fmt.Print(mensaje.Usuario, ": ")
		fmt.Println(mensaje.Men)
		l := len(ListaMensaje)
		for i := 0; i < l; i++ {
			m := ListaMensaje[i]
			*reply = append(*reply, m)
		}
		mensajeRecibido = true
		return nil
	}

}

func (this *Server) RecibirUsuarioArchivo(user string, reply *string) error {
	usuarioArchivo = user
	return nil
}

func (this *Server) RecibirArchivo(aux []byte, reply *string) error {
	if aux == nil {
		return errors.New("El archivo esta vacio")
	} else {
		archivo, error := os.Create("archivoEnServidorDe" + usuarioArchivo + ".txt")
		if error != nil {
			fmt.Println(error)
			return errors.New("No se puede guardar el archivo")
			//return nil
		}
		defer archivo.Close()
		archivo.Write(aux)
		*reply = "Archivo guardado de manera exitosa"
		return nil
	}
	return nil
}

func (this *Server) EnviarMensaje(Bandera bool, Reply *[]Mensajes) error {
	if mensajeRecibido {
		time.Sleep(time.Millisecond * 300)
		l := len(ListaMensaje)
		for i := 0; i < l; i++ {
			m := ListaMensaje[i]
			*Reply = append(*Reply, m)
		}
		mensajeRecibido = false
		return nil
	}
	return nil
}

func server() {
	rpc.Register(new(Server))
	ln, error := net.Listen("tcp", ":9999")
	if error != nil {
		fmt.Println(error)
	}
	for {
		c, error := ln.Accept()
		if error != nil {
			fmt.Println(error)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func respaldoChat(){
	archivo, error := os.Create("respaldoDeChat.txt")
	if error != nil {
		fmt.Println(error)
		//return nil
	}
	defer archivo.Close()
	l:=len(ListaMensaje)
	if l<1{
		fmt.Println("No hay mensaje todavia")
	}else{
		for i:=0;i<l;i++{
			aux:=ListaMensaje[i].Usuario
			archivo.WriteString(aux + ": ")
			aux=ListaMensaje[i].Men
			archivo.WriteString(aux + "\n")
		}
	}
	//*reply = "Archivo guardado de manera exitosa"
}

func main() {
	go server()
	opc := 1

	for opc != 0 {
		fmt.Println("1.- Respaldar Chat")
		fmt.Println("0. Cerrar Servidor")
		fmt.Scanln(&opc)
		if opc==1{
			respaldoChat()
		}
	}
}
