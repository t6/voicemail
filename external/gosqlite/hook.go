// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite

/*
#include <sqlite3.h>

void* goSqlite3CommitHook(sqlite3 *db, void *udp);
void* goSqlite3RollbackHook(sqlite3 *db, void *udp);
void* goSqlite3UpdateHook(sqlite3 *db, void *udp);
//void* goSqlite3WalHook(sqlite3 *db, void *udp);
*/
import "C"

import (
	"unsafe"
)

// CommitHook is the callback function signature.
// If the callback on a commit hook function returns true, then the commit is converted into a rollback.
type CommitHook func(udp interface{}) (rollback bool)

type sqliteCommitHook struct {
	f   CommitHook
	udp interface{}
}

//export goXCommitHook
func goXCommitHook(udp unsafe.Pointer) C.int {
	arg := (*sqliteCommitHook)(udp)
	return btocint(arg.f(arg.udp))
}

// CommitHook registers a callback function to be invoked whenever a transaction is committed.
// (See http://sqlite.org/c3ref/commit_hook.html)
func (c *Conn) CommitHook(f CommitHook, udp interface{}) {
	if f == nil {
		c.commitHook = nil
		C.sqlite3_commit_hook(c.db, nil, nil)
		return
	}
	// To make sure it is not gced, keep a reference in the connection.
	c.commitHook = &sqliteCommitHook{f, udp}
	C.goSqlite3CommitHook(c.db, unsafe.Pointer(c.commitHook))
}

// RollbackHook is the callback function signature.
type RollbackHook func(udp interface{})

type sqliteRollbackHook struct {
	f   RollbackHook
	udp interface{}
}

//export goXRollbackHook
func goXRollbackHook(udp unsafe.Pointer) {
	arg := (*sqliteRollbackHook)(udp)
	arg.f(arg.udp)
}

// RollbackHook registers a callback to be invoked each time a transaction is rolled back.
// (See http://sqlite.org/c3ref/commit_hook.html)
func (c *Conn) RollbackHook(f RollbackHook, udp interface{}) {
	if f == nil {
		c.rollbackHook = nil
		C.sqlite3_rollback_hook(c.db, nil, nil)
		return
	}
	// To make sure it is not gced, keep a reference in the connection.
	c.rollbackHook = &sqliteRollbackHook{f, udp}
	C.goSqlite3RollbackHook(c.db, unsafe.Pointer(c.rollbackHook))
}

// UpdateHook is the callback function signature.
type UpdateHook func(udp interface{}, a Action, dbName, tableName string, rowID int64)

type sqliteUpdateHook struct {
	f   UpdateHook
	udp interface{}
}

//export goXUpdateHook
func goXUpdateHook(udp unsafe.Pointer, action int, dbName, tableName *C.char, rowID C.sqlite3_int64) {
	arg := (*sqliteUpdateHook)(udp)
	arg.f(arg.udp, Action(action), C.GoString(dbName), C.GoString(tableName), int64(rowID))
}

// UpdateHook registers a callback to be invoked each time a row is updated,
// inserted or deleted using this database connection.
// (See http://sqlite.org/c3ref/update_hook.html)
func (c *Conn) UpdateHook(f UpdateHook, udp interface{}) {
	if f == nil {
		c.updateHook = nil
		C.sqlite3_update_hook(c.db, nil, nil)
		return
	}
	// To make sure it is not gced, keep a reference in the connection.
	c.updateHook = &sqliteUpdateHook{f, udp}
	C.goSqlite3UpdateHook(c.db, unsafe.Pointer(c.updateHook))
}

/*
type WalHook func(udp interface{}, c *Conn, dbName string, nEntry int) int

type sqliteWalHook struct {
	f   WalHook
	udp interface{}
}

//export goXWalHook
func goXWalHook(udp, db unsafe.Pointer, dbName *C.char, nEntry C.int) C.int {
	return 0
}

// Register a callback to be invoked each time a transaction is written
// into the write-ahead-log by this database connection.
// (See http://sqlite.org/c3ref/wal_hook.html)
func (c *Conn) WalHook(f WalHook, udp interface{}) {
	if f == nil {
		c.walHook = nil
		C.sqlite3_wal_hook(c.db, nil, nil)
		return
	}
	// To make sure it is not gced, keep a reference in the connection.
	c.walHook = &sqliteWalHook{f, udp}
	C.goSqlite3WalHook(c.db, unsafe.Pointer(c.walHook))
}
*/
