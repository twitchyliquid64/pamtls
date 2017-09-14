#!/bin/bash

function build
{
  go build -o pamgo.so -buildmode c-shared
}

function clean
{
  go clean
  rm -f pamgo.h pamgo.so
}

function install_test
{
  export wd=`pwd`
  echo "auth required ${wd}/pamgo.so logger=syslog test" > /etc/pam.d/test_pamgo
}

function uninstall_test
{
  rm -f /etc/pam.d/test_pamgo
}

case $1 in
  "clean")
    clean
    ;;
  "clean_build")
    clean
    build
    ;;
  "build")
    build
    ;;

  "install_test")
    install_test
    ;;
  "uninstall_test")
    uninstall_test
    ;;

  *)
    echo "No argument specified. Defaulting to 'clean_build'."
    echo "In future, invoke like: \n  ${0} clean/clean_build/build/install_test/uninstall_test"
    build
    ;;
esac
