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

type Program struct {
	Name   string
	Source string
}

type FunctionCall struct {
	ProgramName string
	Name        string
	Args        string
}

func saveToDisk(pgm Program, root string) (*os.File, error) {
	dirpath := filepath.Join(root, pgm.Name)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dirpath, fmt.Sprintf("%s.clar", pgm.Name))
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if _, err := file.WriteString(pgm.Source); err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return file, err
}

func Init(env Env, pgm Program) error {
	file, err := saveToDisk(pgm, env.SourcePrefix)
	if err != nil {
		return err
	}
	return kalgo(env, "init", file.Name())
}

func Call(env Env, fn FunctionCall) error {
	return kalgo(env, "call", fmt.Sprintf(".%s", fn.ProgramName), fn.Name, fn.Args)
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

type BatchItem struct {
	Program      *Program
	FunctionCall *FunctionCall
}

// FIXME: this probably doesn't belong here?
func BatchExecute(env Env, batch []BatchItem) error {
	for _, item := range batch {
		if pgm := item.Program; pgm != nil {
			if err := Init(env, *pgm); err != nil {
				return err
			}
		} else if fn := item.FunctionCall; fn != nil {
			if err := Call(env, *fn); err != nil {
				return err
			}
		}
	}
	return nil
}
