package bpe

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
)

func ExampleTrain() {
	// Any reader could be used for reading training data.
	// To train model using file as a source:
	// source, err := os.Open("/path/to/file.txt")
	// if err != nil { /* Handle error */ }
	// model, err := Train(source)
	// Check available TrainOption for customization.
	source := strings.NewReader("Lorem Ipsum")
	ctx := context.Background()
	model, err := Train(ctx, source)
	if err != nil {
		log.Fatalln(err)
	}

	// Now you have trained model.
	// You can start using it or export to save time for future usages.
	// Check Export() function.

	fmt.Printf("%d", len(model.vocab))
	// Output: 29
}

func ExampleExport() {
	// To export model we need to create one.
	source := strings.NewReader("x")
	ctx := context.Background()
	model, err := Train(ctx, source)

	// For simplicity we're going to export model to bytes.Buffer,
	// but you can export it wherever you want. E.g. to file:
	// destination, err := os.Open("/path/to/exported-model.json")
	// if err != nil { /* Handle error */ }
	destination := bytes.NewBuffer(nil)

	// If you want to store model in other format than JSON
	// you can create your one encoder and use it with WithEncoder() option.
	err = Export(model, destination)
	if err != nil {
		log.Fatalln(err)
	}

	// vocab = <w>x</w>, default JSON formatter converts string to unicode sequences.
	fmt.Println(destination)
	// Output: {"max_token_length":8,"vocab":["\u003cw\u003ex\u003c/w\u003e"]}
}

func ExampleImport() {
	// You can import model from various sources using io.Reader.
	// E.g. import from file:
	// source, err := os.Open("/path/to/file.json")
	// if err != nil { /* Handle error */ }
	// model, err := Import(source)
	source := strings.NewReader(`{"vocab":["token"]}`)
	model, err := Import(source)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", model.vocab)
	// Output: map[token:{}]
}

func ExampleBPE_Encode() {
	source := strings.NewReader(`Lorem Ipsum`)
	ctx := context.Background()
	model, err := Train(ctx, source)
	if err != nil {
		log.Fatalln(err)
	}

	text := strings.NewReader(`Ipsum Lorem`)
	tokens, err := model.Encode(text)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%v", tokens)
	// Output: [<s> <w>Ipsum</w> <w>Lorem</w> </s>]
}

func ExampleBPE_Decode() {
	model := &BPE{}
	tokens := []string{
		"<s>", "<w>Th", "is</w>", "<w>is</w>", "<w>j", "u", "st</w>", "<w>an</w>", "<w>ex", "ample", ".</w>", "</s>",
	}

	result, err := model.Decode(tokens)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%v", result)
	// Output: This is just an example.
}
