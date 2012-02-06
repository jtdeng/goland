package core

type NetworkElement interface {
    Login(username, password string)
    Logout()
    SendCommand(request string) string
}