/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "ogfl.h"

#if defined(G2D_OGFL_WIN32)

#define OGFL_CDATA_ID "vbsw.g2d.ogfl"
#define CLASS_NAME TEXT("ogfl_dummy")

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

typedef void* (ogfl_load_func_t) (void *obj, const char *name, long long *err);
typedef struct { ogfl_load_func_t *load_func; void *obj; } oglf_t;
typedef struct { HINSTANCE instance; HWND hndl; HDC dc; HGLRC rc; } oglf_obj_t;

static void* ogfl_load(void *const obj, const char *const name, long long *const err) {
	PROC const proc = wglGetProcAddress(name);
	const DWORD last_error = GetLastError();
	if (last_error == 0) {
		return (void*)proc;
	} else {
		err[0] = (long long)last_error;
		return NULL;
	}
}

void vbsw_ogfl_init(const int pass, cdata_t *const cdata) {
	cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
	cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
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
										oglf_t *const oglf = (oglf_t*)malloc(sizeof(oglf_t) + sizeof(oglf_obj_t));
										if (oglf) {
											oglf_obj_t *const obj = (oglf_obj_t*)((size_t)oglf + sizeof(oglf_t));
											oglf[0].load_func = ogfl_load;
											oglf[0].obj = (void*)obj;
											obj[0].instance = instance;
											obj[0].hndl = dummy_hndl;
											obj[0].dc = dummy_dc;
											obj[0].rc = dummy_rc;
											set(cdata, (void*)oglf, OGFL_CDATA_ID);
										} else {
											cdata[0].err1 = 10;
											wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
											DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
										}
									} else {
										cdata[0].err1 = 1007; cdata[0].err2 = (long long)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
									}
								} else {
									cdata[0].err1 = 1006; cdata[0].err2 = (long long)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
								}
							} else {
								cdata[0].err1 = 1005; cdata[0].err2 = (long long)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
							}
						} else {
							cdata[0].err1 = 1004; cdata[0].err2 = (long long)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
						}
					} else {
						cdata[0].err1 = 1003;
						DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
					}
				} else {
					cdata[0].err1 = 1002; cdata[0].err2 = (long long)GetLastError();
					UnregisterClass(CLASS_NAME, instance);
				}
			} else {
				cdata[0].err1 = 1001; cdata[0].err2 = (long long)GetLastError();
			}
		} else {
			cdata[0].err1 = 1000; cdata[0].err2 = (long long)GetLastError();
		}
	} else if (pass == 1) {
		oglf_t *const oglf = (oglf_t*)get(cdata, OGFL_CDATA_ID);
		if (oglf) {
			oglf_obj_t *const obj = (oglf_obj_t*)oglf[0].obj;
			if (wglGetCurrentContext() == obj[0].rc)
				wglMakeCurrent(NULL, NULL);
			wglDeleteContext(obj[0].rc); ReleaseDC(obj[0].hndl, obj[0].dc);
			DestroyWindow(obj[0].hndl); UnregisterClass(CLASS_NAME, obj[0].instance);
			free(oglf);
		} else {
			cdata[0].err1 = 1008;
		}
	} else if (pass < 0) {
		oglf_t *const oglf = (oglf_t*)get(cdata, OGFL_CDATA_ID);
		if (oglf) {
			oglf_obj_t *const obj = (oglf_obj_t*)oglf[0].obj;
			if (wglGetCurrentContext() == obj[0].rc)
				wglMakeCurrent(NULL, NULL);
			wglDeleteContext(obj[0].rc); ReleaseDC(obj[0].hndl, obj[0].dc);
			DestroyWindow(obj[0].hndl); UnregisterClass(CLASS_NAME, obj[0].instance);
			free(oglf);
		}
	}
}

#elif defined(G2D_OGFL_LINUX)

#endif
