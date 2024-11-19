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
	// TODO: Document
	AddExec(exec ukcore.Exec, spec ukspec.Parameters, target ...string) error

	// TODO: Document
	AddInfo(info any, target ...string) error

	// TODO: Document
	AddRule(rule ukinit.Rule)

	// TODO: Document
	LoadSpec(t reflect.Type, opts ...ukspec.Option) (ukspec.Parameters, error)

	loadEntry(target []string) (ukexec.Entry, bool)
	runDecode(ukcore.Input, any, ...ukdec.Option) error
	runInit(any, ...ukspec.Option) error
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

func (s *state) AddExec(exec ukcore.Exec, spec ukspec.Parameters, target ...string) error {
	return s.execTree.AddExec(exec, spec, target...)
}

func (s *state) AddInfo(info any, target ...string) error {
	return s.execTree.AddInfo(info, target...)
}

func (s *state) AddRule(rule ukinit.Rule) {
	rule.Register(s.ruleSet)
}

func (s *state) LoadSpec(t reflect.Type, opts ...ukspec.Option) (ukspec.Parameters, error) {
	opts = append(s.config.Spec, opts...)
	return ukspec.NewParameters(t, opts...)
}

func (s *state) loadEntry(target []string) (ukexec.Entry, bool) {
	return s.execTree.LoadEntry(target...)
}

func (s *state) runDecode(input ukcore.Input, v any, opts ...ukdec.Option) error {
	opts = append(s.config.Decode, opts...)
	return ukdec.Decode(input, v, opts...)
}

func (s *state) runInit(v any, opts ...ukspec.Option) error {
	opts = append(s.config.Spec, opts...)
	spec, err := ukspec.ParametersOf(v, opts...)
	if err != nil {
		return err
	}

	return s.ruleSet.Process(spec, v)
}

// =============================================================================
// Context
// =============================================================================

var _ Context = inputContext{}

type Context interface {
	context.Context

	// TODO: Document
	LoadEntry(target ...string) (ukexec.Entry, bool)

	// TODO: Document
	LoadSpec(t reflect.Type, opts ...ukspec.Option) (ukspec.Parameters, error)

	// TODO: Document
	Decode(input ukcore.Input, v any, opts ...ukdec.Option) error

	// TODO: Document
	Initialize(v any, opts ...ukspec.Option) error
}

type inputContext struct {
	context.Context
	state State
}

func newInputContext(ctx context.Context, state State) inputContext {
	return inputContext{Context: ctx, state: state}
}

func (ic inputContext) LoadEntry(target ...string) (ukexec.Entry, bool) {
	return ic.state.loadEntry(target)
}

func (ic inputContext) LoadSpec(t reflect.Type, opts ...ukspec.Option) (ukspec.Parameters, error) {
	return ic.state.LoadSpec(t, opts...)
}

func (ic inputContext) Decode(input ukcore.Input, v any, opts ...ukdec.Option) error {
	return ic.state.runDecode(input, v, opts...)
}

func (ic inputContext) Initialize(v any, opts ...ukspec.Option) error {
	return ic.state.runInit(v, opts...)
}
