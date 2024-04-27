/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#if defined(G2D_IMPL_WIN32)

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>
#include "win32.h"

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

typedef struct {
	int bla;
} engine_t;

static const WPARAM const g2d_EVENT = (WPARAM)"g2d";
static LPCTSTR const class_name = TEXT("g2d");
static LPCTSTR const class_name_dummy = TEXT("g2d_dummy");

static HINSTANCE instance = NULL;

//#include "win32_window.h"

static void free_engine(engine_t *const engine) {
	free(engine);
}

void g2d_win32_init(void **const e, void **const d, void **const fs, const int d_len, int *const max_t_size, int *const err_num, g2d_ul_t *const err_win32) {
	if (instance == NULL)
		instance = GetModuleHandle(NULL);
	if (instance) {
		engine_t *const engine = (engine_t*)malloc(sizeof(engine_t));
		if (engine) {
			ZeroMemory(engine, sizeof(engine_t));
			/* dummy class */
			WNDCLASSEX cls;
			ZeroMemory(&cls, sizeof(WNDCLASSEX));
			cls.cbSize = sizeof(WNDCLASSEX);
			cls.style = CS_OWNDC | CS_HREDRAW | CS_VREDRAW;
			cls.lpfnWndProc = DefWindowProc;
			cls.hInstance = instance;
			cls.lpszClassName = class_name_dummy;
			if (RegisterClassEx(&cls) != INVALID_ATOM) {
				/* dummy window */
				HWND const dummy_hndl = CreateWindow(class_name_dummy, TEXT("Dummy"), WS_OVERLAPPEDWINDOW, 0, 0, 1, 1, NULL, NULL, instance, NULL);
				if (dummy_hndl) {
					/* dummy context */
					HDC const dummy_dc = GetDC(dummy_hndl);
					if (dummy_dc) {
						int pixelFormat;
						PIXELFORMATDESCRIPTOR pixelFormatDesc;
						ZeroMemory(&pixelFormatDesc, sizeof(PIXELFORMATDESCRIPTOR));
						pixelFormatDesc.nSize = sizeof(PIXELFORMATDESCRIPTOR);
						pixelFormatDesc.nVersion = 1;
						pixelFormatDesc.dwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL;
						pixelFormatDesc.iPixelType = PFD_TYPE_RGBA;
						pixelFormatDesc.cColorBits = 32;
						pixelFormatDesc.cAlphaBits = 8;
						pixelFormatDesc.cDepthBits = 24;
						pixelFormat = ChoosePixelFormat(dummy_dc, &pixelFormatDesc);
						if (pixelFormat) {
							if (SetPixelFormat(dummy_dc, pixelFormat, &pixelFormatDesc)) {
								HGLRC const dummy_rc = wglCreateContext(dummy_dc);
								if (dummy_rc) {
									if (wglMakeCurrent(dummy_dc, dummy_rc)) {
										glGetIntegerv(GL_MAX_TEXTURE_SIZE, max_t_size);
										/* destroy dummy */
										if (!wglMakeCurrent(NULL, NULL) && err_num[0] == 0) {
											err_num[0] = 108; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										if (!wglDeleteContext(dummy_rc) && err_num[0] == 0) {
											err_num[0] = 109; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										ReleaseDC(dummy_hndl, dummy_dc);
										if (!DestroyWindow(dummy_hndl) && err_num[0] == 0) {
											err_num[0] = 110; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										if (!UnregisterClass(class_name_dummy, instance) && err_num[0] == 0) {
											err_num[0] = 111; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										if (err_num[0] == 0) {
											e[0] = (void*)engine;
										} else {
											// TODO destroy modules
										}
									} else {
										err_num[0] = 107; err_win32[0] = (g2d_ul_t)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance); free_engine(engine);
									}
								} else {
									err_num[0] = 106; err_win32[0] = (g2d_ul_t)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance); free_engine(engine);
								}
							} else {
								err_num[0] = 105; err_win32[0] = (g2d_ul_t)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance); free_engine(engine);
							}
						} else {
							err_num[0] = 104; err_win32[0] = (g2d_ul_t)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance); free_engine(engine);
						}
					} else {
						err_num[0] = 103;
						DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance); free_engine(engine);
					}
				} else {
					err_num[0] = 102; err_win32[0] = (g2d_ul_t)GetLastError();
					UnregisterClass(class_name_dummy, instance); free_engine(engine);
				}
			} else {
				err_num[0] = 101; err_win32[0] = (g2d_ul_t)GetLastError(); free_engine(engine);
			}
		} else {
			err_num[0] = 1;
		}
	} else {
		err_num[0] = 100; err_win32[0] = (g2d_ul_t)GetLastError();
	}
}

/* #if defined(G2D_IMPL_WIN32) */
#endif
