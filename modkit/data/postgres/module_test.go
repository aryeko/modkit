package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync"
	"testing"

	"github.com/go-modkit/modkit/modkit/data/sqlmodule"
	"github.com/go-modkit/modkit/modkit/testkit"
)

var testDrv = &countingDriver{}

func init() {
	sql.Register(driverName, testDrv)
}

type countingDriver struct {
	mu          sync.Mutex
	openCount   int
	pingCount   int
	closeCount  int
	pingErr     error
	sawDeadline bool
}

func (d *countingDriver) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	c := countingDriver{}
	d.openCount = c.openCount
	d.pingCount = c.pingCount
	d.closeCount = c.closeCount
	d.pingErr = nil
	d.sawDeadline = false
}

func (d *countingDriver) SetPingErr(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pingErr = err
}

func (d *countingDriver) Snapshot() (open, ping, closed int, sawDeadline bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.openCount, d.pingCount, d.closeCount, d.sawDeadline
}

func (d *countingDriver) Open(_ string) (driver.Conn, error) {
	d.mu.Lock()
	d.openCount++
	d.mu.Unlock()
	return &countingConn{d: d}, nil
}

type countingConn struct {
	d *countingDriver
}

func (c *countingConn) Prepare(_ string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *countingConn) Close() error {
	c.d.mu.Lock()
	c.d.closeCount++
	c.d.mu.Unlock()
	return nil
}

func (c *countingConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func (c *countingConn) Ping(ctx context.Context) error {
	c.d.mu.Lock()
	c.d.pingCount++
	if _, ok := ctx.Deadline(); ok {
		c.d.sawDeadline = true
	}
	err := c.d.pingErr
	c.d.mu.Unlock()
	return err
}

func TestModuleExportsDialectAndDBTokens(t *testing.T) {
	testDrv.Reset()
	t.Setenv("POSTGRES_DSN", "test")
	t.Setenv("POSTGRES_CONNECT_TIMEOUT", "0")

	h := testkit.New(t, NewModule(Options{}))
	_ = testkit.Get[*sql.DB](t, h, sqlmodule.TokenDB)
	dialect := testkit.Get[sqlmodule.Dialect](t, h, sqlmodule.TokenDialect)
	if dialect != sqlmodule.DialectPostgres {
		t.Fatalf("unexpected dialect: %q", dialect)
	}
}

func TestConnectTimeoutZeroSkipsPing(t *testing.T) {
	testDrv.Reset()
	t.Setenv("POSTGRES_DSN", "test")
	t.Setenv("POSTGRES_CONNECT_TIMEOUT", "0")

	h := testkit.New(t, NewModule(Options{}))
	_ = testkit.Get[*sql.DB](t, h, sqlmodule.TokenDB)

	open, ping, _, _ := testDrv.Snapshot()
	if open != 0 {
		t.Fatalf("expected open=0, got %d", open)
	}
	if ping != 0 {
		t.Fatalf("expected ping=0, got %d", ping)
	}
}

func TestConnectTimeoutNonZeroPingsWithTimeout(t *testing.T) {
	testDrv.Reset()
	t.Setenv("POSTGRES_DSN", "test")
	t.Setenv("POSTGRES_CONNECT_TIMEOUT", "25ms")

	h := testkit.New(t, NewModule(Options{}))
	_ = testkit.Get[*sql.DB](t, h, sqlmodule.TokenDB)

	open, ping, _, sawDeadline := testDrv.Snapshot()
	if open == 0 {
		t.Fatalf("expected open>0, got %d", open)
	}
	if ping != 1 {
		t.Fatalf("expected ping=1, got %d", ping)
	}
	if !sawDeadline {
		t.Fatalf("expected ping to observe a context deadline")
	}
}

func TestPingFailureReturnsTypedBuildErrorAndClosesDB(t *testing.T) {
	testDrv.Reset()
	pingErr := errors.New("ping failed")
	testDrv.SetPingErr(pingErr)
	t.Setenv("POSTGRES_DSN", "test")
	t.Setenv("POSTGRES_CONNECT_TIMEOUT", "25ms")

	h := testkit.New(t, NewModule(Options{}))
	_, err := testkit.GetE[*sql.DB](h, sqlmodule.TokenDB)
	if err == nil {
		t.Fatalf("expected error")
	}

	var be *BuildError
	if !errors.As(err, &be) {
		t.Fatalf("expected BuildError, got %T", err)
	}
	if be.Stage != StagePing {
		t.Fatalf("expected stage=%s, got %s", StagePing, be.Stage)
	}
	if be.Token != sqlmodule.TokenDB {
		t.Fatalf("expected token=%q, got %q", sqlmodule.TokenDB, be.Token)
	}
	if !errors.Is(err, pingErr) {
		t.Fatalf("expected error to wrap ping error")
	}

	_, _, closed, _ := testDrv.Snapshot()
	if closed == 0 {
		t.Fatalf("expected ping failure path to close the DB")
	}
}
