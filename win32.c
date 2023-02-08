/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#if defined(G2D_WIN32)

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>
#include "g2d.h"

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* Exported functions from Go are:                          */
/* g2dProcessMessage                                        */

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

#define WGL_SWAP_METHOD_EXT               0x2007
#define WGL_SWAP_EXCHANGE_EXT             0x2028
#define WGL_SWAP_COPY_EXT                 0x2029
#define WGL_SWAP_UNDEFINED_EXT            0x202A

// copied from glcorearb.h
#define GL_ARRAY_BUFFER                   0x8892
#define GL_STATIC_DRAW                    0x88E4
#define GL_FRAGMENT_SHADER                0x8B30
#define GL_VERTEX_SHADER                  0x8B31
#define GL_COMPILE_STATUS                 0x8B81
#define GL_INFO_LOG_LENGTH                0x8B84
#define GL_LINK_STATUS                    0x8B82
#define GL_VALIDATE_STATUS                0x8B83
#define GL_CLAMP_TO_BORDER                0x812D

/* wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions). */
#define LOAD_FUNC(t, n, e) if (err_num[0] == 0) { PROC const proc = wglGetProcAddress(#n); const DWORD last_error = GetLastError(); if (last_error == 0) n = (t) proc; else { err_num[0] = e; err_win32[0] = (g2d_ul_t)last_error; }}

/* from wglext.h */
typedef BOOL(WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC hdc, const int *piAttribIList, const FLOAT *pfAttribFList, UINT nMaxFormats, int *piFormats, UINT *nNumFormats);
typedef HGLRC(WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC hDC, HGLRC hShareContext, const int *attribList);
typedef BOOL(WINAPI * PFNWGLSWAPINTERVALEXTPROC) (int interval);
typedef int (WINAPI * PFNWGLGETSWAPINTERVALEXTPROC) (void);

// from glcorearb.h
typedef char GLchar;
typedef ptrdiff_t GLsizeiptr;
typedef GLuint(APIENTRY *PFNGLCREATESHADERPROC) (GLenum type);
typedef void (APIENTRY *PFNGLSHADERSOURCEPROC) (GLuint shader, GLsizei count, const GLchar *const*string, const GLint *length);
typedef void (APIENTRY *PFNGLCOMPILESHADERPROC) (GLuint shader);
typedef void (APIENTRY *PFNGLGETSHADERIVPROC) (GLuint shader, GLenum pname, GLint *params);
typedef void (APIENTRY *PFNGLGETSHADERINFOLOGPROC) (GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
typedef GLuint(APIENTRY *PFNGLCREATEPROGRAMPROC) (void);
typedef void (APIENTRY *PFNGLATTACHSHADERPROC) (GLuint program, GLuint shader);
typedef void (APIENTRY *PFNGLLINKPROGRAMPROC) (GLuint program);
typedef void (APIENTRY *PFNGLVALIDATEPROGRAMPROC) (GLuint program);
typedef void (APIENTRY *PFNGLGETPROGRAMIVPROC) (GLuint program, GLenum pname, GLint *params);
typedef void (APIENTRY *PFNGLGETPROGRAMINFOLOGPROC) (GLuint program, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
typedef void (APIENTRY *PFNGLGENBUFFERSPROC) (GLsizei n, GLuint *buffers);
typedef void (APIENTRY *PFNGLGENVERTEXARRAYSPROC) (GLsizei n, GLuint *arrays);
typedef GLint(APIENTRY *PFNGLGETATTRIBLOCATIONPROC) (GLuint program, const GLchar *name);
typedef void (APIENTRY *PFNGLBINDVERTEXARRAYPROC) (GLuint array);
typedef void (APIENTRY *PFNGLENABLEVERTEXATTRIBARRAYPROC) (GLuint index);
typedef void (APIENTRY *PFNGLVERTEXATTRIBPOINTERPROC) (GLuint index, GLint size, GLenum type, GLboolean normalized, GLsizei stride, const GLvoid *pointer);
typedef void (APIENTRY *PFNGLBINDBUFFERPROC) (GLenum target, GLuint buffer);
typedef void (APIENTRY *PFNGLBUFFERDATAPROC) (GLenum target, GLsizeiptr size, const GLvoid *data, GLenum usage);
typedef void (APIENTRY *PFNGLGETVERTEXATTRIBPOINTERVPROC) (GLuint index, GLenum pname, GLvoid **pointer);
typedef void (APIENTRY *PFNGLUSEPROGRAMPROC) (GLuint program);
typedef void (APIENTRY *PFNGLDELETEVERTEXARRAYSPROC) (GLsizei n, const GLuint *arrays);
typedef void (APIENTRY *PFNGLDELETEBUFFERSPROC) (GLsizei n, const GLuint *buffers);
typedef void (APIENTRY *PFNGLDELETEPROGRAMPROC) (GLuint program);
typedef void (APIENTRY *PFNGLDELETESHADERPROC) (GLuint shader);
typedef GLint(APIENTRY *PFNGLGETUNIFORMLOCATIONPROC) (GLuint program, const GLchar *name);
typedef void (APIENTRY *PFNGLUNIFORM1FVPROC) (GLint location, GLsizei count, const GLfloat *value);
typedef void (APIENTRY *PFNGLUNIFORMMATRIX4FVPROC) (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
typedef void (APIENTRY *PFNGLUNIFORMMATRIX3FVPROC) (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
typedef void (APIENTRY *PFNGLUNIFORMMATRIX2X3FVPROC) (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
typedef void (APIENTRY *PFNGLACTIVETEXTUREPROC) (GLenum texture);
typedef void (APIENTRY *PFNGLGENERATEMIPMAPPROC) (GLenum target);

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
	GLuint vs_id, fs_id, id;
	GLint position_att, projection_unif, data_unif;
} program_t;

typedef struct {
	window_t wnd;
	client_t client;
	client_t client_bak;
	config_t config;
	state_t state;
	int key_repeated[255];
	int cb_id;
	GLuint fs_id;
	program_t prog_rect;
	GLuint vao, vbo;
	float mat_projection[4*4];
} window_data_t;

static const WPARAM const g2d_EVENT = (WPARAM)"g2d";
static LPCTSTR const class_name = TEXT("g2d");
static LPCTSTR const class_name_dummy = TEXT("g2d_dummy");

static HINSTANCE instance = NULL;
static DWORD thread_id;
static BOOL thread_id_set = FALSE;
static BOOL initialized = FALSE;
static int windows_count = 0;

static PFNWGLCHOOSEPIXELFORMATARBPROC    wglChoosePixelFormatARB    = NULL;
static PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
static PFNWGLSWAPINTERVALEXTPROC         wglSwapIntervalEXT         = NULL;
static PFNWGLGETSWAPINTERVALEXTPROC      wglGetSwapIntervalEXT      = NULL;

static PFNGLCREATESHADERPROC             glCreateShader             = NULL;
static PFNGLSHADERSOURCEPROC             glShaderSource             = NULL;
static PFNGLCOMPILESHADERPROC            glCompileShader            = NULL;
static PFNGLGETSHADERIVPROC              glGetShaderiv              = NULL;
static PFNGLGETSHADERINFOLOGPROC         glGetShaderInfoLog         = NULL;
static PFNGLCREATEPROGRAMPROC            glCreateProgram            = NULL;
static PFNGLATTACHSHADERPROC             glAttachShader             = NULL;
static PFNGLLINKPROGRAMPROC              glLinkProgram              = NULL;
static PFNGLVALIDATEPROGRAMPROC          glValidateProgram          = NULL;
static PFNGLGETPROGRAMIVPROC             glGetProgramiv             = NULL;
static PFNGLGETPROGRAMINFOLOGPROC        glGetProgramInfoLog        = NULL;
static PFNGLGENBUFFERSPROC               glGenBuffers               = NULL;
static PFNGLGENVERTEXARRAYSPROC          glGenVertexArrays          = NULL;
static PFNGLGETATTRIBLOCATIONPROC        glGetAttribLocation        = NULL;
static PFNGLBINDVERTEXARRAYPROC          glBindVertexArray          = NULL;
static PFNGLENABLEVERTEXATTRIBARRAYPROC  glEnableVertexAttribArray  = NULL;
static PFNGLVERTEXATTRIBPOINTERPROC      glVertexAttribPointer      = NULL;
static PFNGLBINDBUFFERPROC               glBindBuffer               = NULL;
static PFNGLBUFFERDATAPROC               glBufferData               = NULL;
static PFNGLGETVERTEXATTRIBPOINTERVPROC  glGetVertexAttribPointerv  = NULL;
static PFNGLUSEPROGRAMPROC               glUseProgram               = NULL;
static PFNGLDELETEVERTEXARRAYSPROC       glDeleteVertexArrays       = NULL;
static PFNGLDELETEBUFFERSPROC            glDeleteBuffers            = NULL;
static PFNGLDELETEPROGRAMPROC            glDeleteProgram            = NULL;
static PFNGLDELETESHADERPROC             glDeleteShader             = NULL;
static PFNGLGETUNIFORMLOCATIONPROC       glGetUniformLocation       = NULL;
static PFNGLUNIFORM1FVPROC               glUniform1fv               = NULL;
static PFNGLUNIFORMMATRIX4FVPROC         glUniformMatrix4fv         = NULL;
static PFNGLUNIFORMMATRIX3FVPROC         glUniformMatrix3fv         = NULL;
static PFNGLUNIFORMMATRIX2X3FVPROC       glUniformMatrix2x3fv       = NULL;
static PFNGLGENERATEMIPMAPPROC           glGenerateMipmap           = NULL;
static PFNGLACTIVETEXTUREPROC            glActiveTexture            = NULL;

void g2d_free(void *const data) {
	free(data);
}

void g2d_test(int *array) {
	goDebug(array[0], array[0], 0, 0);
}

void g2d_to_tstr(void **const str, void *const go_cstr, int *const err_num) {
	LPTSTR str_new = NULL;
	size_t length;
	if (go_cstr)
		length = strlen(go_cstr);
	else
		length = 0;
	#ifdef UNICODE
	str_new = (LPTSTR)malloc(sizeof(WCHAR) * (length + 1));
	if (str_new) {
		if (length > 0)
			MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS, (const char*)go_cstr, length, str_new, length);
	#else
	str_new = (LPTSTR)malloc(sizeof(char) * (length + 1));
	if (str_new) {
		if (length > 0)
			memcpy(str_new, go_cstr, length);
	#endif
		str_new[length] = 0;
	}
	else
		err_num[0] = 120;
	str[0] = (void*)str_new;
}

#include "win32_graphics.h"
#include "win32_keys.h"
#include "win32_window.h"

void g2d_init(int *const err_num, g2d_ul_t *const err_win32) {
	if (!initialized) {
		/* module */
		instance = GetModuleHandle(NULL);
		if (instance) {
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
										/* wgl functions */
										LOAD_FUNC(PFNWGLCHOOSEPIXELFORMATARBPROC, wglChoosePixelFormatARB, 200)
										LOAD_FUNC(PFNWGLCREATECONTEXTATTRIBSARBPROC, wglCreateContextAttribsARB, 201)
										LOAD_FUNC(PFNWGLSWAPINTERVALEXTPROC, wglSwapIntervalEXT, 202)
										LOAD_FUNC(PFNWGLGETSWAPINTERVALEXTPROC, wglGetSwapIntervalEXT, 203)
										/* ogl functions */
										LOAD_FUNC(PFNGLCREATESHADERPROC, glCreateShader, 204)
										LOAD_FUNC(PFNGLSHADERSOURCEPROC, glShaderSource, 205)
										LOAD_FUNC(PFNGLCOMPILESHADERPROC, glCompileShader, 206)
										LOAD_FUNC(PFNGLGETSHADERIVPROC, glGetShaderiv, 207)
										LOAD_FUNC(PFNGLGETSHADERINFOLOGPROC, glGetShaderInfoLog, 208)
										LOAD_FUNC(PFNGLCREATEPROGRAMPROC, glCreateProgram, 209)
										LOAD_FUNC(PFNGLATTACHSHADERPROC, glAttachShader, 210)
										LOAD_FUNC(PFNGLLINKPROGRAMPROC, glLinkProgram, 211)
										LOAD_FUNC(PFNGLVALIDATEPROGRAMPROC, glValidateProgram, 212)
										LOAD_FUNC(PFNGLGETPROGRAMIVPROC, glGetProgramiv, 213)
										LOAD_FUNC(PFNGLGETPROGRAMINFOLOGPROC, glGetProgramInfoLog, 214)
										LOAD_FUNC(PFNGLGENBUFFERSPROC, glGenBuffers, 215)
										LOAD_FUNC(PFNGLGENVERTEXARRAYSPROC, glGenVertexArrays, 216)
										LOAD_FUNC(PFNGLGETATTRIBLOCATIONPROC, glGetAttribLocation, 217)
										LOAD_FUNC(PFNGLBINDVERTEXARRAYPROC, glBindVertexArray, 218)
										LOAD_FUNC(PFNGLENABLEVERTEXATTRIBARRAYPROC, glEnableVertexAttribArray, 219)
										LOAD_FUNC(PFNGLVERTEXATTRIBPOINTERPROC, glVertexAttribPointer, 220)
										LOAD_FUNC(PFNGLBINDBUFFERPROC, glBindBuffer, 221)
										LOAD_FUNC(PFNGLBUFFERDATAPROC, glBufferData, 222)
										LOAD_FUNC(PFNGLGETVERTEXATTRIBPOINTERVPROC, glGetVertexAttribPointerv, 223)
										LOAD_FUNC(PFNGLUSEPROGRAMPROC, glUseProgram, 224)
										LOAD_FUNC(PFNGLDELETEVERTEXARRAYSPROC, glDeleteVertexArrays, 225)
										LOAD_FUNC(PFNGLDELETEBUFFERSPROC, glDeleteBuffers, 226)
										LOAD_FUNC(PFNGLDELETEPROGRAMPROC, glDeleteProgram, 227)
										LOAD_FUNC(PFNGLDELETESHADERPROC, glDeleteShader, 228)
										LOAD_FUNC(PFNGLGETUNIFORMLOCATIONPROC, glGetUniformLocation, 229)
										LOAD_FUNC(PFNGLUNIFORMMATRIX3FVPROC, glUniformMatrix3fv, 230)
										LOAD_FUNC(PFNGLUNIFORM1FVPROC, glUniform1fv, 231)
										LOAD_FUNC(PFNGLUNIFORMMATRIX4FVPROC, glUniformMatrix4fv, 232)
										LOAD_FUNC(PFNGLUNIFORMMATRIX2X3FVPROC, glUniformMatrix2x3fv, 233)
										LOAD_FUNC(PFNGLGENERATEMIPMAPPROC, glGenerateMipmap, 234)
										LOAD_FUNC(PFNGLACTIVETEXTUREPROC, glActiveTexture, 235)
										/* destroy dummy */
										if (!wglMakeCurrent(NULL, NULL) && err_num[0] == 0) {
											err_num[0] = 9; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										if (!wglDeleteContext(dummy_rc) && err_num[0] == 0) {
											err_num[0] = 10; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										ReleaseDC(dummy_hndl, dummy_dc);
										if (!DestroyWindow(dummy_hndl) && err_num[0] == 0) {
											err_num[0] = 11; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										if (!UnregisterClass(class_name_dummy, instance) && err_num[0] == 0) {
											err_num[0] = 12; err_win32[0] = (g2d_ul_t)GetLastError();
										}
										initialized = (BOOL)(err_num[0] == 0);
									} else {
										err_num[0] = 8; err_win32[0] = (g2d_ul_t)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
									}
								} else {
									err_num[0] = 7; err_win32[0] = (g2d_ul_t)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
								}
							} else {
								err_num[0] = 6; err_win32[0] = (g2d_ul_t)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
							}
						} else {
							err_num[0] = 5; err_win32[0] = (g2d_ul_t)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
						}
					} else {
						err_num[0] = 4;
						DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
					}
				} else {
					err_num[0] = 3; err_win32[0] = (g2d_ul_t)GetLastError();
					UnregisterClass(class_name_dummy, instance);
				}
			} else {
				err_num[0] = 2; err_win32[0] = (g2d_ul_t)GetLastError();
			}
		} else {
			err_num[0] = 1; err_win32[0] = (g2d_ul_t)GetLastError();
		}
	}
}

void g2d_process_messages() {
	MSG msg; BOOL ret_code;
	thread_id = GetCurrentThreadId(); thread_id_set = TRUE;
	while ((ret_code = GetMessage(&msg, NULL, 0, 0)) > 0) {
		if (msg.message == WM_APP && msg.wParam == g2d_EVENT) {
			g2dProcessMessage();
		} else {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
	thread_id_set = FALSE;
}

void g2d_post_message(int *const err_num, g2d_ul_t *const err_win32) {
	if (thread_id_set) {
		if (!PostThreadMessage(thread_id, WM_APP, g2d_EVENT, 0)) {
			err_num[0] = 82; err_win32[0] = (g2d_ul_t)GetLastError();
		}
	} else {
		if (!PostMessage(NULL, WM_APP, g2d_EVENT, 0)) {
			err_num[0] = 82; err_win32[0] = (g2d_ul_t)GetLastError();
		}
	}
}

void g2d_quit_message_queue() {
	PostQuitMessage(0);
}

void g2d_context_make_current(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (!wglMakeCurrent(wnd_data[0].wnd.ctx.dc, wnd_data[0].wnd.ctx.rc)) {
		err_num[0] = 56; err_win32[0] = (g2d_ul_t)GetLastError();
	}
}

void g2d_context_release(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wnd_data[0].wnd.ctx.rc == wglGetCurrentContext() && !wglMakeCurrent(NULL, NULL)) {
		err_num[0] = 19; err_win32[0] = (g2d_ul_t)GetLastError();
	}
}

/* #if defined(G2D_WIN32) */
#endif
