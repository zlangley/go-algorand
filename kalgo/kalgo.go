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

type Cmd interface {
	Run(env Env) error
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

type InitCmd struct {
	Id      string
	Source  string
	Sender  *string
	Address *string
}

func (cmd *InitCmd) Run(env Env) error {
	file, err := saveToDisk(cmd.Id, cmd.Source, env.SourcePrefix)
	if err != nil {
		return err
	}
	kargs := baseArgs(env, cmd.Sender, cmd.Address)
	kargs = append(kargs, file.Name())
	return command(env, "init", kargs...).Run()
}

type CallCmd struct {
	Id       string
	Function string
	Args     string
	Sender   *string
	Address  *string
}

func (cmd *CallCmd) Run(env Env) error {
	kargs := baseArgs(env, cmd.Sender, cmd.Address)
	kargs = append(kargs, fmt.Sprintf(".%s", cmd.Id), cmd.Function, cmd.Args)
	return command(env, "call", kargs...).Run()
}

func baseArgs(env Env, sender *string, address *string) []string {
	args := []string{"--prefix", env.SourcePrefix}
	if sender != nil {
		args = append(args, "--sender", *sender)
	}
	if address != nil {
		args = append(args, "--address", *address)
	}
	return args
}

func command(env Env, subcmd string, args ...string) *exec.Cmd {
	cmd := exec.Command("./kalgo", append([]string{subcmd}, args...)...)
	cmd.Dir = os.Getenv("KALGO_PREFIX")
	cmd.Env = append(os.Environ(),
		"ALGOD_ADDRESS="+env.AlgodAddress,
		"ALGOD_TOKEN="+env.AlgodToken,
		"SPECULATION_TOKEN="+env.SpeculationToken,
	)
	return cmd
}
