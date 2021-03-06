package suffuse

import (
  "os/exec"
  "github.com/shirou/gopsutil/process"
  "github.com/shirou/gopsutil/host"
)

func PsutilHostDump () {
  hostInfo, _ := host.HostInfo()
  Println(JsonPretty(*hostInfo))
  // Printfln("%#v", ProcessTable())
}

func ProcessTable() map[int32]string {
  res := make(map[int32]string)
  pids, _ := process.Pids()

  for i := 0 ; i < len(pids) ; i++ {
    pid := pids[i]
    proc, _ := process.NewProcess(pid)
    cline, _ := proc.Cmdline()
    res[pid] = cline
  }
  return res
}

type ExecResult struct {
  Err error
  Stdout []byte
}

func (x ExecResult) OneLine() string { return x.Lines().JoinWords()  }
func (x ExecResult) Lines() Lines    { return BytesToLines(x.Stdout) }
func (x ExecResult) Slurp() string   { return string(x.Stdout)       }
func (x ExecResult) Success() bool   { return x.Err == nil           }

func Exec(args ...string) ExecResult                { return ExecIn(Cwd(), args...)            }
func ExecBash(script string) ExecResult             { return ExecBashIn(Cwd(), script)         }
func ExecBashIn(cwd Path, script string) ExecResult { return ExecIn(cwd, "bash", "-c", script) }

func ExecIn(cwd Path, args ...string) ExecResult {
  var res = ExecResult{}

  if len(args) == 0 {
    res.Err = NewErr("no command")
  } else {
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Dir = string(cwd)

    bytes, err := cmd.Output()

    if err != nil {
      res.Err = err
    } else {
      res.Stdout = bytes
    }
  }

  return res
}

func GitWordDiff(p1 Path, p2 Path) ExecResult {
  return Exec("git", "diff", "--no-index", "--word-diff", string(p1), string(p2))
}
