package kalgo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Env struct {
	AlgodAddress     string
	AlgodToken       string
	SpeculationToken string
	SourcePrefix     string
}

func saveToDisk(id, source, root string) (*os.File, error) {
	dirpath := filepath.Join(root, id)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dirpath, fmt.Sprintf("%s.clar", id))
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if _, err := file.WriteString(source); err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return file, err
}

func kalgoBaseArgs(env Env, sender *string, address *string) []string {
	args := []string{"--prefix", env.SourcePrefix}
	if sender != nil {
		args = append(args, "--sender", *sender)
	}
	if address != nil {
		args = append(args, "--address", *address)
	}
	return args
}

func kalgo(env Env, subcmd string, args ...string) error {
	kargs := []string{}
	kargs = append(kargs, subcmd)
	kargs = append(kargs, "--prefix", env.SourcePrefix)
	kargs = append(kargs, args...)

	cmd := exec.Command("./kalgo", kargs...)
	cmd.Dir = os.Getenv("KALGO_PREFIX")
	cmd.Env = append(os.Environ(),
		"ALGOD_ADDRESS="+env.AlgodAddress,
		"ALGOD_TOKEN="+env.AlgodToken,
		"SPECULATION_TOKEN="+env.SpeculationToken,
	)
	return cmd.Run()
}

type InitArgs struct {
	Name    string
	Source  string
	Sender  *string
	Address *string
}

func Init(env Env, args InitArgs) error {
	file, err := saveToDisk(args.Name, args.Source, env.SourcePrefix)
	if err != nil {
		return err
	}
	kargs := kalgoBaseArgs(env, args.Sender, args.Address)
	kargs = append(kargs, file.Name())
	return kalgo(env, "init", kargs...)
}

type CallArgs struct {
	ProgramName string
	Name        string
	Args        string
	Sender      *string
	Address     *string
}

func Call(env Env, args CallArgs) error {
	kargs := kalgoBaseArgs(env, args.Sender, args.Address)
	kargs = append(kargs, fmt.Sprintf(".%s", args.ProgramName), args.Name, args.Args)
	return kalgo(env, "call", kargs...)
}

type Command struct {
	InitArgs *InitArgs
	CallArgs *CallArgs
}

// FIXME: this probably doesn't belong here?
func BatchExecute(env Env, batch []Command) error {
	for _, item := range batch {
		if args := item.InitArgs; args != nil {
			if err := Init(env, *args); err != nil {
				return err
			}
		} else if args := item.CallArgs; args != nil {
			if err := Call(env, *args); err != nil {
				return err
			}
		}
	}
	return nil
}
