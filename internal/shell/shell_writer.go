package shells

//go:generate mockery --name=ShellWriter --inpackage
type ShellWriter interface {
	Command(command string, arguments ...string)
	Line(text string)
	CheckForErrors()

	Cd(path string)

	Noticef(format string, arguments ...interface{})

	Finish(trace bool) string
}
