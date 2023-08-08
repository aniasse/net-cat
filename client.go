package main

import (
	"fmt"
	"net"
	"strings"
)

var (
	User       = make(map[net.Conn]string)
	MyClient   []string
	Allmessage []string
)

func Client(conn net.Conn) {

	if len(MyClient) <= 10 {
		for {
			avatar := ChargeNetcat("avatar.txt") //Chargement de l'avatar
			conn.Write([]byte("Welcome to TCP-Chat!" + avatar))
			conn.Write([]byte("[ENTER YOUR NAME]:"))
			username := ReadUsername(conn)

			check := CheckUsername(MyClient, username)

			for !check || username == "" || username == " " { //Control du nom d'utilisateur saisie
				conn.Write([]byte("[ENTER YOUR NAME]:"))
				username = ReadUsername(conn)
				check = CheckUsername(MyClient, username)
			}
			MyClient = append(MyClient, username)

			if len(MyClient) != 1 && len(Allmessage) == 0 { //Alerte au cas ou il y a plusieurs utilisateur qui se connecte en meme temps
				for connex, usr := range User {
					if connex != conn { //Pour l'utilisateur en cours
						connex.Write([]byte("\n" + username + " has joined our chat...\n"))
						text := WriteMsg(usr, "")
						connex.Write([]byte(text))
					}
					//Pour les autres utilisateurs present dans le chat
					text1 := WriteMsg(username, "")
					conn.Write([]byte(text1))
					conn.Write([]byte("\n" + usr + " has joined our chat...\n"))
				}
			} else {

				for connex, usr := range User {
					if connex != conn {
						connex.Write([]byte("\n" + username + " has joined our chat...\n"))
						text := WriteMsg(usr, "")
						connex.Write([]byte(text))
					}
				}
			}
			User[conn] = username // le nom de chaque utilisateur

			if len(MyClient) != 1 && len(Allmessage) != 0 { //S'il y a deja un utilisateur et que le chat n'est pas vide
				text := WriteMsg(username, "")
				// conn.Write([]byte(text + "\n"))

				for _, mesg := range Allmessage { //Affiche des anciens messages pour le nouveau utilisateur
					conn.Write([]byte(mesg))
				}
				conn.Write([]byte(text))
			} else {
				text := WriteMsg(username, "")
				conn.Write([]byte(text))
			}

			for {

				msg := make([]byte, 1024)
				n, err := conn.Read(msg)

				if err != nil {
					for connex, usr := range User {
						if connex != conn {
							connex.Write([]byte("\n" + username + " has left our chat...\n"))
							ident := WriteMsg(usr, "")
							connex.Write([]byte(ident))
						}
					}
					MyClient = Delete(MyClient, username)
					conn.Close()
					break
				}
				str := string(msg[:n])
				str = strings.Trim(str, "\n")

				if str == "./chngusr" { //Pour changer de nom d'utilisateur
					conn.Write([]byte("[CHANGE YOUR NAME]:"))
					newusername := ReadUsername(conn)
					check := CheckUsername(MyClient, newusername)
					for !check || newusername == "" || newusername == " " { //Control du nom d'utilisateur saisie
						conn.Write([]byte("[CHANGE YOUR NAME]:"))
						newusername = ReadUsername(conn)
						check = CheckUsername(MyClient, newusername)
					}
					for connex, usr := range User {
						if connex != conn {
							text := fmt.Sprintf("\n"+username+" changed its user name to %s\n", newusername)
							connex.Write([]byte(text))
							ident := WriteMsg(usr, "")
							connex.Write([]byte(ident))
						}
					}
					MyClient = Delete(MyClient, username)
					username = newusername
					MyClient = append(MyClient, username)
					User[conn] = username
				}
				message := ""
				if str != "" && str != "./chngusr" { //stockage des messages
					message = WriteMsg(username, string(msg))
					Allmessage = append(Allmessage, message)
				}

				for con := range isconnexion {
					if con != conn && str != "" && str != "./chngusr" { //Si la connexion est differente de la connexion en cours,->
						con.Write([]byte("\n" + message)) //on affiche le message
						for key, name := range User {     //Permettre a l'utilisateur de pouvoir a nouveau saisir un message
							if key == con && str != "" {
								mes := WriteMsg(name, "")
								con.Write([]byte(mes))
							}
						}
					} else if con == conn {
						mes := WriteMsg(username, "")
						con.Write([]byte(mes))
					}
				}
			}
		}
	} else {
		conn.Write([]byte("Sorry 😌 the group has reached its limit"))
		conn.Close()
	}
}
