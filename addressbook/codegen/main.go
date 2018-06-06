package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// protoc --plugin=protoc-gen-custom=codegen --custom_out=. -I ../pb ../pb/addressbook.proto

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err, "reading input")
	}

	var request plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(data, &request); err != nil {
		log.Fatalln(err, "parsing input proto")
	}

	if len(request.FileToGenerate) == 0 {
		log.Fatalln("no files to generate")
	}

	requestJSON, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	var response plugin.CodeGeneratorResponse
	response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String("raw.json"),
		Content: proto.String(string(requestJSON)),
	})

	data, err = proto.Marshal(&response)
	if err != nil {
		log.Fatalln(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		log.Fatalln(err, "failed to write output proto")
	}
}
