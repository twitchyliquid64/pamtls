package main

import (
	"fmt"
	"strings"

	"github.com/donpark/pam"
)

type mypam struct {
	// your pam vars
}

func (mp *mypam) Authenticate(hdl pam.Handle, args pam.Args) pam.Value {
	coreInit(args)
	err := tlsInit(args)
	if err != nil {
		fatal("TLSINIT", err)
		return pam.AuthError
	}

	usr, err := hdl.GetUser()
	if err != nil {
		return pam.AuthError
	}
	if isDebugMode {
		info("DEBUG-AUTH", "Got request for user: ", usr)
	}

	var prompts *getPromptsResponse
	if strings.ToLower(args["prompt"]) == "password" {
		prompts = &getPromptsResponse{Prompts: []pam.Message{pam.Message{Msg: "Password: ", Style: pam.MessageEchoOff}}}
	} else {
		prompts, err = getAuthPrompts(usr, args["token"])
		if err != nil {
			fatal("GET-PROMPTS", err)
			return pam.AuthInfoUnavailable
		}
	}

	responses := make([][]string, len(prompts.Prompts))
	for i, prompt := range prompts.Prompts {
		resp, errConvo := hdl.Conversation(prompt)
		if isDebugMode {
			info("DEBUG-PROMPT", fmt.Sprintf("prompt=%d - response:%+v - err:%v", i, resp, errConvo))
		}
		if errConvo != nil {
			fatal("CONVERSATION", errConvo)
			return pam.AuthInfoUnavailable
		}
		responses[i] = resp
	}

	authResponse, err := authenticate(usr, args["token"], responses)
	if err != nil {
		fatal("GET-AUTH", err)
		return pam.AuthInfoUnavailable
	}
	if authResponse.Message != "" {
		_, errConvo := hdl.Conversation(pam.Message{Msg: authResponse.Message, Style: pam.MessageTextInfo})
		if errConvo != nil {
			info("MSG-ERR", errConvo)
		}
	}

	if authResponse.Success {
		return pam.Success
	}
	return pam.PermissionDenied
}

func (mp *mypam) SetCredential(hdl pam.Handle, args pam.Args) pam.Value {
	// fmt.Println("SetCredential:", args)
	return pam.Success
}

var mp mypam

func init() {
	pam.RegisterAuthHandler(&mp)
}

func main() {}
