package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/olgoncharov/otbook/internal/usecase/post/query/feed.cacheRepo -o ./internal/usecase/post/query/feed/mocks/cache_repo_mock.go -n CacheRepoMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/olgoncharov/otbook/internal/repository/dto"
)

// CacheRepoMock implements feed.cacheRepo
type CacheRepoMock struct {
	t minimock.Tester

	funcGetCelebrityFriends          func(ctx context.Context, username string) (sa1 []string, err error)
	inspectFuncGetCelebrityFriends   func(ctx context.Context, username string)
	afterGetCelebrityFriendsCounter  uint64
	beforeGetCelebrityFriendsCounter uint64
	GetCelebrityFriendsMock          mCacheRepoMockGetCelebrityFriends

	funcGetPostFeed          func(ctx context.Context, username string, limit uint) (pa1 []dto.PostShortInfo, err error)
	inspectFuncGetPostFeed   func(ctx context.Context, username string, limit uint)
	afterGetPostFeedCounter  uint64
	beforeGetPostFeedCounter uint64
	GetPostFeedMock          mCacheRepoMockGetPostFeed
}

// NewCacheRepoMock returns a mock for feed.cacheRepo
func NewCacheRepoMock(t minimock.Tester) *CacheRepoMock {
	m := &CacheRepoMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCelebrityFriendsMock = mCacheRepoMockGetCelebrityFriends{mock: m}
	m.GetCelebrityFriendsMock.callArgs = []*CacheRepoMockGetCelebrityFriendsParams{}

	m.GetPostFeedMock = mCacheRepoMockGetPostFeed{mock: m}
	m.GetPostFeedMock.callArgs = []*CacheRepoMockGetPostFeedParams{}

	return m
}

type mCacheRepoMockGetCelebrityFriends struct {
	mock               *CacheRepoMock
	defaultExpectation *CacheRepoMockGetCelebrityFriendsExpectation
	expectations       []*CacheRepoMockGetCelebrityFriendsExpectation

	callArgs []*CacheRepoMockGetCelebrityFriendsParams
	mutex    sync.RWMutex
}

// CacheRepoMockGetCelebrityFriendsExpectation specifies expectation struct of the cacheRepo.GetCelebrityFriends
type CacheRepoMockGetCelebrityFriendsExpectation struct {
	mock    *CacheRepoMock
	params  *CacheRepoMockGetCelebrityFriendsParams
	results *CacheRepoMockGetCelebrityFriendsResults
	Counter uint64
}

// CacheRepoMockGetCelebrityFriendsParams contains parameters of the cacheRepo.GetCelebrityFriends
type CacheRepoMockGetCelebrityFriendsParams struct {
	ctx      context.Context
	username string
}

// CacheRepoMockGetCelebrityFriendsResults contains results of the cacheRepo.GetCelebrityFriends
type CacheRepoMockGetCelebrityFriendsResults struct {
	sa1 []string
	err error
}

// Expect sets up expected params for cacheRepo.GetCelebrityFriends
func (mmGetCelebrityFriends *mCacheRepoMockGetCelebrityFriends) Expect(ctx context.Context, username string) *mCacheRepoMockGetCelebrityFriends {
	if mmGetCelebrityFriends.mock.funcGetCelebrityFriends != nil {
		mmGetCelebrityFriends.mock.t.Fatalf("CacheRepoMock.GetCelebrityFriends mock is already set by Set")
	}

	if mmGetCelebrityFriends.defaultExpectation == nil {
		mmGetCelebrityFriends.defaultExpectation = &CacheRepoMockGetCelebrityFriendsExpectation{}
	}

	mmGetCelebrityFriends.defaultExpectation.params = &CacheRepoMockGetCelebrityFriendsParams{ctx, username}
	for _, e := range mmGetCelebrityFriends.expectations {
		if minimock.Equal(e.params, mmGetCelebrityFriends.defaultExpectation.params) {
			mmGetCelebrityFriends.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetCelebrityFriends.defaultExpectation.params)
		}
	}

	return mmGetCelebrityFriends
}

// Inspect accepts an inspector function that has same arguments as the cacheRepo.GetCelebrityFriends
func (mmGetCelebrityFriends *mCacheRepoMockGetCelebrityFriends) Inspect(f func(ctx context.Context, username string)) *mCacheRepoMockGetCelebrityFriends {
	if mmGetCelebrityFriends.mock.inspectFuncGetCelebrityFriends != nil {
		mmGetCelebrityFriends.mock.t.Fatalf("Inspect function is already set for CacheRepoMock.GetCelebrityFriends")
	}

	mmGetCelebrityFriends.mock.inspectFuncGetCelebrityFriends = f

	return mmGetCelebrityFriends
}

// Return sets up results that will be returned by cacheRepo.GetCelebrityFriends
func (mmGetCelebrityFriends *mCacheRepoMockGetCelebrityFriends) Return(sa1 []string, err error) *CacheRepoMock {
	if mmGetCelebrityFriends.mock.funcGetCelebrityFriends != nil {
		mmGetCelebrityFriends.mock.t.Fatalf("CacheRepoMock.GetCelebrityFriends mock is already set by Set")
	}

	if mmGetCelebrityFriends.defaultExpectation == nil {
		mmGetCelebrityFriends.defaultExpectation = &CacheRepoMockGetCelebrityFriendsExpectation{mock: mmGetCelebrityFriends.mock}
	}
	mmGetCelebrityFriends.defaultExpectation.results = &CacheRepoMockGetCelebrityFriendsResults{sa1, err}
	return mmGetCelebrityFriends.mock
}

//Set uses given function f to mock the cacheRepo.GetCelebrityFriends method
func (mmGetCelebrityFriends *mCacheRepoMockGetCelebrityFriends) Set(f func(ctx context.Context, username string) (sa1 []string, err error)) *CacheRepoMock {
	if mmGetCelebrityFriends.defaultExpectation != nil {
		mmGetCelebrityFriends.mock.t.Fatalf("Default expectation is already set for the cacheRepo.GetCelebrityFriends method")
	}

	if len(mmGetCelebrityFriends.expectations) > 0 {
		mmGetCelebrityFriends.mock.t.Fatalf("Some expectations are already set for the cacheRepo.GetCelebrityFriends method")
	}

	mmGetCelebrityFriends.mock.funcGetCelebrityFriends = f
	return mmGetCelebrityFriends.mock
}

// When sets expectation for the cacheRepo.GetCelebrityFriends which will trigger the result defined by the following
// Then helper
func (mmGetCelebrityFriends *mCacheRepoMockGetCelebrityFriends) When(ctx context.Context, username string) *CacheRepoMockGetCelebrityFriendsExpectation {
	if mmGetCelebrityFriends.mock.funcGetCelebrityFriends != nil {
		mmGetCelebrityFriends.mock.t.Fatalf("CacheRepoMock.GetCelebrityFriends mock is already set by Set")
	}

	expectation := &CacheRepoMockGetCelebrityFriendsExpectation{
		mock:   mmGetCelebrityFriends.mock,
		params: &CacheRepoMockGetCelebrityFriendsParams{ctx, username},
	}
	mmGetCelebrityFriends.expectations = append(mmGetCelebrityFriends.expectations, expectation)
	return expectation
}

// Then sets up cacheRepo.GetCelebrityFriends return parameters for the expectation previously defined by the When method
func (e *CacheRepoMockGetCelebrityFriendsExpectation) Then(sa1 []string, err error) *CacheRepoMock {
	e.results = &CacheRepoMockGetCelebrityFriendsResults{sa1, err}
	return e.mock
}

// GetCelebrityFriends implements feed.cacheRepo
func (mmGetCelebrityFriends *CacheRepoMock) GetCelebrityFriends(ctx context.Context, username string) (sa1 []string, err error) {
	mm_atomic.AddUint64(&mmGetCelebrityFriends.beforeGetCelebrityFriendsCounter, 1)
	defer mm_atomic.AddUint64(&mmGetCelebrityFriends.afterGetCelebrityFriendsCounter, 1)

	if mmGetCelebrityFriends.inspectFuncGetCelebrityFriends != nil {
		mmGetCelebrityFriends.inspectFuncGetCelebrityFriends(ctx, username)
	}

	mm_params := &CacheRepoMockGetCelebrityFriendsParams{ctx, username}

	// Record call args
	mmGetCelebrityFriends.GetCelebrityFriendsMock.mutex.Lock()
	mmGetCelebrityFriends.GetCelebrityFriendsMock.callArgs = append(mmGetCelebrityFriends.GetCelebrityFriendsMock.callArgs, mm_params)
	mmGetCelebrityFriends.GetCelebrityFriendsMock.mutex.Unlock()

	for _, e := range mmGetCelebrityFriends.GetCelebrityFriendsMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.sa1, e.results.err
		}
	}

	if mmGetCelebrityFriends.GetCelebrityFriendsMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetCelebrityFriends.GetCelebrityFriendsMock.defaultExpectation.Counter, 1)
		mm_want := mmGetCelebrityFriends.GetCelebrityFriendsMock.defaultExpectation.params
		mm_got := CacheRepoMockGetCelebrityFriendsParams{ctx, username}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetCelebrityFriends.t.Errorf("CacheRepoMock.GetCelebrityFriends got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetCelebrityFriends.GetCelebrityFriendsMock.defaultExpectation.results
		if mm_results == nil {
			mmGetCelebrityFriends.t.Fatal("No results are set for the CacheRepoMock.GetCelebrityFriends")
		}
		return (*mm_results).sa1, (*mm_results).err
	}
	if mmGetCelebrityFriends.funcGetCelebrityFriends != nil {
		return mmGetCelebrityFriends.funcGetCelebrityFriends(ctx, username)
	}
	mmGetCelebrityFriends.t.Fatalf("Unexpected call to CacheRepoMock.GetCelebrityFriends. %v %v", ctx, username)
	return
}

// GetCelebrityFriendsAfterCounter returns a count of finished CacheRepoMock.GetCelebrityFriends invocations
func (mmGetCelebrityFriends *CacheRepoMock) GetCelebrityFriendsAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCelebrityFriends.afterGetCelebrityFriendsCounter)
}

// GetCelebrityFriendsBeforeCounter returns a count of CacheRepoMock.GetCelebrityFriends invocations
func (mmGetCelebrityFriends *CacheRepoMock) GetCelebrityFriendsBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCelebrityFriends.beforeGetCelebrityFriendsCounter)
}

// Calls returns a list of arguments used in each call to CacheRepoMock.GetCelebrityFriends.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetCelebrityFriends *mCacheRepoMockGetCelebrityFriends) Calls() []*CacheRepoMockGetCelebrityFriendsParams {
	mmGetCelebrityFriends.mutex.RLock()

	argCopy := make([]*CacheRepoMockGetCelebrityFriendsParams, len(mmGetCelebrityFriends.callArgs))
	copy(argCopy, mmGetCelebrityFriends.callArgs)

	mmGetCelebrityFriends.mutex.RUnlock()

	return argCopy
}

// MinimockGetCelebrityFriendsDone returns true if the count of the GetCelebrityFriends invocations corresponds
// the number of defined expectations
func (m *CacheRepoMock) MinimockGetCelebrityFriendsDone() bool {
	for _, e := range m.GetCelebrityFriendsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCelebrityFriendsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCelebrityFriendsCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCelebrityFriends != nil && mm_atomic.LoadUint64(&m.afterGetCelebrityFriendsCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetCelebrityFriendsInspect logs each unmet expectation
func (m *CacheRepoMock) MinimockGetCelebrityFriendsInspect() {
	for _, e := range m.GetCelebrityFriendsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheRepoMock.GetCelebrityFriends with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCelebrityFriendsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCelebrityFriendsCounter) < 1 {
		if m.GetCelebrityFriendsMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CacheRepoMock.GetCelebrityFriends")
		} else {
			m.t.Errorf("Expected call to CacheRepoMock.GetCelebrityFriends with params: %#v", *m.GetCelebrityFriendsMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCelebrityFriends != nil && mm_atomic.LoadUint64(&m.afterGetCelebrityFriendsCounter) < 1 {
		m.t.Error("Expected call to CacheRepoMock.GetCelebrityFriends")
	}
}

type mCacheRepoMockGetPostFeed struct {
	mock               *CacheRepoMock
	defaultExpectation *CacheRepoMockGetPostFeedExpectation
	expectations       []*CacheRepoMockGetPostFeedExpectation

	callArgs []*CacheRepoMockGetPostFeedParams
	mutex    sync.RWMutex
}

// CacheRepoMockGetPostFeedExpectation specifies expectation struct of the cacheRepo.GetPostFeed
type CacheRepoMockGetPostFeedExpectation struct {
	mock    *CacheRepoMock
	params  *CacheRepoMockGetPostFeedParams
	results *CacheRepoMockGetPostFeedResults
	Counter uint64
}

// CacheRepoMockGetPostFeedParams contains parameters of the cacheRepo.GetPostFeed
type CacheRepoMockGetPostFeedParams struct {
	ctx      context.Context
	username string
	limit    uint
}

// CacheRepoMockGetPostFeedResults contains results of the cacheRepo.GetPostFeed
type CacheRepoMockGetPostFeedResults struct {
	pa1 []dto.PostShortInfo
	err error
}

// Expect sets up expected params for cacheRepo.GetPostFeed
func (mmGetPostFeed *mCacheRepoMockGetPostFeed) Expect(ctx context.Context, username string, limit uint) *mCacheRepoMockGetPostFeed {
	if mmGetPostFeed.mock.funcGetPostFeed != nil {
		mmGetPostFeed.mock.t.Fatalf("CacheRepoMock.GetPostFeed mock is already set by Set")
	}

	if mmGetPostFeed.defaultExpectation == nil {
		mmGetPostFeed.defaultExpectation = &CacheRepoMockGetPostFeedExpectation{}
	}

	mmGetPostFeed.defaultExpectation.params = &CacheRepoMockGetPostFeedParams{ctx, username, limit}
	for _, e := range mmGetPostFeed.expectations {
		if minimock.Equal(e.params, mmGetPostFeed.defaultExpectation.params) {
			mmGetPostFeed.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetPostFeed.defaultExpectation.params)
		}
	}

	return mmGetPostFeed
}

// Inspect accepts an inspector function that has same arguments as the cacheRepo.GetPostFeed
func (mmGetPostFeed *mCacheRepoMockGetPostFeed) Inspect(f func(ctx context.Context, username string, limit uint)) *mCacheRepoMockGetPostFeed {
	if mmGetPostFeed.mock.inspectFuncGetPostFeed != nil {
		mmGetPostFeed.mock.t.Fatalf("Inspect function is already set for CacheRepoMock.GetPostFeed")
	}

	mmGetPostFeed.mock.inspectFuncGetPostFeed = f

	return mmGetPostFeed
}

// Return sets up results that will be returned by cacheRepo.GetPostFeed
func (mmGetPostFeed *mCacheRepoMockGetPostFeed) Return(pa1 []dto.PostShortInfo, err error) *CacheRepoMock {
	if mmGetPostFeed.mock.funcGetPostFeed != nil {
		mmGetPostFeed.mock.t.Fatalf("CacheRepoMock.GetPostFeed mock is already set by Set")
	}

	if mmGetPostFeed.defaultExpectation == nil {
		mmGetPostFeed.defaultExpectation = &CacheRepoMockGetPostFeedExpectation{mock: mmGetPostFeed.mock}
	}
	mmGetPostFeed.defaultExpectation.results = &CacheRepoMockGetPostFeedResults{pa1, err}
	return mmGetPostFeed.mock
}

//Set uses given function f to mock the cacheRepo.GetPostFeed method
func (mmGetPostFeed *mCacheRepoMockGetPostFeed) Set(f func(ctx context.Context, username string, limit uint) (pa1 []dto.PostShortInfo, err error)) *CacheRepoMock {
	if mmGetPostFeed.defaultExpectation != nil {
		mmGetPostFeed.mock.t.Fatalf("Default expectation is already set for the cacheRepo.GetPostFeed method")
	}

	if len(mmGetPostFeed.expectations) > 0 {
		mmGetPostFeed.mock.t.Fatalf("Some expectations are already set for the cacheRepo.GetPostFeed method")
	}

	mmGetPostFeed.mock.funcGetPostFeed = f
	return mmGetPostFeed.mock
}

// When sets expectation for the cacheRepo.GetPostFeed which will trigger the result defined by the following
// Then helper
func (mmGetPostFeed *mCacheRepoMockGetPostFeed) When(ctx context.Context, username string, limit uint) *CacheRepoMockGetPostFeedExpectation {
	if mmGetPostFeed.mock.funcGetPostFeed != nil {
		mmGetPostFeed.mock.t.Fatalf("CacheRepoMock.GetPostFeed mock is already set by Set")
	}

	expectation := &CacheRepoMockGetPostFeedExpectation{
		mock:   mmGetPostFeed.mock,
		params: &CacheRepoMockGetPostFeedParams{ctx, username, limit},
	}
	mmGetPostFeed.expectations = append(mmGetPostFeed.expectations, expectation)
	return expectation
}

// Then sets up cacheRepo.GetPostFeed return parameters for the expectation previously defined by the When method
func (e *CacheRepoMockGetPostFeedExpectation) Then(pa1 []dto.PostShortInfo, err error) *CacheRepoMock {
	e.results = &CacheRepoMockGetPostFeedResults{pa1, err}
	return e.mock
}

// GetPostFeed implements feed.cacheRepo
func (mmGetPostFeed *CacheRepoMock) GetPostFeed(ctx context.Context, username string, limit uint) (pa1 []dto.PostShortInfo, err error) {
	mm_atomic.AddUint64(&mmGetPostFeed.beforeGetPostFeedCounter, 1)
	defer mm_atomic.AddUint64(&mmGetPostFeed.afterGetPostFeedCounter, 1)

	if mmGetPostFeed.inspectFuncGetPostFeed != nil {
		mmGetPostFeed.inspectFuncGetPostFeed(ctx, username, limit)
	}

	mm_params := &CacheRepoMockGetPostFeedParams{ctx, username, limit}

	// Record call args
	mmGetPostFeed.GetPostFeedMock.mutex.Lock()
	mmGetPostFeed.GetPostFeedMock.callArgs = append(mmGetPostFeed.GetPostFeedMock.callArgs, mm_params)
	mmGetPostFeed.GetPostFeedMock.mutex.Unlock()

	for _, e := range mmGetPostFeed.GetPostFeedMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.pa1, e.results.err
		}
	}

	if mmGetPostFeed.GetPostFeedMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetPostFeed.GetPostFeedMock.defaultExpectation.Counter, 1)
		mm_want := mmGetPostFeed.GetPostFeedMock.defaultExpectation.params
		mm_got := CacheRepoMockGetPostFeedParams{ctx, username, limit}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetPostFeed.t.Errorf("CacheRepoMock.GetPostFeed got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetPostFeed.GetPostFeedMock.defaultExpectation.results
		if mm_results == nil {
			mmGetPostFeed.t.Fatal("No results are set for the CacheRepoMock.GetPostFeed")
		}
		return (*mm_results).pa1, (*mm_results).err
	}
	if mmGetPostFeed.funcGetPostFeed != nil {
		return mmGetPostFeed.funcGetPostFeed(ctx, username, limit)
	}
	mmGetPostFeed.t.Fatalf("Unexpected call to CacheRepoMock.GetPostFeed. %v %v %v", ctx, username, limit)
	return
}

// GetPostFeedAfterCounter returns a count of finished CacheRepoMock.GetPostFeed invocations
func (mmGetPostFeed *CacheRepoMock) GetPostFeedAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetPostFeed.afterGetPostFeedCounter)
}

// GetPostFeedBeforeCounter returns a count of CacheRepoMock.GetPostFeed invocations
func (mmGetPostFeed *CacheRepoMock) GetPostFeedBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetPostFeed.beforeGetPostFeedCounter)
}

// Calls returns a list of arguments used in each call to CacheRepoMock.GetPostFeed.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetPostFeed *mCacheRepoMockGetPostFeed) Calls() []*CacheRepoMockGetPostFeedParams {
	mmGetPostFeed.mutex.RLock()

	argCopy := make([]*CacheRepoMockGetPostFeedParams, len(mmGetPostFeed.callArgs))
	copy(argCopy, mmGetPostFeed.callArgs)

	mmGetPostFeed.mutex.RUnlock()

	return argCopy
}

// MinimockGetPostFeedDone returns true if the count of the GetPostFeed invocations corresponds
// the number of defined expectations
func (m *CacheRepoMock) MinimockGetPostFeedDone() bool {
	for _, e := range m.GetPostFeedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetPostFeedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetPostFeedCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetPostFeed != nil && mm_atomic.LoadUint64(&m.afterGetPostFeedCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetPostFeedInspect logs each unmet expectation
func (m *CacheRepoMock) MinimockGetPostFeedInspect() {
	for _, e := range m.GetPostFeedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheRepoMock.GetPostFeed with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetPostFeedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetPostFeedCounter) < 1 {
		if m.GetPostFeedMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CacheRepoMock.GetPostFeed")
		} else {
			m.t.Errorf("Expected call to CacheRepoMock.GetPostFeed with params: %#v", *m.GetPostFeedMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetPostFeed != nil && mm_atomic.LoadUint64(&m.afterGetPostFeedCounter) < 1 {
		m.t.Error("Expected call to CacheRepoMock.GetPostFeed")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CacheRepoMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetCelebrityFriendsInspect()

		m.MinimockGetPostFeedInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CacheRepoMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *CacheRepoMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetCelebrityFriendsDone() &&
		m.MinimockGetPostFeedDone()
}
