package test

import (
	"fmt"
	"os"
	"os/exec"
)

func GetTestHelperProcessArgs() []string {
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	return args
}

func SkipTestHelperProcess() bool {
	return os.Getenv("GO_WANT_HELPER_PROCESS") != "1"
}

func StubExecCommand(testHelper string, desiredOutput string) func(...string) *exec.Cmd {
	return func(args ...string) *exec.Cmd {
		cs := []string{
			fmt.Sprintf("-test.run=%s", testHelper),
			"--", desiredOutput}
		cs = append(cs, args...)
		env := []string{
			"GO_WANT_HELPER_PROCESS=1",
		}

		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(env, os.Environ()...)
		return cmd
	}
}

type execStub struct {
	testHelper    string
	desiredOutput string
}

type execCall struct {
	arguments []string
}

type FakeExecer struct {
	calls []*execCall
	count int
	stubs []*execStub
}

func (fe *FakeExecer) StubbedExec(args ...string) *exec.Cmd {
	call := fe.count
	fe.count += 1
	fe.calls = append(fe.calls, &execCall{arguments: args})

	if len(fe.stubs) <= call {
		err := fmt.Sprintf("fake execer received more calls than it has stubs. most recent call: %v", fe.calls[len(fe.calls)-1])
		panic(err)
	}

	stub := fe.stubs[call]
	cs := []string{
		fmt.Sprintf("-test.run=%s", stub.testHelper),
		"--", stub.desiredOutput}
	cs = append(cs, args...)
	env := []string{
		"GO_WANT_HELPER_PROCESS=1",
	}

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}

func (fe *FakeExecer) StubExec(testHelper, desiredOutput string) {
	fe.stubs = append(fe.stubs, &execStub{testHelper, desiredOutput})
}
