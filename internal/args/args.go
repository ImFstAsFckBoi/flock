package args

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ImFstAsFckBoi/flock/pkg/utils"
	"github.com/akamensky/argparse"
)

type ArgsRaw struct {
	mode     string
	file_in  string
	file_out string
}

type ArgsProcessed struct {
	mode     MODE
	file_in  io.Reader
	file_out io.Writer
}

type MODE int

const (
	MODE_AUTO MODE = iota
	MODE_INFO
	MODE_LOCK
	MODE_UNLOCK
	MODE_VALIDATE
)

func ReadArgs() (*ArgsRaw, error) {
	parser := argparse.NewParser("flock", "Quickly password lock files and directories")
	args := ArgsRaw{}
	mode := parser.String("m", "mode",
		&argparse.Options{
			Required: false,
			Validate: nil,
			Help:     "Operating mode: [AUTO|LOCK|UNLOCK]",
			Default:  "AUTO",
		})

	file_in := parser.String("f", "file",
		&argparse.Options{
			Required: false,
			Validate: nil,
			Help:     "Input file: [<file path>|'stdin']",
			Default:  "stdin",
		})

	file_out := parser.String("o", "out",
		&argparse.Options{
			Required: false,
			Validate: nil,
			Help:     "Output file: [<file path>|'stdout']. Default: \n\tif file=stdin: stdout\n\tif file='<path>±.locked': <path>±.locked",
		})

	parser.Parse(os.Args)

	if file_out == nil {
		if *file_in == "stdin" {
			args.file_out = "stdout"
		}

		if *mode == "LOCK" {
			args.file_out = *file_in + ".locked"
		} else if *mode == "UNLOCK" {
			found := false
			args.file_out, found = strings.CutSuffix(*file_in, ".locked")
			if !found {
				return nil, errors.New("File output ambiguous, please provide explicitly")
			}
		} else {
			args.file_out = ""
		}
	}

	args.mode = strings.ToUpper(*mode)

	return &args, nil
}

func ModeFromString(modeStr string) (MODE, error) {
	var mode MODE
	switch strings.ToUpper(modeStr) {
	case "AUTO":
		mode = MODE_AUTO
	case "INFO":
		mode = MODE_INFO
	case "LOCK":
		mode = MODE_LOCK
	case "UNLOCK":
		mode = MODE_UNLOCK
	case "VALIDATE":
		mode = MODE_VALIDATE
	default:
		return 100, errors.New("Invalid modeStr")
	}

	return mode, nil
}

func ProcessArgs(arg *ArgsRaw) (*ArgsProcessed, error) {
	argP := ArgsProcessed{}

	var openMode os.FileMode
	if arg.file_in == "stdin" {
		argP.file_in = os.Stdin
		openMode = 0600
	} else {

		fIn, err1 := os.Open(arg.file_in)
		fInInfo, err2 := fIn.Stat()
		openMode = fInInfo.Mode()
		if errors.Join(err1, err2) != nil {
			return nil, errors.New("Could not open input file '" + arg.file_in + "'")
		}

		argP.file_in = fIn
	}

	if arg.file_out != "" {
		_, err := os.Stat(arg.file_out)
		if err == nil {
			yn, err := utils.YNPrompt("File "+arg.file_out+" already exists. Do you want to overwrite it?", utils.BIAS_NO)
			if err != nil {
				return nil, err
			} else if !yn {
				return nil, errors.New("User cancelled program")
			}
		}

		argP.file_out, err = os.OpenFile(arg.file_out, os.O_CREATE|os.O_RDWR|os.O_TRUNC, openMode)
		if err != nil {
			return nil, err
		}
	}

	return &argP, nil
}

func GetArgs() (*ArgsProcessed, error) {
	rawArgs, err := ReadArgs()
	if err != nil {
		return nil, err
	}

	args, err := ProcessArgs(rawArgs)

	if err != nil {
		return nil, err
	}

	return args, nil
}
