/*
 *        Copyright 2023, 2025 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "openglloader.h"

#if defined(G2D_OPENGLLOADER_WIN32)

#define OPENGLLOADER_ERR_ALLOC 101
#define OPENGLLOADER_ERR_1 1000001
#define OPENGLLOADER_ERR_2 1000002
#define OPENGLLOADER_ERR_3 1000003
#define OPENGLLOADER_ERR_4 1000004
#define OPENGLLOADER_ERR_5 1000005
#define OPENGLLOADER_ERR_6 1000006
#define OPENGLLOADER_ERR_7 1000007
#define OPENGLLOADER_ERR_8 1000008
#define OPENGLLOADER_ERR_9 1000009

#define OPENGLLOADER_CDATA_ID "vbsw.g2d.openglloader"
#define CLASS_NAME TEXT("openglloader_dummy")

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>

/* external */
typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

/* protected */
typedef struct { void *load_func; void *glGetIntegerv; } openglloader_t;
typedef void* (openglloader_load_func_t) (const char *name, long long *err);
typedef void (openglloader_glGetIntegerv_t) (unsigned int pname, int *params);

/* private */
typedef struct { void *load_func; void *glGetIntegerv; HINSTANCE instance; HWND hndl; HDC dc; HGLRC rc; } openglloader_data_t;
typedef union { openglloader_t loader; openglloader_data_t data; } openglloader_union_t;

static void* openglloader_load(const char *const name, long long *const err) {
	/* wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions). */
	PROC const proc = wglGetProcAddress(name);
	const DWORD last_error = GetLastError();
	if (last_error == 0) {
		return (void*)proc;
	} else {
		err[0] = (long long)last_error;
		return NULL;
	}
}

static void openglloader_glGetIntegerv(const unsigned int pname, int *const params) {
	glGetIntegerv((GLenum) pname, (GLint *)params);
}

void g2d_openglloader_init(const int pass, cdata_t *const cdata) {
	cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
	cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
	/* init and set current OpenGL context */
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
										openglloader_union_t *const loader = (openglloader_union_t*)malloc(sizeof(openglloader_union_t));
										if (loader) {
											loader[0].data.load_func = (void*)openglloader_load;
											loader[0].data.glGetIntegerv = (void*)openglloader_glGetIntegerv;
											loader[0].data.instance = instance;
											loader[0].data.hndl = dummy_hndl;
											loader[0].data.dc = dummy_dc;
											loader[0].data.rc = dummy_rc;
											set(cdata, (void*)loader, OPENGLLOADER_CDATA_ID);
										} else {
											cdata[0].err1 = OPENGLLOADER_ERR_ALLOC;
											wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
											DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
										}
									} else {
										cdata[0].err1 = OPENGLLOADER_ERR_8; cdata[0].err2 = (long long)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
									}
								} else {
									cdata[0].err1 = OPENGLLOADER_ERR_7; cdata[0].err2 = (long long)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
								}
							} else {
								cdata[0].err1 = OPENGLLOADER_ERR_6; cdata[0].err2 = (long long)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
							}
						} else {
							cdata[0].err1 = OPENGLLOADER_ERR_5; cdata[0].err2 = (long long)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
						}
					} else {
						cdata[0].err1 = OPENGLLOADER_ERR_4;
						DestroyWindow(dummy_hndl); UnregisterClass(CLASS_NAME, instance);
					}
				} else {
					cdata[0].err1 = OPENGLLOADER_ERR_3; cdata[0].err2 = (long long)GetLastError();
					UnregisterClass(CLASS_NAME, instance);
				}
			} else {
				cdata[0].err1 = OPENGLLOADER_ERR_2; cdata[0].err2 = (long long)GetLastError();
			}
		} else {
			cdata[0].err1 = OPENGLLOADER_ERR_1; cdata[0].err2 = (long long)GetLastError();
		}
	/* clean up (after success) */
	} else if (pass == 1) {
		openglloader_data_t *const data = (openglloader_data_t*)get(cdata, OPENGLLOADER_CDATA_ID);
		if (data) {
			if (wglGetCurrentContext() == data[0].rc)
				wglMakeCurrent(NULL, NULL);
			wglDeleteContext(data[0].rc); ReleaseDC(data[0].hndl, data[0].dc);
			DestroyWindow(data[0].hndl); UnregisterClass(CLASS_NAME, data[0].instance);
			free(data);
		} else {
			cdata[0].err1 = OPENGLLOADER_ERR_9;
		}
	/* clean up (after error) */
	} else if (pass < 0) {
		openglloader_data_t *const data = (openglloader_data_t*)get(cdata, OPENGLLOADER_CDATA_ID);
		if (data) {
			if (wglGetCurrentContext() == data[0].rc)
				wglMakeCurrent(NULL, NULL);
			wglDeleteContext(data[0].rc); ReleaseDC(data[0].hndl, data[0].dc);
			DestroyWindow(data[0].hndl); UnregisterClass(CLASS_NAME, data[0].instance);
			free(data);
		}
	}
}

#elif defined(G2D_OPENGLLOADER_LINUX)

#endif
