package shells

import (
	"bytes"
	"fmt"
	"strings"

	"gitm/internal/helpers"
)

type stringQuoter func(string) string

type Shell interface {
	GenerateScript(script []string) (string, error)
}

func NewBashShell(workDir string) *BashShell {
	return &BashShell{workDir: workDir}
}

type BashShell struct {
	workDir string
}

func (bs *BashShell) GenerateScript(script []string) (string, error) {
	w := NewBashWriter()
	bs.goToWorkDir(w)
	return bs.generateScript(w, script)
}

func (bs *BashShell) generateScript(w *BashWriter, script []string) (string, error) {
	bs.writeCommands(w, script...)

	return w.Finish(false), nil
}

func (bs *BashShell) writeCommands(w *BashWriter, script ...string) {
	for _, command := range script {
		command = strings.TrimSpace(command)

		if command == "" {
			continue
		}

		w.Noticef("$ %s", command)
		w.Line(command)
		w.CheckForErrors()
	}
}

func (bs *BashShell) goToWorkDir(w *BashWriter) {
	w.Cd(bs.workDir)
}

func (bs *BashShell) cleanup(w *BashWriter) {
	w.RmDir(bs.workDir)
}

func NewBashWriter() *BashWriter {
	return &BashWriter{}
}

type BashWriter struct {
	bytes.Buffer
	checkError bool
}

func (b *BashWriter) Command(command string, arguments ...string) {
	b.Line(b.buildCommand(helpers.ShellEscape, command, arguments...))
	b.CheckForErrors()
}

func (b *BashWriter) Line(text string) {
	_, _ = b.WriteString(text + "\n")
}

func (b *BashWriter) CheckForErrors() {
	b.Line("_gitm_code=$?; if [ ${_gitm_code} -ne 0 ]; then ${exit _gitm_code}; fi")
}

func (b *BashWriter) Finish(trace bool) string {
	var buf strings.Builder

	if trace {
		buf.WriteString("set -o xtrace\n")
	}

	buf.WriteString("if set -o | grep pipefail > /dev/null; then set -o pipefail; fi; set -o errexit\n")
	buf.WriteString("set +o noclobber\n")

	buf.WriteString(": | (eval " + b.escape(b.String()) + ")\n")

	buf.WriteString("exit 0\n")

	return buf.String()
}

func (b *BashWriter) escape(input string) string {
	return helpers.ShellEscape(input)
}

func (b *BashWriter) buildCommand(quoter stringQuoter, command string, args ...string) string {
	list := make([]string, 0, len(args)+1)

	list = append(list, quoter(command))

	for _, arg := range args {
		list = append(list, quoter(arg))
	}

	return strings.Join(list, " ")
}

func (b *BashWriter) Noticef(format string, arguments ...interface{}) {
	const boldGreen = "\033[32;1m"
	const reset = "\033[0;m"
	coloredText := boldGreen + fmt.Sprintf(format, arguments...) + reset

	b.Line("echo " + b.escape(coloredText))
}

func (b *BashWriter) Cd(path string) {
	b.Command("cd", path)
}

func (b *BashWriter) RmDir(path string) {
	b.Command("rm", "-r", "-f", path)
}
