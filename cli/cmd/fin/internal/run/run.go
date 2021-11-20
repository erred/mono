package run

import "flag"

type Cmd func(Options, []string) error

type Options struct {
	File string
}

func NewOptions(fs *flag.FlagSet) Options {
	var o Options
	o.InitFlags(fs)
	return o
}

func (o *Options) InitFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.File, "f", "fin.prototext", "name of transactions file")
}
