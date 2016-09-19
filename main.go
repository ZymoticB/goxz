package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/ZymoticB/goxz/compress"
	"github.com/ZymoticB/goxz/decompress"
	"github.com/ZymoticB/goxz/output"
	"github.com/ZymoticB/goxz/xz"
)

var errExit = errors.New("sentinel error used to exit cleanly")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	parseAndRun(consoleOutput{os.Stdout})
}

func parseAndRun(out output) {
	opts, err := getOptions(os.Args[1:], out)
	if err != nil {
		if err == errExit {
			return
		}
		out.Fatalf("Failed to parse options: %v", err)
	}
	runWithOptions(*opts, out)
}

func getMethod(opts Options, out output) string {
	method := opts.GOpts.Method
	if method != "" {
		return method
	}

	// TODO: tylerd detect filetype based on magic rather than extension
	// Will need to peek first 6 bytes and then stitch it back in, or maybe close the file
	// instead of trying to pass a multiwriter all the way through to the xz module
	inputFilePath := opts.FOpts.Input
	outputFilePath := opts.FOpts.Output
	out.Printf("inputPath: %s, outputPath: %s\n", inputFilePath, outputFilePath)

	inputIsXZ := strings.HasSuffix(inputFilePath, ".xz")
	outputIsXZ := strings.HasSuffix(outputFilePath, ".xz")
	out.Printf("inputIsXZ: %v, outputIsXZ: %v\n", inputIsXZ, outputIsXZ)

	if inputIsXZ && outputIsXZ {
		out.Fatal("Input and Output are both xz unable to determine method")
	}

	if !inputIsXZ && !outputIsXZ {
		out.Fatal("Neither Input or Output are xz format; unable to determine method")
	}

	if inputIsXZ {
		method = "decompress"
	}

	if outputIsXZ {
		method = "compress"
	}
	return method
}

func runWithOptions(opts Options, out output) {
	method := getMethod(opts, out)

	inputFilePath := opts.FOpts.Input
	outputFilePath := opts.FOpts.Output

	if method == "compress" {
		err := compress.RunCompress(inputFilePath, outputFilePath)
		if err != nil {
			out.Fatalf("Failed while running compress: %v", err)
		} else {
			os.Exit(0)
		}
	}

	if method == "decompress" {
		err := decompress.RunDecompress(inputFilePath, outputFilePath)
		if err != nil {
			out.Fatalf("Failed while running compress: %v", err)
		} else {
			os.Exit(0)
		}
	}

	if method == "headers" {
		inputIsXZ := strings.HasSuffix(inputFilePath, ".xz")
		if inputIsXZ {
			err := xz.OpenFile(inputFilePath)
			os.Exit(0)
		} else {
			out.Fatalf("Headers command requires input file to be in xz format")
		}
	}

}

func newParser() (*flags.Parser, *Options) {
	opts := newOptions()
	return flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash), opts
}

func getOptions(args []string, out output) (*Options, error) {
	parser, opts := newParser()
	parser.Usage = "-i <file>.xz -o <file> -m decompress"
	parser.ShortDescription = "LZMA2 based [de]compressor"
	parser.LongDescription = `
goxz is a go implementation of LZMA2 which supports compressing and decompressing LZMA2 streams. Currently supports the xz file format
`

	if len(args) == 0 {
		parser.WriteHelp(out)
		return opts, errExit
	}

	_, err := parser.ParseArgs(args)
	if err != nil {
		if ferr, ok := err.(*flags.Error); ok {
			if ferr.Type == flags.ErrHelp {
				parser.WriteHelp(out)
				return opts, errExit
			}
		}
		return opts, err
	}

	return opts, nil
}
