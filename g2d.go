/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "C"
import (
	"time"
	"sync"
)

var (
	errs []error
	mutex      sync.Mutex
	errGen ErrorGenerator
	errLog ErrorLogger
	errHandler tErrorHandler
	initialized bool
	startTime   time.Time
)

func Errors() []error {
	mutex.Lock()
	defer mutex.Unlock()
	return errs
}

func appendError(err error) {
	mutex.Lock()
	errs = append(errs, err)
	errLog.LogError(err)
	mutex.Unlock()
}

func deltaNanos() int64 {
	timeNow := time.Now()
	d := timeNow.Sub(startTime)
	return d.Nanoseconds()
}
