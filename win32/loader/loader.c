/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "loader.h"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#define OGFL_ID "vbsw.g2d.loader"
#define CLASS_NAME TEXT("loader_dummy")

/* from github.com/vbsw/golib/cdata/cdata.c */
typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

/* for external usage */
typedef void* (loader_load_func_t) (void *obj, const char *name, long long *err2);
typedef struct { loader_load_func_t *load_func; void *obj; } loader_t;

/* for internal usage */
typedef struct { HINSTANCE instance; HWND hndl; HDC dc; HGLRC rc; } obj_t;

static void* loader_load(void *const obj, const char *const name, long long *const err2) {
	// wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions).
	PROC const proc = wglGetProcAddress(name);
	const DWORD last_error = GetLastError();
	if (last_error == 0) {
		return (void*)proc;
	} else {
		err2[0] = (long long)last_error;
		return NULL;
	}
}

void vbsw_loader_init(const int pass, cdata_t *const cdata) {
	if (pass == 0) {
		HINSTANCE const instance = GetModuleHandle(NULL);
		if (instance) {
			/* dummy class */
			WNDCLASSEX cls;
			ZeroMemory(&cls, sizeof(WNDCLASSEX));
			cls.cbSize = sizeof(WNDCLASSEX);
			cls.style = CS_OWNDC;
			cls.lpfnWndProc = DefWindowProc;
			cls.hInstance = instance;
			cls.lpszClassName = CLASS_NAME;
			if (RegisterClassEx(&cls) != INVALID_ATOM) {
				/* dummy window */
				HWND const dummy_hndl = CreateWindow(CLASS_NAME, TEXT("Dummy"), WS_OVERLAPPEDWINDOW, 0, 0, 1, 1, NULL, NULL, instance, NULL);
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
										loader_t *const loader = (loader_t*)malloc(sizeof(loader_t) + sizeof(obj_t));
										if (loader) {
											cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
											obj_t *const obj = (obj_t*)((size_t)loader + sizeof(loader_t));
											loader[0].load_func = loader_load;
											loader[0].obj = (void*)obj;
											obj[0].instance = instance;
											obj[0].hndl = dummy_hndl;
											obj[0].dc = dummy_dc;
											obj[0].rc = dummy_rc;
											set(cdata, (void*)loader, OGFL_ID);
										} else {
											cdata[0].err1 = 10;
											wglMakeCurrent(NULL, NULL);
											wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
											DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
										}
									} else {
										cdata[0].err1 = 1000007; cdata[0].err2 = (long long)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
									}
								} else {
									cdata[0].err1 = 1000006; cdata[0].err2 = (long long)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
								}
							} else {
								cdata[0].err1 = 1000005; cdata[0].err2 = (long long)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
							}
						} else {
							cdata[0].err1 = 1000004; cdata[0].err2 = (long long)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
						}
					} else {
						cdata[0].err1 = 1000003;
						DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
					}
				} else {
					cdata[0].err1 = 1000002; cdata[0].err2 = (long long)GetLastError();
					UnregisterClass(CLASS_NAME, instance);
				}
			} else {
				cdata[0].err1 = 1000001; cdata[0].err2 = (long long)GetLastError();
			}
		} else {
			cdata[0].err1 = 1000000; cdata[0].err2 = (long long)GetLastError();
		}
	} else if (pass < 0 || pass == 1) {
		cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
		loader_t *const loader = (loader_t*)get(cdata, OGFL_ID);
		if (loader) {
			obj_t *const obj = (obj_t*)loader[0].obj;
			if (wglGetCurrentContext() == obj[0].rc)
				wglMakeCurrent(NULL, NULL);
			wglDeleteContext(obj[0].rc); ReleaseDC(obj[0].hndl, obj[0].dc);
			DestroyWindow(obj[0].hndl); UnregisterClass(CLASS_NAME, obj[0].instance);
			free(loader);
		} else if (pass == 1) {
			cdata[0].err1 = 1000008;
		}
	}
}
