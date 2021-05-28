package apollo

var (
	defaultNamespace          = "application"
	defaultBackFile           = "/tmp/.go-apollo"
	defaultLongPollerInterval = 1
)

type Option func(*Options)

type Options struct {
	ConvertStruct      bool   //自动转为struct
	BackFile           string //备份文件目录
	LongPollerInterval int64    //轮训时间间隔
	DefaultNamespace   string //默认namespace
}

func DefaultNamespace(namespace string) Option {
	return func(o *Options) {
		o.DefaultNamespace = namespace
	}
}
func WithoutConvertStruct() Option {
	return func(o *Options) {
		o.ConvertStruct = false
	}
}

func BackFile(backFile string) Option {
	return func(o *Options) {
		o.BackFile = backFile
	}
}

func newOption(opts ...Option) Options {
	var options = Options{
		ConvertStruct:      true,
		BackFile:           defaultBackFile,
		LongPollerInterval: int64(defaultLongPollerInterval),
		DefaultNamespace:   defaultNamespace,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return options
}
