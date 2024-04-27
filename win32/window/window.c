/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "window.h"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* Exported Go functions:                                   */
/* g2dResize                                                */
/* g2dKeyDown                                               */
/* g2dKeyUp                                                 */
/* g2dClose                                                 */

#define LOADER_ID "g2d.loader"
#define WINDOW_ID "g2d.window"

#define CLASS_NAME TEXT("vbsw_g2d_window")

#define LOAD_FUNC(t, n, e) if (cdata[0].err2 == 0) { n = (t) loader[0].load_func(#n, &cdata[0].err2); if (cdata[0].err2) cdata[0].err1 = e; }

/* from github.com/vbsw/golib/cdata/cdata.c */
typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

/* from github.com/vbsw/g2d/win32/loader.c */
typedef void* (loader_load_func_t) (const char *name, long long *err2);
typedef struct { loader_load_func_t *load_func; void *obj; } loader_t;

/* for external usage */

/* for internal usage */
typedef struct { HINSTANCE instance; HWND hndl; HDC dc; HGLRC rc; } obj_t;

typedef struct {
	HDC dc;
	HGLRC rc;
} context_t;

typedef struct {
	HWND hndl;
	context_t ctx;
} window_t;

typedef struct {
	int x, y, width, height;
} client_t;

typedef struct {
	int width_min, height_min, width_max, height_max;
	int borderless, dragable, fullscreen, resizable, locked;
	DWORD style;
} config_t;

typedef struct {
	int dragging, dragging_cust, locked;
	int minimized, maximized, resizing;
	int focus, shown;
} state_t;

typedef struct {
	window_t wnd;
	client_t client;
	client_t client_bak;
	config_t config;
	state_t state;
	int key_repeated[255];
	int cb_id;
	float projection_mat[4*4];
} window_data_t;

/* from wglext.h */
typedef BOOL(WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC hdc, const int *piAttribIList, const FLOAT *pfAttribFList, UINT nMaxFormats, int *piFormats, UINT *nNumFormats);
typedef HGLRC(WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC hDC, HGLRC hShareContext, const int *attribList);
typedef BOOL(WINAPI * PFNWGLSWAPINTERVALEXTPROC) (int interval);
typedef int (WINAPI * PFNWGLGETSWAPINTERVALEXTPROC) (void);

static PFNWGLCHOOSEPIXELFORMATARBPROC    wglChoosePixelFormatARB    = NULL;
static PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
static PFNWGLSWAPINTERVALEXTPROC         wglSwapIntervalEXT         = NULL;
static PFNWGLGETSWAPINTERVALEXTPROC      wglGetSwapIntervalEXT      = NULL;

static HINSTANCE instance = NULL;
static int windows_count = 0;

static LPCTSTR ensure_title(void *const title) {
	return title? (LPCTSTR)title : TEXT("OpenGL");
}

static void monitor_metrics(HMONITOR const monitor, int *const x, int *const y, int *const w, int *const h) {
	MONITORINFO mi = { sizeof(mi) }; GetMonitorInfo(monitor, &mi);
	x[0] = mi.rcMonitor.left; y[0] = mi.rcMonitor.top;
	w[0] = mi.rcMonitor.right - mi.rcMonitor.left; h[0] = mi.rcMonitor.bottom - mi.rcMonitor.top;
}

static void window_metrics(window_data_t *const wnd_data, int *const x, int *const y, int *const w, int *const h) {
	RECT rect = { wnd_data[0].client.x, wnd_data[0].client.y, wnd_data[0].client.x + wnd_data[0].client.width, wnd_data[0].client.y + wnd_data[0].client.height };
	AdjustWindowRect(&rect, wnd_data[0].config.style, FALSE);
	x[0] = rect.left; y[0] = rect.top; w[0] = rect.right - rect.left; h[0] = rect.bottom - rect.top;
}

static void style_update(window_data_t *const wnd_data) {
	if (wnd_data[0].config.borderless)
		if (wnd_data[0].config.resizable)
			wnd_data[0].config.style = WS_POPUP;
		else
			wnd_data[0].config.style = WS_POPUP;
	else
		if (wnd_data[0].config.resizable)
			wnd_data[0].config.style = WS_OVERLAPPEDWINDOW;
		else
			wnd_data[0].config.style = WS_OVERLAPPEDWINDOW & ~WS_THICKFRAME & ~WS_MAXIMIZEBOX;
}

static void client_update(window_data_t *const wnd_data) {
	POINT point = {0, 0};
	RECT rect = {0, 0, 0, 0};
	ClientToScreen(wnd_data[0].wnd.hndl, &point);
	GetClientRect(wnd_data[0].wnd.hndl, &rect);
	wnd_data[0].client.x = point.x;
	wnd_data[0].client.y = point.y;
	wnd_data[0].client.width = (int)(rect.right - rect.left);
	wnd_data[0].client.height = (int)(rect.bottom - rect.top);
}

static LRESULT CALLBACK windowProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam) {
	LRESULT result = 0;
	if (message == WM_NCCREATE) {
		window_data_t *const wnd_data = (window_data_t*)(((CREATESTRUCT*)lParam)->lpCreateParams);
		if (wnd_data)
			SetWindowLongPtr(hWnd, GWLP_USERDATA, (LONG_PTR)wnd_data);
		result = DefWindowProc(hWnd, message, wParam, lParam);
	} else {
		window_data_t *const wnd_data = (window_data_t*)GetWindowLongPtr(hWnd, GWLP_USERDATA);
		if (wnd_data) {
			switch (message) {
			case WM_SIZE:
				client_update(wnd_data);
				g2dResize(wnd_data[0].cb_id);
				result = DefWindowProc(hWnd, message, wParam, lParam);
				break;
			case WM_CLOSE:
				g2dClose(wnd_data[0].cb_id);
				break;
			case WM_KEYDOWN:
				if (!key_down_process(wnd_data, message, wParam, lParam))
					result = DefWindowProc(hWnd, message, wParam, lParam);
				break;
			case WM_KEYUP:
				if (!key_up_process(wnd_data, message, wParam, lParam))
					result = DefWindowProc(hWnd, message, wParam, lParam);
				break;
			default:
				result = DefWindowProc(hWnd, message, wParam, lParam);
			}
		} else {
			result = DefWindowProc(hWnd, message, wParam, lParam);
		}
	}
	return result;
}

void g2d_window_init(const int pass, cdata_t *const cdata) {
	if (pass == 0) {
		cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
		loader_t *const loader = get(cdata, LOADER_ID);
		if (loader) {
			LOAD_FUNC(PFNWGLCHOOSEPIXELFORMATARBPROC,    wglChoosePixelFormatARB,    1000201)
			LOAD_FUNC(PFNWGLCREATECONTEXTATTRIBSARBPROC, wglCreateContextAttribsARB, 1000202)
			LOAD_FUNC(PFNWGLSWAPINTERVALEXTPROC,         wglSwapIntervalEXT,         1000203)
			LOAD_FUNC(PFNWGLGETSWAPINTERVALEXTPROC,      wglGetSwapIntervalEXT,      1000204)
			if (cdata[0].err2 == 0) {
				instance = GetModuleHandle(NULL);
				if (instance) {
					obj_t *const obj = (obj_t*)malloc(sizeof(obj_t));
					if (obj) {
						cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
						set(cdata, (void*)obj, WINDOW_ID);
					} else {
						cdata[0].err1 = 30;
					}
				} else {
					cdata[0].err1 = 1000205; cdata[0].err2 = (long long)GetLastError();
				}
			}
		} else {
			cdata[0].err1 = 1000200;
		}
	} else if (pass < 0 || pass == 1) {
		cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
		obj_t *const obj = (obj_t*)get(cdata, WINDOW_ID);
		if (obj) {
			obj_t *const obj = (obj_t*)obj[0].obj;
			free(obj);
		} else if (pass == 1) {
			cdata[0].err1 = 1000213;
		}
	}
}

void g2d_window_create(void **const data, const int cb_id, const int x, const int y, const int w, const int h, const int wn, const int hn, const int wx, const int hx,
	const int b, const int d, const int r, const int f, const int l, const int c, void *t, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)malloc(sizeof(window_data_t));
	if (wnd_data) {
		ZeroMemory(wnd_data, sizeof(window_data_t));
		wnd_data[0].cb_id = cb_id;
		wnd_data[0].client.x = x;
		wnd_data[0].client.y = y;
		wnd_data[0].client.width = w;
		wnd_data[0].client.height = h;
		wnd_data[0].config.width_min = wn;
		wnd_data[0].config.height_min = hn;
		wnd_data[0].config.width_max = wx;
		wnd_data[0].config.height_max = hx;
		wnd_data[0].config.borderless = b;
		wnd_data[0].config.dragable = d;
		wnd_data[0].config.fullscreen = f;
		wnd_data[0].config.resizable = r;
		wnd_data[0].config.locked = l;
		style_update(wnd_data);
		if (c) {
			int wx, wy, ww, wh, mx, my, mw, mh;
			window_metrics(wnd_data, &wx, &wy, &ww, &wh);
			monitor_metrics(MonitorFromWindow(NULL, MONITOR_DEFAULTTOPRIMARY), &mx, &my, &mw, &mh);
			wnd_data[0].client.x = mx + (mw - ww) / 2 + (wnd_data[0].client.x - wx);
			wnd_data[0].client.y = my + (mh - wh) / 2 + (wnd_data[0].client.y - wy);
		}
		if (windows_count == 0) {
			WNDCLASSEX cls;
			ZeroMemory(&cls, sizeof(WNDCLASSEX));
			cls.cbSize = sizeof(WNDCLASSEX);
			cls.style = CS_OWNDC | CS_HREDRAW | CS_VREDRAW | CS_DBLCLKS;
			cls.lpfnWndProc = windowProc;
			cls.hInstance = instance;
			cls.hIcon = LoadIcon(NULL, IDI_WINLOGO);
			cls.hCursor = LoadCursor(NULL, IDC_ARROW);
			cls.lpszClassName = class_name;
			if (RegisterClassEx(&cls) != INVALID_ATOM) {
				windows_count++;
			} else {
				err1[0] = 13; err2[0] = (long long)GetLastError(); free(wnd_data); if (windows_count <= 0) PostQuitMessage(0);
			}
		} else {
			windows_count++;
		}
		if (err1[0] == 0) {
			int x, y, w, h; window_metrics(wnd_data, &x, &y, &w, &h);
			const DWORD style = wnd_data[0].config.style;
			wnd_data[0].wnd.hndl = CreateWindow(class_name, ensure_title(t), style, x, y, w, h, NULL, NULL, instance, (LPVOID)wnd_data);
			if (wnd_data[0].wnd.hndl) {
				wnd_data[0].wnd.ctx.dc = GetDC(wnd_data[0].wnd.hndl);
				if (wnd_data[0].wnd.ctx.dc) {
					int pixelFormat;
					BOOL status = FALSE;
					UINT numFormats = 0;
					const int pixelAttribs[] = {
						WGL_DRAW_TO_WINDOW_ARB, GL_TRUE,
						WGL_SUPPORT_OPENGL_ARB, GL_TRUE,
						WGL_DOUBLE_BUFFER_ARB, GL_TRUE,
						/* WGL_SWAP_COPY_ARB might have update problems in fullscreen */
						/* WGL_SWAP_EXCHANGE_ARB might have problems with start menu in fullscreen */
						WGL_SWAP_METHOD_ARB, WGL_SWAP_EXCHANGE_ARB,
						WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
						WGL_ACCELERATION_ARB, WGL_FULL_ACCELERATION_ARB,
						WGL_COLOR_BITS_ARB, 32,
						WGL_ALPHA_BITS_ARB, 8,
						WGL_DEPTH_BITS_ARB, 24,
						/* anti aliasing */
						//WGL_SAMPLE_BUFFERS_ARB, 1,
						//WGL_SAMPLES_ARB, 4,
						0
					};
					const int contextAttributes[] = {
						WGL_CONTEXT_MAJOR_VERSION_ARB, 3,
						WGL_CONTEXT_MINOR_VERSION_ARB, 0,
						WGL_CONTEXT_PROFILE_MASK_ARB, WGL_CONTEXT_CORE_PROFILE_BIT_ARB,
						0
					};
					status = wglChoosePixelFormatARB(wnd_data[0].wnd.ctx.dc, pixelAttribs, NULL, 1, &pixelFormat, &numFormats);
					if (status && numFormats) {
						PIXELFORMATDESCRIPTOR pfd;
						ZeroMemory(&pfd, sizeof(PIXELFORMATDESCRIPTOR));
						DescribePixelFormat(wnd_data[0].wnd.ctx.dc, pixelFormat, sizeof(PIXELFORMATDESCRIPTOR), &pfd);
						if (SetPixelFormat(wnd_data[0].wnd.ctx.dc, pixelFormat, &pfd)) {
							wnd_data[0].wnd.ctx.rc = wglCreateContextAttribsARB(wnd_data[0].wnd.ctx.dc, 0, contextAttributes);
							if (wnd_data[0].wnd.ctx.rc) {
								data[0] = (void*)wnd_data;
							} else {
								err1[0] = 18; err2[0] = (long long)GetLastError(); windows_count--;
								ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.ctx.dc); DestroyWindow(wnd_data[0].wnd.hndl);
								free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); PostQuitMessage(0); }
							}
						} else {
							err1[0] = 17; err2[0] = (long long)GetLastError(); windows_count--;
							ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.ctx.dc); DestroyWindow(wnd_data[0].wnd.hndl);
							free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); PostQuitMessage(0); }
						}
					} else {
						err1[0] = 16; err2[0] = (long long)GetLastError(); windows_count--;
						ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.ctx.dc); DestroyWindow(wnd_data[0].wnd.hndl);
						free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); PostQuitMessage(0); }
					}
				} else {
					err1[0] = 15; windows_count--; DestroyWindow(wnd_data[0].wnd.hndl);
					free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); PostQuitMessage(0); }
				}
			} else {
				err1[0] = 14; err2[0] = (long long)GetLastError(); windows_count--;
				free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); PostQuitMessage(0); }
			}
		}
	} else {
		err1[0] = 121;
	}
}

void g2d_window_show(void *data, long long *err1, long long *err2) {
	if (data) {
		window_data_t *const wnd_data = (window_data_t*)data;
		ShowWindow(wnd_data[0].wnd.hndl, SW_SHOWDEFAULT);
	}
}

void g2d_window_destroy(void *const data, long long *err1, long long *err2) {
	if (data) {
		window_data_t *const wnd_data = (window_data_t*)data;
		if (!wglDeleteContext(wnd_data[0].wnd.ctx.rc) && err1[0] == 0) {
			err1[0] = 20; err2[0] = (long long)GetLastError();
		}
		ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.ctx.dc);
		if (!DestroyWindow(wnd_data[0].wnd.hndl) && err1[0] == 0) {
			err1[0] = 21; err2[0] = (long long)GetLastError();
		}
		windows_count--; free(wnd_data);
		if (windows_count <= 0) { 
			if (!UnregisterClass(class_name, instance) && err1[0] == 0) {
				err1[0] = 22; err2[0] = (long long)GetLastError();
			}
			PostQuitMessage(0);
		}
	}
}

void g2d_window_props(void *const data, int *const x, int *const y, int *const w, int *const h, int *const wn, int *const hn,
	int *const wx, int *const hx, int *const b, int *const d, int *const r, int *const f, int *const l) {
	window_data_t *const wnd_data = (window_data_t*)data;
	*x = wnd_data[0].client.x;
	*y = wnd_data[0].client.y;
	*w = wnd_data[0].client.width;
	*h = wnd_data[0].client.height;
	*wn = wnd_data[0].config.width_min;
	*hn = wnd_data[0].config.height_min;
	*wx = wnd_data[0].config.width_max;
	*hx = wnd_data[0].config.height_max;
	*b = wnd_data[0].config.borderless;
	*d = wnd_data[0].config.dragable;
	*r = wnd_data[0].config.resizable;
	*f = wnd_data[0].config.fullscreen;
	*l = wnd_data[0].config.locked;
}

void g2d_window_to_tstr(void **const str, void *const go_cstr, const size_t length, long long *err1) {
	LPTSTR const str_new = (LPTSTR)malloc(sizeof(TCHAR) * (length + 1));
	if (str_new) {
		if (length > 0)
			#ifdef UNICODE
			MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, (const char*)go_cstr, length, str_new, length);
			#else
			memcpy(str_new, go_cstr, length);
			#endif
		str_new[length] = 0;
	} else {
		err1[0] = 120;
	}
	str[0] = (void*)str_new;
}

void g2d_window_free(void *const data) {
	free(data);
}
