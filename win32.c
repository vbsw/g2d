/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#if defined(G2D_WIN32)

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>
#include "g2d.h"
#include "win32_errors.h"

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

#define WGL_SWAP_METHOD_EXT               0x2007
#define WGL_SWAP_EXCHANGE_EXT             0x2028
#define WGL_SWAP_COPY_EXT                 0x2029
#define WGL_SWAP_UNDEFINED_EXT            0x202A

// copied from glcorearb.h
#define GL_TEXTURE0                       0x84C0
#define GL_ARRAY_BUFFER                   0x8892
#define GL_ELEMENT_ARRAY_BUFFER           0x8893
#define GL_STATIC_DRAW                    0x88E4
#define GL_DYNAMIC_DRAW                   0x88E8
#define GL_FRAGMENT_SHADER                0x8B30
#define GL_VERTEX_SHADER                  0x8B31
#define GL_COMPILE_STATUS                 0x8B81
#define GL_INFO_LOG_LENGTH                0x8B84
#define GL_LINK_STATUS                    0x8B82
#define GL_VALIDATE_STATUS                0x8B83
#define GL_CLAMP_TO_BORDER                0x812D
#define GL_MAX_TEXTURE_IMAGE_UNITS        0x8872

/* from wglext.h */
typedef BOOL(WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC hdc, const int *piAttribIList, const FLOAT *pfAttribFList, UINT nMaxFormats, int *piFormats, UINT *nNumFormats);
typedef HGLRC(WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC hDC, HGLRC hShareContext, const int *attribList);
typedef BOOL(WINAPI * PFNWGLSWAPINTERVALEXTPROC) (int interval);
typedef int (WINAPI * PFNWGLGETSWAPINTERVALEXTPROC) (void);

// from glcorearb.h
typedef char GLchar;
typedef ptrdiff_t GLsizeiptr;
typedef ptrdiff_t GLintptr;
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
typedef void (APIENTRY *PFNGLBUFFERSUBDATAPROC) (GLenum target, GLintptr offset, GLsizeiptr size, const void *data);
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
static PFNGLBUFFERSUBDATAPROC            glBufferSubData            = NULL;
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

static const WPARAM g2d_REQUEST_EVENT  = (WPARAM)"g2dc";
static const WPARAM g2d_QUIT_EVENT    = (WPARAM)"g2dq";
static LPCTSTR const class_name       = TEXT("g2d");
static LPCTSTR const class_name_dummy = TEXT("g2d_dummy");

static HINSTANCE instance = NULL;
static BOOL initialized   = FALSE;
static int windows_count  = 0;
static DWORD thread_id    = 0;
static BOOL stop          = FALSE;

/*
static struct {
	int count;
	BOOL force_destroy;
} active_windows = {0, FALSE};
*/

void g2d_free(void *const data) {
	free(data);
}

#include "win32_init.h"
#include "win32_main_loop.h"
#include "win32_window.h"

void g2d_post_request(long long *const err1, long long *const err2) {
	if (!PostThreadMessage(thread_id, WM_APP, g2d_REQUEST_EVENT, 0)) {
		err1[0] = 3999;
		err2[0] = (long long)GetLastError();
	}
}

void g2d_post_quit(long long *const err1, long long *const err2) {
	if (!PostThreadMessage(thread_id, WM_APP, g2d_QUIT_EVENT, 0)) {
		err1[0] = 3999;
		err2[0] = (long long)GetLastError();
	}
	stop = TRUE;
}

void g2d_clean_up() {
	MSG msg;
	while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE));
}

/*
static void *error_new(const int err_num, const DWORD err_win32, char *const err_str) {
	error_t *const err = (error_t*)malloc(sizeof(error_t));
	if (err) {
		err[0].err_num = err_num;
		err[0].err_win32 = (g2d_ul_t)err_win32;
		err[0].err_str = err_str;
		return (void*)err;
	}
	if (err_str)
		free(err_str);
	return (void*)&err_no_mem;
}

static BOOL is_class_registered() {
	WNDCLASSEX wcx;
	if (GetClassInfoEx(instance, CLASS_NAME, &wcx))
		return TRUE;
	return FALSE;
}

#include "win32_debug.h"
#include "win32_keys.h"
#include "win32_init.h"
#include "win32_window.h"

void g2d_error(void *const err, int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
	error_t *const error = (error_t*)err;
	err_num[0] = error->err_num;
	err_win32[0] = error->err_win32;
	err_str[0] = error->err_str;
}

void g2d_error_free(void *const err) {
	error_t *const err_t = (error_t*)err;
	if (err_t[0].err_str) {
		free(err_t[0].err_str);
		err_t[0].err_str = NULL;
	}
	if (err_t != &err_no_mem)
		free(err);
}

void *g2d_string_new(void **const str, void *const go_cstr) {
	void *err = NULL;
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
			MultiByteToWideChar(CP_UTF8, 0, (const char*)go_cstr, length, str_new, length);
	#else
	str_new = (LPTSTR)malloc(sizeof(char) * (length + 1));
	if (str_new) {
		if (length > 0)
			memcpy(str_new, go_cstr, length);
	#endif
		str_new[length] = 0;
	}
	else
		err = (void*)&err_no_mem;
	str[0] = (void*)str_new;
	return err;
}

void g2d_string_free(void *const str) {
	if (str)
		free(str);
}

void *g2d_process_events() {
	if (active_windows > 0) {
		MSG msg;
		while (err_static == NULL && GetMessage(&msg, NULL, 0, 0) > 0) {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
	return (void*)err_static;
}

void g2d_err_static_set(const int go_obj) {
	err_static = error_new(100, (DWORD) go_obj, NULL);
}

void *g2d_post_close(void *const data) {
	if (!PostMessage(((window_data_t*)data)[0].wnd.hndl, WM_CLOSE, 0, 0))
		return error_new(66, 0, NULL);
	return NULL;
}

void *g2d_post_update(void *const data) {
	if (!PostMessage(((window_data_t*)data)[0].wnd.hndl, WM_APP, MSG_UPDATE, 0))
		return error_new(67, 0, NULL);
	return NULL;
}

void *g2d_post_props(void *const data) {
	if (!PostMessage(((window_data_t*)data)[0].wnd.hndl, WM_APP, MSG_PROPS, 0))
		return error_new(68, 0, NULL);
	return NULL;
}

void *g2d_post_err(void *const data) {
	if (!PostMessage(((window_data_t*)data)[0].wnd.hndl, WM_APP, MSG_ERROR, 0))
		return error_new(68, 0, NULL);
	return NULL;
}
*/

/* #if defined(G2D_WIN32) */
#endif
