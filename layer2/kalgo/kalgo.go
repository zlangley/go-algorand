package kalgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/algorand/go-algorand/daemon/algod/api/server/v2/generated"
	"github.com/algorand/go-algorand/data/basics"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Env struct {
	AlgodAddress     string
	AlgodToken       string
	SpeculationToken string
	SourcePrefix     string
}

type Runner interface {
	Run(env Env) ([]byte, error)
}

func saveToDisk(name, source, root string) (*os.File, error) {
	dirpath := filepath.Join(root, name)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dirpath, fmt.Sprintf("%s.clar", name))
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

type Cmd struct {
	Name      string
	Sender  basics.Address
	Address basics.Address
}

type InitCmd struct {
	Cmd
	Source  string
}

func (cmd *InitCmd) Run(env Env) ([]byte, error) {
	file, err := saveToDisk(cmd.Name, cmd.Source, env.SourcePrefix)
	if err != nil {
		return nil, err
	}
	kargs := baseArgs(env, cmd.Sender, cmd.Address)
	kargs = append(kargs, file.Name())
	return command(env, "init", kargs...)
}

type CallCmd struct {
	Cmd
	Function string
	Args     string
}

func (cmd *CallCmd) Run(env Env) ([]byte, error) {
	kargs := baseArgs(env, cmd.Sender, cmd.Address)
	kargs = append(kargs, fmt.Sprintf(".%s", cmd.Name), cmd.Function, cmd.Args)
	return command(env, "call", kargs...)
}

func baseArgs(env Env, sender basics.Address, address basics.Address) []string {
	args := []string{"--prefix", env.SourcePrefix}
	if !sender.IsZero() {
		args = append(args, "--sender", "@" + sender.GetUserAddress())
	}
	if !address.IsZero() {
		args = append(args, "--address", "@" + address.GetUserAddress())
	}
	return args
}

type Commitment struct {
	Contract string
	Value    string
}

type Output struct {
	Commitments    []generated.ContractCommitment
	CommitmentsRaw string `xml:"commitments"`
}

func command(env Env, subcmd string, args ...string) ([]byte, error) {
	cmd := exec.Command("./kalgo", append([]string{subcmd}, args...)...)
	cmd.Dir = os.Getenv("KALGO_HOME")
	cmd.Env = append(os.Environ(),
		"ALGOD_ADDRESS="+env.AlgodAddress,
		"ALGOD_TOKEN="+env.AlgodToken,
		"SPECULATION_TOKEN="+env.SpeculationToken,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%v: %w", cmd, err)
	}
	return out, nil
}

func ParseOutput(raw []byte) (*Output, error) {
	var out Output

	decoder := xml.NewDecoder(bytes.NewReader(raw))
	decoder.Strict = false
	err := decoder.Decode(&out)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out.CommitmentsRaw), "\n")
	commitments := make([]generated.ContractCommitment, len(lines))
	for i, line := range lines {
		// Each line is either:
		//    .contract |-> Commitment ( "<PREVIOUS COMMITMENT>" , "<NEW COMMITMENT>" )
		//    .contract |-> InitialCommitment ( "<INITIAL COMMITMENT>" )
		//    .contract |-> InitialCommitmentPromise
		split := strings.Split(strings.TrimSpace(line), " ")
		if len(split) < 4 {
			continue
		}
		var prevCommit, newCommit []byte
		if split[2] == "Commitment" {
			unq, _ := strconv.Unquote(split[4])
			prevCommit = []byte(unq)
			unq, _ = strconv.Unquote(split[6])
			newCommit = []byte(unq)
		} else if split[2] == "InitialCommitment" {
			unq, _ := strconv.Unquote(split[4])
			newCommit = []byte(unq)
		}
		commitments[i] = generated.ContractCommitment{
			Contract:           strings.TrimLeft(split[0], "."),
			PreviousCommitment: prevCommit,
			NewCommitment:      newCommit,
		}
	}
	out.Commitments = commitments
	return &out, nil
}
