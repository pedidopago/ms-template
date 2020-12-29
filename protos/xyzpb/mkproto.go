// +build ignore

// mkproto generates service.pb.go
package main

import (
	"os"
	"path"
	"strings"

	"go.pedidopago.com.br/microservices/xcmd"
)

func main() {
	// GOPATH must exist, and `go env GOPATH` always works (better than `os.Getenv("GOPATH")`)
	gopath, err := xcmd.CombinedString("go", "env", "GOPATH")
	if err != nil || gopath == "" {
		println("mkproto.go error: GOPATH not set; check `go env GOPATH`")
		os.Exit(1)
	}
	microsvcProtoParent, err := xcmd.CombinedString("go", "list", "-m", "-f", "{{.Dir}}", "go.pedidopago.com.br/microservices")
	if err != nil || gopath == "" {
		println("mkproto.go error: microservices package not found")
		os.Exit(1)
	}
	microsvcProtoParent = strings.TrimSpace(microsvcProtoParent)
	// if _, err := os.Stat(path.Join(gopath, "src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis")); err != nil {
	// 	println("mkproto.go error: grpc-gateway NOT INSTALLED. Check https://github.com/grpc-ecosystem/grpc-gateway")
	// 	os.Exit(1)
	// }

	cargs := []string{
		"-I.",                           // include this folder
		"-I/usr/local/include",          // TIP: this is the usual path for source files
		"-I" + path.Join(gopath, "src"), // include GOPATH/src for .proto definitions
		"-I" + path.Join(microsvcProtoParent, "proto/opt"), // include microservices opt for .proto type definitions
	}

	if vs := os.Getenv("MKPROTO_ARGS"); vs != "" {
		vsl := strings.Split(vs, ",")
		cargs = append(cargs, vsl...)
	} else {
		cargs = append(cargs, "--go_out=plugins=grpc:.") // output to this folder
	}

	cargs = append(cargs, "--go_opt=paths=source_relative") // uses option folder in source relative folder instead of absolute

	cargs = append(cargs, "service.proto") //TIP: you can add more proto files here

	// run protoc to generate the gRPC with grpc-gateway
	err = xcmd.RunX("", os.Stdout, os.Stderr, "protoc", cargs...)
	if err != nil {
		println("mkproto.go error: ", err.Error())
		os.Exit(1)
	}
}
