package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/zxdvd/go-libs/dag"
	"github.com/zxdvd/go-libs/future"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	defaultLogger, _ := zap.NewDevelopment()
	logger = defaultLogger
}

type TaskHook func(*task) error

type task struct {
	Task
	preRunHooks  []TaskHook
	postRunHooks []TaskHook
	dependOn     []*task
	m            sync.Mutex
	done         bool
	pool         *RunnerPool
}

func NewTask(typ string, t map[string]interface{}) (*task, error) {
	t1, err := newTask(typ, t)
	if err != nil {
		return nil, err
	}
	if t1 == nil {
		return nil, nil
	}
	return &task{
		Task: t1,
	}, nil
}

func (t *task) Nexts() []dag.Node {
	nodes := make([]dag.Node, len(t.dependOn))
	for i, dep := range t.dependOn {
		nodes[i] = dep
	}
	return nodes
}

func (t *task) run(ctx context.Context) error {
	if t.done {
		return nil
	}
	if t.pool != nil {
		t.pool.Get()
		defer t.pool.Put()
	}
	t.m.Lock()
	defer t.m.Unlock()
	if t.done {
		return nil
	}
	logger.Debug("run task", zap.String("name", t.Name()))
	err := t.Task.Run(ctx)
	logger.Debug("run task finished", zap.Error(err), zap.String("name", t.Name()))
	t.done = true
	return err
}

func (t *task) Run(ctx context.Context) error {
	futures := future.NewN(len(t.dependOn))
	for i, dep := range t.dependOn {
		go func(i int, t *task) {
			err := t.Run(ctx)
			if err != nil {
				futures[i].SetError(err)
			} else {
				futures[i].SetResult(true)
			}
		}(i, dep)
	}
	// wait and resolve all depends tasks
	if _, err := future.GetAll(futures); err != nil {
		return err
	}
	for _, fn := range t.preRunHooks {
		if err := fn(t); err != nil {
			return errors.Wrap(err, "preRunHooks fails")
		}
	}
	if err := t.run(ctx); err != nil {
		return err
	}
	for _, fn := range t.postRunHooks {
		if err := fn(t); err != nil {
			return errors.Wrap(err, "postRunHooks fails")
		}
	}
	return nil
}

var _ dag.Node = &task{}

var defaultConcurrentLimit = 3

type DagTaskConfig struct {
	Tasks           []map[string]interface{}
	ConcurrentLimit int
}

func CreateTaskDag(c DagTaskConfig) (*dagTask, error) {
	if c.ConcurrentLimit == 0 {
		c.ConcurrentLimit = defaultConcurrentLimit
	}
	taskmap := map[string]*task{}
	for _, tc := range c.Tasks {
		t, err := NewTask(tc["type"].(string), tc)
		if err != nil {
			return nil, err
		}
		if t == nil {
			continue
		}
		taskmap[t.Name()] = t
	}
	// deal with task depends
	for _, tc := range c.Tasks {
		taskname, ok := tc["name"].(string)
		if !ok {
			continue
		}
		t := taskmap[taskname]
		if tc["dependOn"] == nil {
			continue
		}
		dependOn := tc["dependOn"].([]interface{})
		depends := make([]*task, 0, len(dependOn))
		for _, depend := range dependOn {
			dep, ok := depend.(string)
			if !ok {
				continue
			}
			if deptask, ok := taskmap[dep]; ok {
				depends = append(depends, deptask)
			} else {
				return nil, fmt.Errorf("depended task %s not found", dep)
			}
		}
		t.dependOn = depends
	}
	dag_ := &dag.Dag{}
	for _, t := range taskmap {
		dag_.Add(t)
	}
	if err := dag_.CircleDetect(); err != nil {
		return nil, err
	}
	return &dagTask{
		Dag:  dag_,
		pool: NewRunnerPool(defaultConcurrentLimit),
	}, nil
}

type dagTask struct {
	*dag.Dag
	pool *RunnerPool
}

type RunnerPool struct {
	limit int
	pool  chan struct{}
}

func NewRunnerPool(n int) *RunnerPool {
	return &RunnerPool{
		limit: n,
		pool:  make(chan struct{}, n),
	}
}

func (p *RunnerPool) Get() {
	p.pool <- struct{}{}
}

func (p *RunnerPool) Put() {
	<-p.pool
}

func (d *dagTask) RunTask(name string) error {
	ctx, _ := context.WithCancel(context.Background())
	nodes := d.Nodes()
	for _, node := range nodes {
		t, ok := node.(*task)
		if !ok {
			panic("task not implement node")
		}
		if t.Name() != name {
			continue
		}
		return t.Run(ctx)
	}
	return nil
}

func (d *dagTask) Run() error {
	defer logger.Sync()
	ctx, cancel := context.WithCancel(context.Background())
	nodes := d.Nodes()
	futures := future.NewN(len(nodes))
	for i, node := range nodes {
		t, ok := node.(*task)
		if !ok {
			panic("task not implement node")
		}
		t.pool = d.pool
		go func(i int) {
			if err := t.Run(ctx); err != nil {
				futures[i].SetError(err)
			} else {
				futures[i].SetResult(true)
			}
		}(i)
	}
	if _, err := future.GetAll(futures); err != nil {
		cancel()
		logger.Debug("error:",
			zap.Error(err), zap.Stack("stack"))
		return err
	}
	return nil
}
