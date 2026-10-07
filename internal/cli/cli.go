package cli

import (
	"fmt"
	"os"
	"strings"
)

type Options struct {
	daemon  bool
	verbose bool
	urls    []string
	// flags for struct.
	// d - daemon, v - versbose, V - version
}

func version() {
	fmt.Println("current version is 0.1")
}
func ParseArgs(argv []string) (*Options, error) {
	opts := &Options{}
	var positionals []string

	i := 0
	for i < len(argv) {
		token := argv[i]

		if len(token) > 1 && token[0] == '-' {
			if err := parseFlags(token, opts); err != nil {
				return nil, err
			}
		} else {
			positionals = append(positionals, token)
		}
		i++
	}

	if len(positionals) == 0 {
		return nil, fmt.Errorf("no urls given")
	}
	opts.urls = []string(positionals)
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
			opts.daemon = true
		case 'v':
			opts.verbose = true
		case 'V':
			version()
			os.Exit(0)
		default:
			return fmt.Errorf("unknown flag: -%c", body[flagIdx])
		}
	}
	return nil
}
