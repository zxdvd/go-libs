package task

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
	"github.com/zxdvd/go-libs/std-helper/str"
)

type fnNewTask func(args ...interface{}) (Task, error)

var registeredTasks = struct {
	constructors map[string]fnNewTask
	lock         sync.Mutex
}{
	constructors: map[string]fnNewTask{},
}

func Register(typ string, newTask func(args ...interface{}) (Task, error)) error {
	registeredTasks.lock.Lock()
	defer registeredTasks.lock.Unlock()
	if _, ok := registeredTasks.constructors[typ]; ok {
		return fmt.Errorf("type %s registered!", typ)
	}
	registeredTasks.constructors[typ] = newTask
	return nil
}

func init() {
	Register("echo", newEchoTask)
	Register("sh", NewShellTask)
	Register("shell", NewShellTask)
	Register("sql", newSqlTask)
	log.Println("echo Register")
}

func newTask(typ string, config map[string]interface{}) (Task, error) {
	registeredTasks.lock.Lock()
	defer registeredTasks.lock.Unlock()
	if new_, ok := registeredTasks.constructors[typ]; ok {
		return new_(config)
	}
	return nil, fmt.Errorf("type %s not registered!", typ)
}

type Task interface {
	Name() string
	Run(context.Context) error
}

type baseTask struct {
	name string
}

func (b baseTask) Name() string {
	return b.name
}

type EchoTask struct {
	baseTask
	str string
}

func newEchoTask(data ...interface{}) (Task, error) {
	conf, ok := data[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to newEchoTask, wrong config")
	}
	return &EchoTask{
		baseTask: baseTask{
			name: conf["name"].(string),
		},
		str: conf["echostr"].(string),
	}, nil
}

func (t *EchoTask) Run(ctx context.Context) error {
	fmt.Println(t.str)
	return nil
}

type ShellTask struct {
	baseTask
	cmd string
	cwd string
}

func NewShellTask(data ...interface{}) (Task, error) {
	conf, ok := data[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to newShellTask, wrong config")
	}
	t := ShellTask{
		baseTask: baseTask{
			name: conf["name"].(string),
		},
		cmd: conf["shellcmd"].(string),
	}
	if cwd, ok := conf["shellcwd"].(string); ok {
		t.cwd = cwd
	}
	return &t, nil
}

func (t *ShellTask) Run(ctx context.Context) error {
	params := map[string]string{
		"name": t.name,
	}
	cmd := str.StrReplace(t.cmd, params)
	log.Println("------cmd", cmd)
	command := exec.Command("sh", "-c", cmd)
	command.Dir = t.cwd
	out, err := command.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

type SqlTask struct {
	baseTask
	dialect string
	uri     string
	Sql     string
}

func newSqlTask(data ...interface{}) (Task, error) {
	conf, ok := data[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to newSqlTask, wrong config")
	}
	return &SqlTask{
		baseTask: baseTask{
			name: conf["name"].(string),
		},
		dialect: conf["dialect"].(string),
		uri:     conf["uri"].(string),
		Sql:     conf["sql"].(string),
	}, nil
}

func (t *SqlTask) Run(ctx context.Context) error {
	db, err := sql.Open(t.dialect, t.uri)
	if err != nil {
		return err
	}
	_, err = db.Exec(t.Sql)
	log.Println("SqlTask error --------", err)
	return err
}
