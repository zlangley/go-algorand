package kalgo

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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
	Run(env Env) (*Output, error)
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

func (cmd *InitCmd) Run(env Env) (*Output, error) {
	file, err := saveToDisk(cmd.Id, cmd.Source, env.SourcePrefix)
	if err != nil {
		return nil, err
	}
	kargs := baseArgs(env, cmd.Sender, cmd.Address)
	kargs = append(kargs, file.Name())
	return command(env, "init", kargs...)
}

type CallCmd struct {
	Id       string
	Function string
	Args     string
	Sender   *string
	Address  *string
}

func (cmd *CallCmd) Run(env Env) (*Output, error) {
	kargs := baseArgs(env, cmd.Sender, cmd.Address)
	kargs = append(kargs, fmt.Sprintf(".%s", cmd.Id), cmd.Function, cmd.Args)
	return command(env, "call", kargs...)
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

type Output struct {
	Commitments string `xml:"commitments"`
}

func command(env Env, subcmd string, args ...string) (*Output, error) {
	cmd := exec.Command("./kalgo", append([]string{subcmd}, args...)...)
	cmd.Dir = os.Getenv("KALGO_PREFIX")
	cmd.Env = append(os.Environ(),
		"ALGOD_ADDRESS="+env.AlgodAddress,
		"ALGOD_TOKEN="+env.AlgodToken,
		"SPECULATION_TOKEN="+env.SpeculationToken,
	)
	rawout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Some random stuff gets printed before the XML output... try to skip over it.
	reader := bufio.NewReader(bytes.NewReader(rawout))
	for {
		next, err := reader.Peek(1)
		if err != nil {
			return nil, err
		}
		if string(next) == "<" {
			break
		}
		reader.ReadLine()
	}

	rawout, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var out Output
	xml.Unmarshal(rawout, &out)
	return &out, nil
}
