package cli

import (
	"fmt"
	"os"
	"strings"
)

type Options struct {
	Daemon   bool
	Verbose  bool
	Config   bool
	Unsafe   bool
	ConfPath string
	Urls     []string
	// flags for struct.
	// d - Daemon, v - versbose, V - version
}

func version() {
	fmt.Println("current version is 0.1")
}
func ParseArgs(argv []string) (*Options, error) {
	opts := &Options{}
	var positionals []string

	i := 0
	for i < len(argv) { // will probably need to replace it with for loop for conf
		token := argv[i]

		if len(token) > 1 && token[0] == '-' {

			if strings.Contains(token[1:], "c") {
				if i+1 >= len(argv) {
					return nil, fmt.Errorf("-c requires config file path")
				}
				opts.Config = true
				opts.ConfPath = argv[i+1]
				i++
			}
			if err := parseFlags(token, opts); err != nil {
				return nil, err
			}
		} else {
			positionals = append(positionals, parseProtocols(token, opts))
		}
		i++
	}

	if len(positionals) == 0 {
		return nil, fmt.Errorf("no Urls given")
	}

	opts.Urls = positionals

	return opts, nil
}

func parseFlags(token string, opts *Options) error {

	if !strings.HasPrefix(token, "-") {
		return fmt.Errorf("not a flag: %s", token)
	}

	body := token[1:]
	if body == "" {
		return fmt.Errorf("unknown flag: -")
	}
	for flagIdx := 0; flagIdx < len(body); flagIdx++ {
		switch body[flagIdx] {
		case 'd':
			opts.Daemon = true
		case 'c':
			opts.Config = true
		case 'v':
			opts.Verbose = true
		case 'V':
			version()
			os.Exit(0)
		case 'u':
			opts.Unsafe = true
		default:
			return fmt.Errorf("unknown flag: -%c", body[flagIdx])
		}
	}
	return nil
}

func parseProtocols(pos string, opts *Options) string {
	if strings.HasPrefix(pos, "https://") || strings.HasPrefix(pos, "http://") {
		return pos
	}
	if opts.Unsafe {
		return "http://" + pos
	}
	return "https://" + pos
}
