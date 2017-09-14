package main

import (
	"fmt"

	"github.com/donpark/pam"
)

type mypam struct {
	// your pam vars
}

func (mp *mypam) Authenticate(hdl pam.Handle, args pam.Args) pam.Value {
	coreInit(args)
	usr, err := hdl.GetUser()
	if err != nil {
		return pam.AuthError
	}
	info("authenticate", "Got request for user: ", usr)
	return pam.Success
}

func (mp *mypam) SetCredential(hdl pam.Handle, args pam.Args) pam.Value {
	fmt.Println("SetCredential:", args)
	return pam.Success
}

var mp mypam

func init() {
	pam.RegisterAuthHandler(&mp)
}
