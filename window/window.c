/*
 *       Copyright 2024, 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "window.h"

#if defined(G2D_WINDOW_WIN32)

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* external */
typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

#define WINDOW_ERR_ALLOC_1 201
#define WINDOW_ERR_1 1000201
#define WINDOW_ERR_2 1000202
#define WINDOW_ERR_3 1000203
#define WINDOW_ERR_4 1000204
#define WINDOW_ERR_5 1000205

#define WINDOW_CDATA_ID "vbsw.g2d.window"

/* wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions). */
#define LOAD_FUNC(t, n, e) if (err1[0] == 0) { PROC const proc = wglGetProcAddress(#n); const DWORD last_err = GetLastError(); if (last_err == 0) n = (t) proc; else { err1[0] = e; err2[0] = (long long)last_err; }}

/* from wglext.h */
typedef BOOL(WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC hdc, const int *piAttribIList, const FLOAT *pfAttribFList, UINT nMaxFormats, int *piFormats, UINT *nNumFormats);
typedef HGLRC(WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC hDC, HGLRC hShareContext, const int *attribList);
typedef BOOL(WINAPI * PFNWGLSWAPINTERVALEXTPROC) (int interval);
typedef int (WINAPI * PFNWGLGETSWAPINTERVALEXTPROC) (void);

typedef struct {} window_t;

static PFNWGLCHOOSEPIXELFORMATARBPROC    wglChoosePixelFormatARB    = NULL;
static PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
static PFNWGLSWAPINTERVALEXTPROC         wglSwapIntervalEXT         = NULL;
static PFNWGLGETSWAPINTERVALEXTPROC      wglGetSwapIntervalEXT      = NULL;

static const WPARAM wnd_CUSTOM_EVENT = (WPARAM)"vbsw_g2dc";
static const WPARAM wnd_QUIT_EVENT   = (WPARAM)"vbsw_g2dq";
static LPCTSTR const class_name      = TEXT("vbsw_g2d_window");

static BOOL initialized   = FALSE;
static HINSTANCE instance = NULL;
static int windows_count  = 0;
static DWORD thread_id    = 0;
static BOOL stop          = FALSE;

static window_t window;

void g2d_window_mainloop_process() {
	MSG msg; BOOL ret_code; stop = FALSE;
	thread_id = GetCurrentThreadId();
	g2dWindowInit();
	while (!stop && (ret_code = GetMessage(&msg, NULL, 0, 0)) > 0) {
		if (msg.message == WM_APP) {
			if (msg.wParam == wnd_CUSTOM_EVENT)
				g2dWindowProcessMessages();
			else if (msg.wParam == wnd_QUIT_EVENT)
				break;
		} else {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
}

void g2d_window_post_custom_msg(long long *const err1, long long *const err2) {
	if (!PostThreadMessage(thread_id, WM_APP, wnd_CUSTOM_EVENT, 0)) {
		err1[0] = 3999;
		err2[0] = (long long)GetLastError();
	}
}

void g2d_window_post_quit_msg(long long *const err1, long long *const err2) {
	if (!PostThreadMessage(thread_id, WM_APP, wnd_QUIT_EVENT, 0)) {
		err1[0] = 3999;
		err2[0] = (long long)GetLastError();
	}
	stop = TRUE;
}

void g2d_window_mainloop_clean_up() {
	MSG msg;
	while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE));
}

void g2d_window_init(const int pass, cdata_t *const cdata) {
	if (pass == 0) {
		if (!initialized) {
			long long *const err1 = &cdata[0].err1;
			long long *const err2 = &cdata[0].err2;
			instance = GetModuleHandle(NULL);
			if (instance) {
				cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
				cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
				/* wgl functions */
				LOAD_FUNC(PFNWGLCHOOSEPIXELFORMATARBPROC,    wglChoosePixelFormatARB,    WINDOW_ERR_2)
				LOAD_FUNC(PFNWGLCREATECONTEXTATTRIBSARBPROC, wglCreateContextAttribsARB, WINDOW_ERR_3)
				LOAD_FUNC(PFNWGLSWAPINTERVALEXTPROC,         wglSwapIntervalEXT,         WINDOW_ERR_4)
				LOAD_FUNC(PFNWGLGETSWAPINTERVALEXTPROC,      wglGetSwapIntervalEXT,      WINDOW_ERR_5)
				if (err1[0] == 0)
					set(cdata, (void*)&window, WINDOW_CDATA_ID);
				if (err1[0] == 0)
					initialized = TRUE;
			} else {
				err1[0] = WINDOW_ERR_1; err2[0] = (long long)GetLastError();
			}
		}
	} else if (pass < 0 && initialized) {
		initialized = FALSE;
	}
}

#elif defined(G2D_WINDOW_LINUX)

#endif
