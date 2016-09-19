package main

type Options struct {
	FOpts FileOptions    `group:"file"`
	GOpts GeneralOptions `group:"general"`
}

type FileOptions struct {
	Input  string `short:"i" long:"input" description:"Path to input file"`
	Output string `short:"o" long:"output" description:"Path to output file"`
}

type GeneralOptions struct {
	Method string `short:"m" long:"method" description:"Method to perform on input, options are: compress, decompress. Defaults to decompress if the input file has a '.xz' postfix. Defaults to compress if the output file has a '.xz' postfix."`
}

func newOptions() *Options {
	var opts Options
	return &opts
}
