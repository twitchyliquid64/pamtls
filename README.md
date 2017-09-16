# pamtls

pamtls allows you to hit a JSON endpoint and ask if the given username+password+other should be allowed or not.

You can use this to move your passwords for \*nix user accounts to a centralised system without LDAP/Kerberos, or check a webservice to deny access (remote lockout).

pamtls also supports client certs to allow full mutual authentication between the client and server.

## Build

#### Install dependencies

You need pam-dev (C headers and libpam) on your system.

Debian, Ubuntu, Linux Mint: `sudo apt-get install libpam0g-dev`

CentOS, Fedora, RHEL: `sudo yum install gcc pam-devel`

#### Run the build script

```shell
export GOPATH=`pwd` #set the GOPATH to the root directory of pamtls
./build.sh
# This will create pamtls.so in pamtls/
```

## Install

1. Copy `pamtls.so` to somewhere in your filesystem. It must be accessible whenever authorization is performed (someone tries to SSH, boots to a lock screen etc).
2. Modify your `/etc/pam.d/common-auth` file to invoke your PAM module. If you don't know how to do this, I recommend some reading about PAM modules, the `/etc/pam.d` directory, and the layout of the PAM configuration files.
3. Decide on your options and append them to the PAM module line. The options are as follows:

| Option name   | Explanation | Example |
| ------------- | ----------- | --------|
| debug         | When set, will print lots of helpful messages in the syslog.  | `debug=true` |
| token         | This is a string which will be sent to the web endpoint for all requests on this machine. Useful for giving a machine a unique identifier. | `token=blah-blah` |
| url           | Base URL which pamtls will make requests to, to ask if the user should be permitted. | `url=https://<mydomain>/auth` |
| verify        | Sets the TLS verification mode. Defaults to full verification against System root certificates. Other values are 'pinned' and 'insecure'. Do not get lazy and use insecure, you might as well just not have a password then. | `verify=pinned` |
| root          | Only valid when `verify=pinned`. This should be the full path to a PEM-encoded certificate. The TLS cert presented by the server must be signed by this cert. | `root=/etc/ca.pem` |
| cert          | Only valid when `key` is set as well. This is a path to the PEM-encoded client certificate to use for TLS connections. | `cert=/etc/client_cert.pem` |
| key           | Only valid when `cert`is set as well. This is a path to the PEM-encoded client key to use for TLS connections. | `cert=/etc/client_key.pem` |
| prompt        | When present and set to `password`, pamtls will not ask the server which questions the user should be asked, and instead just ask for the password. You probably want to set this, unless you are doing funky stuff and want to prompt the user for other credentials as well (such as an OTP/2-factor code). | `prompt=password` |

For example, a pamtls config line which makes requests to `https://example.com/pam/authenticate` and does normal TLS server validation is:

```
auth required /boot/pamtls.so url=https://example.com/pam prompt=password token=computer-x
```

## Server setup

*NOTE: You need to be very careful with your implementation. Make sure you have rate limiting to stop someone bruteforcing, and you should use the TLS-client certificates feature to prevent any devices from talking to your server which are not meant to.*

**TL;DR** you need to implement a HTTPS endpoint which accepts HTTP POST requests containing JSON, and replies with JSON.

#### Authentication endpoint

pamtls will make a request to `url`/authenticate. (eg: if you have `url=https://blah` in your configuration, the request will be made to `https://blah/authenticate`.)

The request will contain JSON like this:
```json
"{"user":"<username>","token":"<token>","responses":[["<password>"]]}"```

You should reply like this:
```json
{"Success": true}
```

If you want to deny access, `Success` should be false (obviously).

If you encounter a system error and want to deny access but also log the failure, send this:

```json
{"Error": "Error message to appear in syslog here"}
```

Lastly, if you want print a message to the user (regardless of the Value of `Success`), return a `Message` field which contains text.

#### Prompts endpoint

To support 2FA, pamtls will request a list of prompts for the user if `prompt=password` is not present in the configuration line. This request will be made to `url`/authPrompts.

I havent documented this yet, but if you look at `getPrompts` and `getPromptsResponse` in `requests.go` you should be able to work out the JSON format. All responses are sent in the `responses` list in the subsequent `/authentication` API call.
