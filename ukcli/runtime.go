package ukcli

import (
	"context"
	"reflect"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukdec"
	"github.com/oligarch316/ukase/ukcore/ukexec"
	"github.com/oligarch316/ukase/ukcore/ukinit"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Runtime
// =============================================================================

type Runtime struct {
	config     Config
	directives []Directive
}

func NewRuntime(opts ...Option) *Runtime {
	config := newConfig(opts)
	return &Runtime{config: config}
}

func (r *Runtime) Add(directives ...Directive) {
	r.directives = append(r.directives, directives...)
}

func (r *Runtime) Execute(ctx context.Context, values []string) error {
	state := newState(r.config)

	if err := r.prepare(state); err != nil {
		return err
	}

	return state.execTree.Execute(ctx, values)
}

func (r *Runtime) prepare(state State) error {
	for _, middleware := range r.config.Middleware {
		state = middleware(state)
	}

	for _, dir := range r.directives {
		if err := dir.UkaseRegister(state); err != nil {
			return err
		}
	}

	return nil
}

// =============================================================================
// State
// =============================================================================

var _ State = (*state)(nil)

type State interface {
	// Execution time utilities
	loadEntry(target []string) (ukexec.Entry, bool)
	loadSpec(t reflect.Type) (ukspec.Parameters, error)
	runDecode(ukcore.Input, any) error
	runInit(any) error

	// Registration time utilities
	AddExec(exec ukcore.Exec, spec ukspec.Parameters, target ...string) error
	AddInfo(info any, target ...string) error
	AddRule(rule ukinit.Rule)
}

type state struct {
	config   Config
	execTree *ukexec.Tree
	ruleSet  *ukinit.RuleSet
}

func newState(config Config) *state {
	return &state{
		config:   config,
		execTree: ukexec.NewTree(config.Exec...),
		ruleSet:  ukinit.NewRuleSet(config.Init...),
	}
}

func (s *state) loadEntry(target []string) (ukexec.Entry, bool) {
	return s.execTree.LoadEntry(target...)
}

func (s *state) loadSpec(t reflect.Type) (ukspec.Parameters, error) {
	return ukspec.NewParameters(t, s.config.Spec...)
}

func (s *state) runDecode(i ukcore.Input, v any) error {
	decoder := ukdec.NewDecoder(i, s.config.Decode...)
	return decoder.Decode(v)
}

func (s *state) runInit(v any) error {
	spec, err := ukspec.ParametersOf(v, s.config.Spec...)
	if err != nil {
		return err
	}

	return s.ruleSet.Process(spec, v)
}

func (s *state) AddExec(exec ukcore.Exec, spec ukspec.Parameters, target ...string) error {
	return s.execTree.AddExec(exec, spec, target...)
}

func (s *state) AddInfo(info any, target ...string) error {
	return s.execTree.AddInfo(info, target...)
}

func (s *state) AddRule(rule ukinit.Rule) {
	rule.Register(s.ruleSet)
}

// =============================================================================
// Input
// =============================================================================

var _ Input = input{}

type Input interface {
	Core() ukcore.Input
	Decode(any) error
	Initialize(any) error
	Lookup(target ...string) (ukexec.Entry, bool)
}

type input struct {
	core  ukcore.Input
	state State
}

func newInput(core ukcore.Input, state State) input {
	return input{core: core, state: state}
}

func (i input) Core() ukcore.Input                      { return i.core }
func (i input) Decode(v any) error                      { return i.state.runDecode(i.core, v) }
func (i input) Initialize(v any) error                  { return i.state.runInit(v) }
func (i input) Lookup(t ...string) (ukexec.Entry, bool) { return i.state.loadEntry(t) }
