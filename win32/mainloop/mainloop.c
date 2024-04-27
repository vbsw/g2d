/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "mainloop.h"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* Exported Go functions:                                   */
/* g2dMainLoopInit                                          */
/* g2dMainLoopProcessCustomEvents                           */

static const WPARAM g2d_CUSTOM_EVENT = (WPARAM)"g2dc";
static const WPARAM g2d_QUIT_EVENT   = (WPARAM)"g2dq";

static DWORD thread_id;
static BOOL stop;

void g2d_mainloop_process_messages() {
	MSG msg; BOOL ret_code; int skip = 0; stop = FALSE;
	thread_id = GetCurrentThreadId();
	g2dMainLoopInit();
	while (!stop && (ret_code = GetMessage(&msg, NULL, 0, 0)) > 0) {
		if (msg.message == WM_APP) {
			if (msg.wParam == g2d_CUSTOM_EVENT)
				if (skip == 0)
					g2dMainLoopProcessCustomEvents(&skip);
				else
					skip--;
			else if (msg.wParam == g2d_QUIT_EVENT)
				break;
		} else {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
}

void g2d_mainloop_post_custom(long long *const err) {
	if (!PostThreadMessage(thread_id, WM_APP, g2d_CUSTOM_EVENT, 0))
		err[0] = (long long)GetLastError();
}

void g2d_mainloop_post_quit(long long *const err) {
	stop = TRUE;
	if (!PostThreadMessage(thread_id, WM_APP, g2d_QUIT_EVENT, 0))
		err[0] = (long long)GetLastError();
}

void g2d_mainloop_clean_up() {
	MSG msg;
	while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE));
}
