package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/olgoncharov/otbook/internal/usecase/access/command/login.passwordChecker -o ./internal/usecase/access/command/login/mocks/password_checker_mock.go -n PasswordCheckerMock

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// PasswordCheckerMock implements login.passwordChecker
type PasswordCheckerMock struct {
	t minimock.Tester

	funcCheck          func(input string, hash string) (b1 bool, err error)
	inspectFuncCheck   func(input string, hash string)
	afterCheckCounter  uint64
	beforeCheckCounter uint64
	CheckMock          mPasswordCheckerMockCheck
}

// NewPasswordCheckerMock returns a mock for login.passwordChecker
func NewPasswordCheckerMock(t minimock.Tester) *PasswordCheckerMock {
	m := &PasswordCheckerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CheckMock = mPasswordCheckerMockCheck{mock: m}
	m.CheckMock.callArgs = []*PasswordCheckerMockCheckParams{}

	return m
}

type mPasswordCheckerMockCheck struct {
	mock               *PasswordCheckerMock
	defaultExpectation *PasswordCheckerMockCheckExpectation
	expectations       []*PasswordCheckerMockCheckExpectation

	callArgs []*PasswordCheckerMockCheckParams
	mutex    sync.RWMutex
}

// PasswordCheckerMockCheckExpectation specifies expectation struct of the passwordChecker.Check
type PasswordCheckerMockCheckExpectation struct {
	mock    *PasswordCheckerMock
	params  *PasswordCheckerMockCheckParams
	results *PasswordCheckerMockCheckResults
	Counter uint64
}

// PasswordCheckerMockCheckParams contains parameters of the passwordChecker.Check
type PasswordCheckerMockCheckParams struct {
	input string
	hash  string
}

// PasswordCheckerMockCheckResults contains results of the passwordChecker.Check
type PasswordCheckerMockCheckResults struct {
	b1  bool
	err error
}

// Expect sets up expected params for passwordChecker.Check
func (mmCheck *mPasswordCheckerMockCheck) Expect(input string, hash string) *mPasswordCheckerMockCheck {
	if mmCheck.mock.funcCheck != nil {
		mmCheck.mock.t.Fatalf("PasswordCheckerMock.Check mock is already set by Set")
	}

	if mmCheck.defaultExpectation == nil {
		mmCheck.defaultExpectation = &PasswordCheckerMockCheckExpectation{}
	}

	mmCheck.defaultExpectation.params = &PasswordCheckerMockCheckParams{input, hash}
	for _, e := range mmCheck.expectations {
		if minimock.Equal(e.params, mmCheck.defaultExpectation.params) {
			mmCheck.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCheck.defaultExpectation.params)
		}
	}

	return mmCheck
}

// Inspect accepts an inspector function that has same arguments as the passwordChecker.Check
func (mmCheck *mPasswordCheckerMockCheck) Inspect(f func(input string, hash string)) *mPasswordCheckerMockCheck {
	if mmCheck.mock.inspectFuncCheck != nil {
		mmCheck.mock.t.Fatalf("Inspect function is already set for PasswordCheckerMock.Check")
	}

	mmCheck.mock.inspectFuncCheck = f

	return mmCheck
}

// Return sets up results that will be returned by passwordChecker.Check
func (mmCheck *mPasswordCheckerMockCheck) Return(b1 bool, err error) *PasswordCheckerMock {
	if mmCheck.mock.funcCheck != nil {
		mmCheck.mock.t.Fatalf("PasswordCheckerMock.Check mock is already set by Set")
	}

	if mmCheck.defaultExpectation == nil {
		mmCheck.defaultExpectation = &PasswordCheckerMockCheckExpectation{mock: mmCheck.mock}
	}
	mmCheck.defaultExpectation.results = &PasswordCheckerMockCheckResults{b1, err}
	return mmCheck.mock
}

//Set uses given function f to mock the passwordChecker.Check method
func (mmCheck *mPasswordCheckerMockCheck) Set(f func(input string, hash string) (b1 bool, err error)) *PasswordCheckerMock {
	if mmCheck.defaultExpectation != nil {
		mmCheck.mock.t.Fatalf("Default expectation is already set for the passwordChecker.Check method")
	}

	if len(mmCheck.expectations) > 0 {
		mmCheck.mock.t.Fatalf("Some expectations are already set for the passwordChecker.Check method")
	}

	mmCheck.mock.funcCheck = f
	return mmCheck.mock
}

// When sets expectation for the passwordChecker.Check which will trigger the result defined by the following
// Then helper
func (mmCheck *mPasswordCheckerMockCheck) When(input string, hash string) *PasswordCheckerMockCheckExpectation {
	if mmCheck.mock.funcCheck != nil {
		mmCheck.mock.t.Fatalf("PasswordCheckerMock.Check mock is already set by Set")
	}

	expectation := &PasswordCheckerMockCheckExpectation{
		mock:   mmCheck.mock,
		params: &PasswordCheckerMockCheckParams{input, hash},
	}
	mmCheck.expectations = append(mmCheck.expectations, expectation)
	return expectation
}

// Then sets up passwordChecker.Check return parameters for the expectation previously defined by the When method
func (e *PasswordCheckerMockCheckExpectation) Then(b1 bool, err error) *PasswordCheckerMock {
	e.results = &PasswordCheckerMockCheckResults{b1, err}
	return e.mock
}

// Check implements login.passwordChecker
func (mmCheck *PasswordCheckerMock) Check(input string, hash string) (b1 bool, err error) {
	mm_atomic.AddUint64(&mmCheck.beforeCheckCounter, 1)
	defer mm_atomic.AddUint64(&mmCheck.afterCheckCounter, 1)

	if mmCheck.inspectFuncCheck != nil {
		mmCheck.inspectFuncCheck(input, hash)
	}

	mm_params := &PasswordCheckerMockCheckParams{input, hash}

	// Record call args
	mmCheck.CheckMock.mutex.Lock()
	mmCheck.CheckMock.callArgs = append(mmCheck.CheckMock.callArgs, mm_params)
	mmCheck.CheckMock.mutex.Unlock()

	for _, e := range mmCheck.CheckMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1, e.results.err
		}
	}

	if mmCheck.CheckMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCheck.CheckMock.defaultExpectation.Counter, 1)
		mm_want := mmCheck.CheckMock.defaultExpectation.params
		mm_got := PasswordCheckerMockCheckParams{input, hash}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCheck.t.Errorf("PasswordCheckerMock.Check got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCheck.CheckMock.defaultExpectation.results
		if mm_results == nil {
			mmCheck.t.Fatal("No results are set for the PasswordCheckerMock.Check")
		}
		return (*mm_results).b1, (*mm_results).err
	}
	if mmCheck.funcCheck != nil {
		return mmCheck.funcCheck(input, hash)
	}
	mmCheck.t.Fatalf("Unexpected call to PasswordCheckerMock.Check. %v %v", input, hash)
	return
}

// CheckAfterCounter returns a count of finished PasswordCheckerMock.Check invocations
func (mmCheck *PasswordCheckerMock) CheckAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCheck.afterCheckCounter)
}

// CheckBeforeCounter returns a count of PasswordCheckerMock.Check invocations
func (mmCheck *PasswordCheckerMock) CheckBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCheck.beforeCheckCounter)
}

// Calls returns a list of arguments used in each call to PasswordCheckerMock.Check.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCheck *mPasswordCheckerMockCheck) Calls() []*PasswordCheckerMockCheckParams {
	mmCheck.mutex.RLock()

	argCopy := make([]*PasswordCheckerMockCheckParams, len(mmCheck.callArgs))
	copy(argCopy, mmCheck.callArgs)

	mmCheck.mutex.RUnlock()

	return argCopy
}

// MinimockCheckDone returns true if the count of the Check invocations corresponds
// the number of defined expectations
func (m *PasswordCheckerMock) MinimockCheckDone() bool {
	for _, e := range m.CheckMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CheckMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCheckCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCheck != nil && mm_atomic.LoadUint64(&m.afterCheckCounter) < 1 {
		return false
	}
	return true
}

// MinimockCheckInspect logs each unmet expectation
func (m *PasswordCheckerMock) MinimockCheckInspect() {
	for _, e := range m.CheckMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PasswordCheckerMock.Check with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CheckMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCheckCounter) < 1 {
		if m.CheckMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PasswordCheckerMock.Check")
		} else {
			m.t.Errorf("Expected call to PasswordCheckerMock.Check with params: %#v", *m.CheckMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCheck != nil && mm_atomic.LoadUint64(&m.afterCheckCounter) < 1 {
		m.t.Error("Expected call to PasswordCheckerMock.Check")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PasswordCheckerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockCheckInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PasswordCheckerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *PasswordCheckerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCheckDone()
}
