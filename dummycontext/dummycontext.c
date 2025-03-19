/*
 *       Copyright 2023, 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "dummycontext.h"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#define DUMMYCONTEXT_ERR_1 1000001
#define DUMMYCONTEXT_ERR_2 1000002
#define DUMMYCONTEXT_ERR_3 1000003
#define DUMMYCONTEXT_ERR_4 1000004
#define DUMMYCONTEXT_ERR_5 1000005
#define DUMMYCONTEXT_ERR_6 1000006
#define DUMMYCONTEXT_ERR_7 1000007
#define DUMMYCONTEXT_ERR_8 1000008

static LPCTSTR const class_name = TEXT("g2d_window_dummy");
static BOOL initialized = FALSE;
static HINSTANCE instance;
static HWND hndl;
static HDC dc;
static HGLRC rc;

void g2d_dummycontext_init(const int pass, cdata_t *const cdata) {
	/* init and set current OpenGL context */
	if (pass == 0) {
		instance = GetModuleHandle(NULL);
		if (instance) {
			/* dummy class */
			WNDCLASSEX cls;
			ZeroMemory(&cls, sizeof(WNDCLASSEX));
			cls.cbSize = sizeof(WNDCLASSEX);
			cls.style = CS_OWNDC;
			cls.lpfnWndProc = DefWindowProc;
			cls.hInstance = instance;
			cls.lpszClassName = class_name;
			if (RegisterClassEx(&cls) != INVALID_ATOM) {
				/* dummy window */
				hndl = CreateWindow(class_name, TEXT("Dummy"), WS_OVERLAPPEDWINDOW, 0, 0, 1, 1, NULL, NULL, instance, NULL);
				if (hndl) {
					/* dummy context */
					dc = GetDC(hndl);
					if (dc) {
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
						pixelFormat = ChoosePixelFormat(dc, &pixelFormatDesc);
						if (pixelFormat) {
							if (SetPixelFormat(dc, pixelFormat, &pixelFormatDesc)) {
								rc = wglCreateContext(dc);
								if (rc) {
									if (wglMakeCurrent(dc, rc)) {
										initialized = TRUE;
									} else {
										cdata[0].err1 = DUMMYCONTEXT_ERR_8; cdata[0].err2 = (long long)GetLastError();
										wglDeleteContext(rc); ReleaseDC(hndl, dc);
										DestroyWindow(hndl); UnregisterClass(class_name, instance);
									}
								} else {
									cdata[0].err1 = DUMMYCONTEXT_ERR_7; cdata[0].err2 = (long long)GetLastError();
									ReleaseDC(hndl, dc); DestroyWindow(hndl); UnregisterClass(class_name, instance);
								}
							} else {
								cdata[0].err1 = DUMMYCONTEXT_ERR_6; cdata[0].err2 = (long long)GetLastError();
								ReleaseDC(hndl, dc); DestroyWindow(hndl); UnregisterClass(class_name, instance);
							}
						} else {
							cdata[0].err1 = DUMMYCONTEXT_ERR_5; cdata[0].err2 = (long long)GetLastError();
							ReleaseDC(hndl, dc); DestroyWindow(hndl); UnregisterClass(class_name, instance);
						}
					} else {
						cdata[0].err1 = DUMMYCONTEXT_ERR_4;
						DestroyWindow(hndl); UnregisterClass(class_name, instance);
					}
				} else {
					cdata[0].err1 = DUMMYCONTEXT_ERR_3; cdata[0].err2 = (long long)GetLastError();
					UnregisterClass(class_name, instance);
				}
			} else {
				cdata[0].err1 = DUMMYCONTEXT_ERR_2; cdata[0].err2 = (long long)GetLastError();
			}
		} else {
			cdata[0].err1 = DUMMYCONTEXT_ERR_1; cdata[0].err2 = (long long)GetLastError();
		}
	/* clean up (after success) */
	} else if (pass == 1) {
		if (wglGetCurrentContext() == rc)
			wglMakeCurrent(NULL, NULL);
		wglDeleteContext(rc); ReleaseDC(hndl, dc);
		DestroyWindow(hndl); UnregisterClass(class_name, instance);
		initialized = FALSE;
	/* clean up (after error) */
	} else if (pass < 0) {
		if (initialized) {
			if (wglGetCurrentContext() == rc)
				wglMakeCurrent(NULL, NULL);
			wglDeleteContext(rc); ReleaseDC(hndl, dc);
			DestroyWindow(hndl); UnregisterClass(class_name, instance);
			initialized = FALSE;
		}
	}
}
