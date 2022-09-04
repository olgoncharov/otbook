package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/olgoncharov/otbook/internal/usecase/profile/query/search.profileRepo -o ./internal/usecase/profile/query/search/mocks/profile_repo_mock.go -n ProfileRepoMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/entity"
)

// ProfileRepoMock implements search.profileRepo
type ProfileRepoMock struct {
	t minimock.Tester

	funcSearchProfiles          func(ctx context.Context, firstNamePrefix string, lastNamePrefix string) (pa1 []entity.Profile, err error)
	inspectFuncSearchProfiles   func(ctx context.Context, firstNamePrefix string, lastNamePrefix string)
	afterSearchProfilesCounter  uint64
	beforeSearchProfilesCounter uint64
	SearchProfilesMock          mProfileRepoMockSearchProfiles
}

// NewProfileRepoMock returns a mock for search.profileRepo
func NewProfileRepoMock(t minimock.Tester) *ProfileRepoMock {
	m := &ProfileRepoMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SearchProfilesMock = mProfileRepoMockSearchProfiles{mock: m}
	m.SearchProfilesMock.callArgs = []*ProfileRepoMockSearchProfilesParams{}

	return m
}

type mProfileRepoMockSearchProfiles struct {
	mock               *ProfileRepoMock
	defaultExpectation *ProfileRepoMockSearchProfilesExpectation
	expectations       []*ProfileRepoMockSearchProfilesExpectation

	callArgs []*ProfileRepoMockSearchProfilesParams
	mutex    sync.RWMutex
}

// ProfileRepoMockSearchProfilesExpectation specifies expectation struct of the profileRepo.SearchProfiles
type ProfileRepoMockSearchProfilesExpectation struct {
	mock    *ProfileRepoMock
	params  *ProfileRepoMockSearchProfilesParams
	results *ProfileRepoMockSearchProfilesResults
	Counter uint64
}

// ProfileRepoMockSearchProfilesParams contains parameters of the profileRepo.SearchProfiles
type ProfileRepoMockSearchProfilesParams struct {
	ctx             context.Context
	firstNamePrefix string
	lastNamePrefix  string
}

// ProfileRepoMockSearchProfilesResults contains results of the profileRepo.SearchProfiles
type ProfileRepoMockSearchProfilesResults struct {
	pa1 []entity.Profile
	err error
}

// Expect sets up expected params for profileRepo.SearchProfiles
func (mmSearchProfiles *mProfileRepoMockSearchProfiles) Expect(ctx context.Context, firstNamePrefix string, lastNamePrefix string) *mProfileRepoMockSearchProfiles {
	if mmSearchProfiles.mock.funcSearchProfiles != nil {
		mmSearchProfiles.mock.t.Fatalf("ProfileRepoMock.SearchProfiles mock is already set by Set")
	}

	if mmSearchProfiles.defaultExpectation == nil {
		mmSearchProfiles.defaultExpectation = &ProfileRepoMockSearchProfilesExpectation{}
	}

	mmSearchProfiles.defaultExpectation.params = &ProfileRepoMockSearchProfilesParams{ctx, firstNamePrefix, lastNamePrefix}
	for _, e := range mmSearchProfiles.expectations {
		if minimock.Equal(e.params, mmSearchProfiles.defaultExpectation.params) {
			mmSearchProfiles.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSearchProfiles.defaultExpectation.params)
		}
	}

	return mmSearchProfiles
}

// Inspect accepts an inspector function that has same arguments as the profileRepo.SearchProfiles
func (mmSearchProfiles *mProfileRepoMockSearchProfiles) Inspect(f func(ctx context.Context, firstNamePrefix string, lastNamePrefix string)) *mProfileRepoMockSearchProfiles {
	if mmSearchProfiles.mock.inspectFuncSearchProfiles != nil {
		mmSearchProfiles.mock.t.Fatalf("Inspect function is already set for ProfileRepoMock.SearchProfiles")
	}

	mmSearchProfiles.mock.inspectFuncSearchProfiles = f

	return mmSearchProfiles
}

// Return sets up results that will be returned by profileRepo.SearchProfiles
func (mmSearchProfiles *mProfileRepoMockSearchProfiles) Return(pa1 []entity.Profile, err error) *ProfileRepoMock {
	if mmSearchProfiles.mock.funcSearchProfiles != nil {
		mmSearchProfiles.mock.t.Fatalf("ProfileRepoMock.SearchProfiles mock is already set by Set")
	}

	if mmSearchProfiles.defaultExpectation == nil {
		mmSearchProfiles.defaultExpectation = &ProfileRepoMockSearchProfilesExpectation{mock: mmSearchProfiles.mock}
	}
	mmSearchProfiles.defaultExpectation.results = &ProfileRepoMockSearchProfilesResults{pa1, err}
	return mmSearchProfiles.mock
}

//Set uses given function f to mock the profileRepo.SearchProfiles method
func (mmSearchProfiles *mProfileRepoMockSearchProfiles) Set(f func(ctx context.Context, firstNamePrefix string, lastNamePrefix string) (pa1 []entity.Profile, err error)) *ProfileRepoMock {
	if mmSearchProfiles.defaultExpectation != nil {
		mmSearchProfiles.mock.t.Fatalf("Default expectation is already set for the profileRepo.SearchProfiles method")
	}

	if len(mmSearchProfiles.expectations) > 0 {
		mmSearchProfiles.mock.t.Fatalf("Some expectations are already set for the profileRepo.SearchProfiles method")
	}

	mmSearchProfiles.mock.funcSearchProfiles = f
	return mmSearchProfiles.mock
}

// When sets expectation for the profileRepo.SearchProfiles which will trigger the result defined by the following
// Then helper
func (mmSearchProfiles *mProfileRepoMockSearchProfiles) When(ctx context.Context, firstNamePrefix string, lastNamePrefix string) *ProfileRepoMockSearchProfilesExpectation {
	if mmSearchProfiles.mock.funcSearchProfiles != nil {
		mmSearchProfiles.mock.t.Fatalf("ProfileRepoMock.SearchProfiles mock is already set by Set")
	}

	expectation := &ProfileRepoMockSearchProfilesExpectation{
		mock:   mmSearchProfiles.mock,
		params: &ProfileRepoMockSearchProfilesParams{ctx, firstNamePrefix, lastNamePrefix},
	}
	mmSearchProfiles.expectations = append(mmSearchProfiles.expectations, expectation)
	return expectation
}

// Then sets up profileRepo.SearchProfiles return parameters for the expectation previously defined by the When method
func (e *ProfileRepoMockSearchProfilesExpectation) Then(pa1 []entity.Profile, err error) *ProfileRepoMock {
	e.results = &ProfileRepoMockSearchProfilesResults{pa1, err}
	return e.mock
}

// SearchProfiles implements search.profileRepo
func (mmSearchProfiles *ProfileRepoMock) SearchProfiles(ctx context.Context, firstNamePrefix string, lastNamePrefix string) (pa1 []entity.Profile, err error) {
	mm_atomic.AddUint64(&mmSearchProfiles.beforeSearchProfilesCounter, 1)
	defer mm_atomic.AddUint64(&mmSearchProfiles.afterSearchProfilesCounter, 1)

	if mmSearchProfiles.inspectFuncSearchProfiles != nil {
		mmSearchProfiles.inspectFuncSearchProfiles(ctx, firstNamePrefix, lastNamePrefix)
	}

	mm_params := &ProfileRepoMockSearchProfilesParams{ctx, firstNamePrefix, lastNamePrefix}

	// Record call args
	mmSearchProfiles.SearchProfilesMock.mutex.Lock()
	mmSearchProfiles.SearchProfilesMock.callArgs = append(mmSearchProfiles.SearchProfilesMock.callArgs, mm_params)
	mmSearchProfiles.SearchProfilesMock.mutex.Unlock()

	for _, e := range mmSearchProfiles.SearchProfilesMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.pa1, e.results.err
		}
	}

	if mmSearchProfiles.SearchProfilesMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSearchProfiles.SearchProfilesMock.defaultExpectation.Counter, 1)
		mm_want := mmSearchProfiles.SearchProfilesMock.defaultExpectation.params
		mm_got := ProfileRepoMockSearchProfilesParams{ctx, firstNamePrefix, lastNamePrefix}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSearchProfiles.t.Errorf("ProfileRepoMock.SearchProfiles got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSearchProfiles.SearchProfilesMock.defaultExpectation.results
		if mm_results == nil {
			mmSearchProfiles.t.Fatal("No results are set for the ProfileRepoMock.SearchProfiles")
		}
		return (*mm_results).pa1, (*mm_results).err
	}
	if mmSearchProfiles.funcSearchProfiles != nil {
		return mmSearchProfiles.funcSearchProfiles(ctx, firstNamePrefix, lastNamePrefix)
	}
	mmSearchProfiles.t.Fatalf("Unexpected call to ProfileRepoMock.SearchProfiles. %v %v %v", ctx, firstNamePrefix, lastNamePrefix)
	return
}

// SearchProfilesAfterCounter returns a count of finished ProfileRepoMock.SearchProfiles invocations
func (mmSearchProfiles *ProfileRepoMock) SearchProfilesAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSearchProfiles.afterSearchProfilesCounter)
}

// SearchProfilesBeforeCounter returns a count of ProfileRepoMock.SearchProfiles invocations
func (mmSearchProfiles *ProfileRepoMock) SearchProfilesBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSearchProfiles.beforeSearchProfilesCounter)
}

// Calls returns a list of arguments used in each call to ProfileRepoMock.SearchProfiles.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSearchProfiles *mProfileRepoMockSearchProfiles) Calls() []*ProfileRepoMockSearchProfilesParams {
	mmSearchProfiles.mutex.RLock()

	argCopy := make([]*ProfileRepoMockSearchProfilesParams, len(mmSearchProfiles.callArgs))
	copy(argCopy, mmSearchProfiles.callArgs)

	mmSearchProfiles.mutex.RUnlock()

	return argCopy
}

// MinimockSearchProfilesDone returns true if the count of the SearchProfiles invocations corresponds
// the number of defined expectations
func (m *ProfileRepoMock) MinimockSearchProfilesDone() bool {
	for _, e := range m.SearchProfilesMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SearchProfilesMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSearchProfilesCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSearchProfiles != nil && mm_atomic.LoadUint64(&m.afterSearchProfilesCounter) < 1 {
		return false
	}
	return true
}

// MinimockSearchProfilesInspect logs each unmet expectation
func (m *ProfileRepoMock) MinimockSearchProfilesInspect() {
	for _, e := range m.SearchProfilesMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ProfileRepoMock.SearchProfiles with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SearchProfilesMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSearchProfilesCounter) < 1 {
		if m.SearchProfilesMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ProfileRepoMock.SearchProfiles")
		} else {
			m.t.Errorf("Expected call to ProfileRepoMock.SearchProfiles with params: %#v", *m.SearchProfilesMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSearchProfiles != nil && mm_atomic.LoadUint64(&m.afterSearchProfilesCounter) < 1 {
		m.t.Error("Expected call to ProfileRepoMock.SearchProfiles")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ProfileRepoMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockSearchProfilesInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ProfileRepoMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ProfileRepoMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockSearchProfilesDone()
}
