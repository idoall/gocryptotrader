// Code generated by SQLBoiler 3.5.0-gct (https://github.com/thrasher-corp/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package postgres

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("AuditEvents", testAuditEvents)
	t.Run("Scripts", testScripts)
	t.Run("WithdrawalHistories", testWithdrawalHistories)
}

func TestDelete(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsDelete)
	t.Run("Scripts", testScriptsDelete)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsQueryDeleteAll)
	t.Run("Scripts", testScriptsQueryDeleteAll)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsSliceDeleteAll)
	t.Run("Scripts", testScriptsSliceDeleteAll)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsExists)
	t.Run("Scripts", testScriptsExists)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesExists)
}

func TestFind(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsFind)
	t.Run("Scripts", testScriptsFind)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesFind)
}

func TestBind(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsBind)
	t.Run("Scripts", testScriptsBind)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesBind)
}

func TestOne(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsOne)
	t.Run("Scripts", testScriptsOne)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesOne)
}

func TestAll(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsAll)
	t.Run("Scripts", testScriptsAll)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesAll)
}

func TestCount(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsCount)
	t.Run("Scripts", testScriptsCount)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesCount)
}

func TestHooks(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsHooks)
	t.Run("Scripts", testScriptsHooks)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesHooks)
}

func TestInsert(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsInsert)
	t.Run("AuditEvents", testAuditEventsInsertWhitelist)
	t.Run("Scripts", testScriptsInsert)
	t.Run("Scripts", testScriptsInsertWhitelist)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesInsert)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsReload)
	t.Run("Scripts", testScriptsReload)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsReloadAll)
	t.Run("Scripts", testScriptsReloadAll)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsSelect)
	t.Run("Scripts", testScriptsSelect)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsUpdate)
	t.Run("Scripts", testScriptsUpdate)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("AuditEvents", testAuditEventsSliceUpdateAll)
	t.Run("Scripts", testScriptsSliceUpdateAll)
	t.Run("WithdrawalHistories", testWithdrawalHistoriesSliceUpdateAll)
}