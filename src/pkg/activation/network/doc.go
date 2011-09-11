/*
an abstraction of network element layers
*/
package network

type NetworkElement interface {
	Login(username, password string)
	Logout()
	SendCommand(request string) string
}