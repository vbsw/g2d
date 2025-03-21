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
#include <gl/GL.h>
#include "window.h"

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

// from wgl.h
#define WGL_SAMPLE_BUFFERS_ARB            0x2041
#define WGL_SAMPLES_ARB                   0x2042
#define WGL_DRAW_TO_WINDOW_ARB            0x2001
#define WGL_SWAP_METHOD_ARB               0x2007
#define WGL_SUPPORT_OPENGL_ARB            0x2010
#define WGL_DOUBLE_BUFFER_ARB             0x2011
#define WGL_PIXEL_TYPE_ARB                0x2013
#define WGL_TYPE_RGBA_ARB                 0x202B
#define WGL_ACCELERATION_ARB              0x2003
#define WGL_FULL_ACCELERATION_ARB         0x2027
#define WGL_SWAP_EXCHANGE_ARB             0x2028
#define WGL_SWAP_COPY_ARB                 0x2029
#define WGL_SWAP_UNDEFINED_ARB            0x202A
#define WGL_COLOR_BITS_ARB                0x2014
#define WGL_ALPHA_BITS_ARB                0x201B
#define WGL_DEPTH_BITS_ARB                0x2022
#define WGL_STENCIL_BITS_ARB              0x2023
#define WGL_CONTEXT_MAJOR_VERSION_ARB     0x2091
#define WGL_CONTEXT_MINOR_VERSION_ARB     0x2092
#define WGL_CONTEXT_PROFILE_MASK_ARB      0x9126
#define WGL_CONTEXT_CORE_PROFILE_BIT_ARB  0x00000001

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

/*
typedef struct {
	GLuint id, vao, vbo, ebo, max_size;
	GLint pos_att, col_att, proj_unif;
	float *buffer;
} rect_program_t;

typedef struct {
	GLuint id, vao, vbo, ebo, max_size;
	GLint pos_att, col_att, tex_att, proj_unif, tex_unif;
	float *buffer;
} image_program_t;

typedef struct {
	GLuint id, vao, vbo, ebo, max_size;
	GLint pos_att, col_att, proj_unif;
	float *buffer;
} program_t;
*/

typedef struct {
	struct { HWND hndl; HDC dc; HGLRC rc; } wnd;
	struct { int x, y, width, height; } client;
	struct { int x, y, width, height; } client_bak;
	struct { int x, y, double_clicked[5]; } mouse;
	struct { int width_min, height_min, width_max, height_max, borderless, dragable, fullscreen, resizable, locked; DWORD style; } config;
	struct { int dragging, dragging_cust, locked, minimized, maximized, resizing, focus, shown; } state;
	int key_repeated[255];
	int cb_id;
/*
	program_t prog;
	rect_program_t rect_prog;
	image_program_t image_prog;
	float projection_mat[4*4];
*/
} window_data_t;

typedef struct {} window_t;

static PFNWGLCHOOSEPIXELFORMATARBPROC    wglChoosePixelFormatARB    = NULL;
static PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
static PFNWGLSWAPINTERVALEXTPROC         wglSwapIntervalEXT         = NULL;
static PFNWGLGETSWAPINTERVALEXTPROC      wglGetSwapIntervalEXT      = NULL;

static const WPARAM wnd_CUSTOM_EVENT = (WPARAM)"vbsw_g2dc";
static const WPARAM wnd_QUIT_EVENT   = (WPARAM)"vbsw_g2dq";
static LPCTSTR const class_name      = TEXT("vbsw_g2d_window");
static LPCTSTR const default_title   = TEXT("g2d - 0.1.0");

static BOOL initialized   = FALSE;
static HINSTANCE instance = NULL;
static int windows_count  = 0;
static DWORD thread_id    = 0;
static BOOL stop          = FALSE;

static window_t window;

static LPTSTR to_tstr(void *const go_cstr, const size_t length, long long *err1) {
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
	return str_new;
}

static void window_metrics(window_data_t *const wnd_data, int *const x, int *const y, int *const w, int *const h) {
	RECT rect = { wnd_data[0].client.x, wnd_data[0].client.y, wnd_data[0].client.x + wnd_data[0].client.width, wnd_data[0].client.y + wnd_data[0].client.height };
	AdjustWindowRect(&rect, wnd_data[0].config.style, FALSE);
	x[0] = rect.left; y[0] = rect.top; w[0] = rect.right - rect.left; h[0] = rect.bottom - rect.top;
}

static void monitor_metrics(HMONITOR const monitor, int *const x, int *const y, int *const w, int *const h) {
	MONITORINFO mi = { sizeof(mi) }; GetMonitorInfo(monitor, &mi);
	x[0] = mi.rcMonitor.left; y[0] = mi.rcMonitor.top;
	w[0] = mi.rcMonitor.right - mi.rcMonitor.left; h[0] = mi.rcMonitor.bottom - mi.rcMonitor.top;
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

static int keycode(const UINT message, const WPARAM wParam, const LPARAM lParam) {
	const int key = (int)(HIWORD(lParam) & 0xff);
	switch (key)
	{
	case 0:  return 0;
	case 1:  return 41;        // ESC         0x29
	case 2:  return 30;        // 1           0x1E
	case 3:  return 31;        // 2           0x1F
	case 4:  return 32;        // 3           0x20
	case 5:  return 33;        // 4           0x21
	case 6:  return 34;        // 5           0x22
	case 7:  return 35;        // 6           0x23
	case 8:  return 36;        // 7           0x24
	case 9:  return 37;        // 8           0x25
	case 10: return 38;        // 9           0x26
	case 11: return 39;        // 0           0x27
	case 12: return 45;        // -           0x2D
	case 13: return 46;        // =           0x2E
	case 14: return 42;        // DELETE      0x2A
	case 15: return 43;        // TAB         0x2B
	case 16: return 20;        // Q           0x14
	case 17: return 26;        // W           0x1A
	case 18: return 8;         // E           0x08
	case 19: return 21;        // R           0x15
	case 20: return 23;        // T           0x17
	case 21: return 28;        // Y           0x1C
	case 22: return 24;        // U           0x18
	case 23: return 12;        // I           0x0C
	case 24: return 18;        // O           0x12
	case 25: return 19;        // P           0x13
	case 26: return 47;        // [           0x2F
	case 27: return 48;        // ]           0x30
	case 28:
		if (HIWORD(lParam) >> 8 & 0x1)
			return 88;         // pad ENTER   0x58
		return 40;             // board ENTER 0x28
	case 29:
		if (wParam == VK_CONTROL) {
			if (HIWORD(lParam) >> 8 & 0x1)
				return 228;    // RCTRL       0xE4
			return 224;        // LCTRL       0xE0
		}
		return 0;
	case 30: return 4;         // A           0x04
	case 31: return 22;        // S           0x16
	case 32: return 7;         // D           0x07
	case 33: return 9;         // F           0x09
	case 34: return 10;        // G           0x0A
	case 35: return 11;        // H           0x0B
	case 36: return 13;        // J           0x0D
	case 37: return 14;        // K           0x0E
	case 38: return 15;        // L           0x0F
	case 39: return 51;        // ;           0x33
	case 40: return 52;        // '           0x34
	case 41: return 53;        // ^           0x35
	case 42: return 225;       // LSHIFT      0xE1
	case 43: return 50;        // ~           0x32
	case 44: return 29;        // Z           0x1D
	case 45: return 27;        // X           0x1B
	case 46: return 6;         // C           0x06
	case 47: return 25;        // V           0x19
	case 48: return 5;         // B           0x05
	case 49: return 17;        // N           0x11
	case 50: return 16;        // M           0x10
	case 51: return 54;        // ,           0x36
	case 52: return 55;        // .           0x37
	case 53:
		if (wParam == VK_DIVIDE)
			return 84;         // pad /       0x54
		return 56;             // /           0x38
	case 54: return 229;       // RSHIFT      0xE5
	case 55: return 85;        // pad *       0x55
	case 56:
		if (message == WM_SYSKEYDOWN || message == WM_SYSKEYUP)
			return 226;        // LALT        0xE2
		return 230;            // RALT        0xE6
	case 57: return 44;        // SPACE       0x2C
	case 58: return 57;        // CAPS        0x39
	case 59: return 58;        // F1          0x3A
	case 60: return 59;        // F2          0x3B
	case 61: return 60;        // F3          0x3C
	case 62: return 61;        // F4          0x3D
	case 63: return 62;        // F5          0x3E
	case 64: return 63;        // F6          0x3F
	case 65: return 64;        // F7          0x40
	case 66: return 65;        // F8          0x41
	case 67: return 66;        // F9          0x42
	case 68: return 67;        // F10         0x43
	case 69:
		if (wParam == VK_PAUSE)
			return 72;         // PAUSE       0x48
		return 83;             // pad LOCK    0x53
	case 70: return 71;        // SCROLL      0x47
	case 71:
		if (wParam == VK_HOME)
			return 74;         // HOME        0x4A
		return 95;             // pad 7       0x5F
	case 72:
		if (wParam == VK_UP)
			return 82;         // UP          0x52
		return 96;             // pad 8       0x60
	case 73:
		if (wParam == VK_PRIOR)
			return 75;         // PAGEUP      0x4B
		return 97;             // pad 9       0x61
	case 74: return 86;        // pad -       0x56
	case 75:
		if (wParam == VK_LEFT)
			return 80;         // LEFT        0x50
		return 92;             // pad 4       0x5C
	case 76: return 93;        // pad 5       0x5D
	case 77:
		if (wParam == VK_RIGHT)
			return 79;         // RIGHT       0x4F
		return 94;             // pad 6       0x5E
	case 78: return 87;        // pad +       0x57
	case 79:
		if (wParam == VK_END)
			return 77;         // END         0x4D
		return 89;             // pad 1       0x59
	case 80:
		if (wParam == VK_DOWN)
			return 81;         // DOWN        0x51
		return 90;             // pad 2       0x5A
	case 81:
		if (wParam == VK_NEXT)
			return 78;         // PAGEDOWN    0x4E
		return 91;             // pad 3       0x5B
	case 82: return 73;        // INSERT      0x49
	case 83:
		if (wParam == VK_DELETE)
			return 76;         // DELETE F    0x4C
		return 99;             // pad DELETE  0x63
	case 84: return 0;
	case 85: return 0;
	case 86: return 100;       // |           0x64
	case 87: return 68;        // F11         0x44
	case 88: return 69;        // F12         0x45
	case 89: return 0;         // LWIN        0xE3
	case 90: return 0;         // RWIN        0xE7
	case 91: return 0;
	case 92: return 0;
	case 93: return 118;       // MENU        0x76
	}
	return key;
}

static BOOL key_down_process(window_data_t *const wnd_data, const UINT message, const WPARAM wParam, const LPARAM lParam) {
	const int code = keycode(message, wParam, lParam);
	if (code) {
		g2dKeyDown(wnd_data[0].cb_id, code, wnd_data[0].key_repeated[code]++, 0);
		return TRUE;
	}
	return FALSE;
}

static BOOL key_up_process(window_data_t *const wnd_data, const UINT message, const WPARAM wParam, const LPARAM lParam) {
	const int code = keycode(message, wParam, lParam);
	if (code) {
		wnd_data[0].key_repeated[code] = 0;
		g2dKeyUp(wnd_data[0].cb_id, code, 0);
		return TRUE;
	}
	return FALSE;
}

void g2d_window_mainloop() {
	MSG msg; BOOL ret_code; stop = FALSE;
	thread_id = GetCurrentThreadId();
	g2dWindowMainLoopStart(0);
	while (!stop && (ret_code = GetMessage(&msg, NULL, 0, 0)) > 0) {
		if (msg.message == WM_APP) {
			if (msg.wParam == wnd_CUSTOM_EVENT)
				g2dWindowMainLoopEvent();
			else if (msg.wParam == wnd_QUIT_EVENT)
				break;
		} else {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
}

void g2d_window_post_custom_msg(long long *const millis, long long *const err1, long long *const err2) {
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
			if (!wnd_data[0].state.minimized) {
				switch (message) {
/*
				case WM_MOVE:
					client_update(wnd_data);
					result = DefWindowProc(hWnd, message, wParam, lParam);
					g2dWindowMove(wnd_data[0].cb_id);
					break;
				case WM_SIZE:
					// avoid resize on show
					if (wnd_data[0].state.shown != 0) {
						client_update(wnd_data);
						g2dWindowResize(wnd_data[0].cb_id);
					} else {
						wnd_data[0].state.shown = 1;
					}
					result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
				case WM_CLOSE:
					g2dClose(wnd_data[0].cb_id);
					break;
*/
				case WM_KEYDOWN:
					if (!key_down_process(wnd_data, message, wParam, lParam))
						result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
				case WM_KEYUP:
					if (!key_up_process(wnd_data, message, wParam, lParam))
						result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
/*
				case WM_SYSCOMMAND:
					if (wParam == SC_MINIMIZE) {
						wnd_data[0].state.minimized = 1;
						g2dWindowMinimize(wnd_data[0].cb_id);
					}
					result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
				case WM_MOUSEMOVE:
					wnd_data[0].mouse.x = ((int)(short)LOWORD(lParam));
					wnd_data[0].mouse.y = ((int)(short)HIWORD(lParam));
					g2dMouseMove(wnd_data[0].cb_id);
					result = DefWindowProc(hWnd, message, wParam, lParam);
*/


/*
					if (state.dragging_cust && !state.maximized) {
						move_window(client.x + (int)(short)LOWORD(lParam) - mouse.x, client.y + (int)(short)HIWORD(lParam) - mouse.y, client.width, client.height);
					} else {
						mouse.x = ((int)(short)LOWORD(lParam));
						mouse.y = ((int)(short)HIWORD(lParam));
						result = DefWindowProc(hWnd, message, wParam, lParam);
					}
					if (config.locked && !state.locked && state.focus)
						update_clip_cursor();
*/

/*
					break;
				case WM_LBUTTONDOWN:
					button_down(wnd_data, 0, 0);
					break;
				case WM_LBUTTONUP:
					button_up(wnd_data, 0);
					break;
				case WM_LBUTTONDBLCLK:
					button_down(wnd_data, 0, 1);
					break;
				case WM_RBUTTONDOWN:
					button_down(wnd_data, 1, 0);
					break;
				case WM_RBUTTONUP:
					button_up(wnd_data, 1);
					break;
				case WM_RBUTTONDBLCLK:
					button_down(wnd_data, 1, 1);
					break;
				case WM_MBUTTONDOWN:
					button_down(wnd_data, 2, 0);
					break;
				case WM_MBUTTONUP:
					button_up(wnd_data, 2);
					break;
				case WM_MBUTTONDBLCLK:
					button_down(wnd_data, 2, 1);
					break;
				case WM_MOUSEWHEEL:
					g2dWheel(wnd_data[0].cb_id, (float)GET_WHEEL_DELTA_WPARAM(wParam) / (float)WHEEL_DELTA);
					break;
				case WM_XBUTTONDOWN:
					if (HIWORD(wParam) == XBUTTON1)
						button_down(wnd_data, 3, 0);
					else if (HIWORD(wParam) == XBUTTON2)
						button_down(wnd_data, 4, 0);
					break;
				case WM_XBUTTONUP:
					if (HIWORD(wParam) == XBUTTON1)
						button_up(wnd_data, 3);
					else if (HIWORD(wParam) == XBUTTON2)
						button_up(wnd_data, 4);
					break;
				case WM_XBUTTONDBLCLK:
					if (HIWORD(wParam) == XBUTTON1)
						button_down(wnd_data, 3, 1);
					else if (HIWORD(wParam) == XBUTTON2)
						button_down(wnd_data, 4, 1);
					break;
*/
				default:
					result = DefWindowProc(hWnd, message, wParam, lParam);
				}
			} else {
				result = DefWindowProc(hWnd, message, wParam, lParam);
				if (message == WM_SETFOCUS) {
					// restore from minimized and avoid move/resize events
/*
					if (wnd_data[0].state.minimized) {
						wnd_data[0].state.minimized = 0;
						g2dWindowRestore(wnd_data[0].cb_id);
					}
*/
				}
			}
		} else {
			result = DefWindowProc(hWnd, message, wParam, lParam);
		}
	}
	return result;
}

void g2d_window_create(void **const data, const int cb_id, const int x, const int y, const int w, const int h, const int wn, const int hn, const int wx, const int hx,
	const int b, const int d, const int r, const int f, const int l, const int c, void *const t, const size_t ts, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)malloc(sizeof(window_data_t));
	if (wnd_data) {
		LPCTSTR const title = (ts == 0) ? default_title : to_tstr(t, ts, err1);
		if (err1[0] == 0) {
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
					err1[0] = 13; err2[0] = (long long)GetLastError(); free(wnd_data);
				}
			} else {
				windows_count++;
			}
			if (err1[0] == 0) {
				int x, y, w, h; window_metrics(wnd_data, &x, &y, &w, &h);
				const DWORD style = wnd_data[0].config.style;
				wnd_data[0].wnd.hndl = CreateWindow(class_name, title, style, x, y, w, h, NULL, NULL, instance, (LPVOID)wnd_data);
				if (wnd_data[0].wnd.hndl) {
					wnd_data[0].wnd.dc = GetDC(wnd_data[0].wnd.hndl);
					if (wnd_data[0].wnd.dc) {
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
						status = wglChoosePixelFormatARB(wnd_data[0].wnd.dc, pixelAttribs, NULL, 1, &pixelFormat, &numFormats);
						if (status && numFormats) {
							PIXELFORMATDESCRIPTOR pfd;
							ZeroMemory(&pfd, sizeof(PIXELFORMATDESCRIPTOR));
							DescribePixelFormat(wnd_data[0].wnd.dc, pixelFormat, sizeof(PIXELFORMATDESCRIPTOR), &pfd);
							if (SetPixelFormat(wnd_data[0].wnd.dc, pixelFormat, &pfd)) {
								wnd_data[0].wnd.rc = wglCreateContextAttribsARB(wnd_data[0].wnd.dc, 0, contextAttributes);
								if (wnd_data[0].wnd.rc) {
									data[0] = (void*)wnd_data;
								} else {
									err1[0] = 18; err2[0] = (long long)GetLastError(); windows_count--;
									ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.dc); DestroyWindow(wnd_data[0].wnd.hndl);
									free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
								}
							} else {
								err1[0] = 17; err2[0] = (long long)GetLastError(); windows_count--;
								ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.dc); DestroyWindow(wnd_data[0].wnd.hndl);
								free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
							}
						} else {
							err1[0] = 16; err2[0] = (long long)GetLastError(); windows_count--;
							ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.dc); DestroyWindow(wnd_data[0].wnd.hndl);
							free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
						}
					} else {
						err1[0] = 15; windows_count--; DestroyWindow(wnd_data[0].wnd.hndl);
						free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
					}
				} else {
					err1[0] = 14; err2[0] = (long long)GetLastError(); windows_count--;
					free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
				}
			}
			if (title != default_title)
				free((void*)title);
		} else {
			err1[0] = 122;
		}
	} else {
		err1[0] = 121;
	}
}

void g2d_window_show(void *const data, long long *const err1, long long *const err2) {
	if (data) {
		window_data_t *const wnd_data = (window_data_t*)data;
		ShowWindow(wnd_data[0].wnd.hndl, SW_SHOWDEFAULT);
	}
}

#elif defined(G2D_WINDOW_LINUX)

#endif
