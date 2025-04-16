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

#define G2D_RESIZE_BORDER 4

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
#define GL_CLAMP_TO_EDGE                  0x812F
#define GL_MAX_TEXTURE_IMAGE_UNITS        0x8872

/* from wglext.h */
typedef BOOL(WINAPI * PFNWGLCHOOSEPIXELFORMATARBPROC) (HDC hdc, const int *piAttribIList, const FLOAT *pfAttribFList, UINT nMaxFormats, int *piFormats, UINT *nNumFormats);
typedef HGLRC(WINAPI * PFNWGLCREATECONTEXTATTRIBSARBPROC) (HDC hDC, HGLRC hShareContext, const int *attribList);
typedef const char *(WINAPI * PFNWGLGETEXTENSIONSSTRINGARBPROC) (HDC hdc);
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
typedef void (APIENTRY *PFNGLUNIFORM1IPROC) (GLint location, GLint v0);
typedef void (APIENTRY *PFNGLUNIFORMMATRIX4FVPROC) (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
typedef void (APIENTRY *PFNGLUNIFORMMATRIX3FVPROC) (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
typedef void (APIENTRY *PFNGLUNIFORMMATRIX2X3FVPROC) (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
typedef void (APIENTRY *PFNGLACTIVETEXTUREPROC) (GLenum texture);
typedef void (APIENTRY *PFNGLGENERATEMIPMAPPROC) (GLenum target);

static PFNWGLCHOOSEPIXELFORMATARBPROC    wglChoosePixelFormatARB    = NULL;
static PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
static PFNWGLGETEXTENSIONSSTRINGARBPROC  wglGetExtensionsStringARB  = NULL;
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
static PFNGLUNIFORM1IPROC                glUniform1i                = NULL;
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
	struct { int dragging, minimized, maximized, resizing, focus, shown; } state;
	unsigned int key_repeated[255];
	int cb_id;
	struct { int r, g, b, w, h, i; GLfloat unif_data[16*3]; } gfx;
	struct { GLuint prog_ref, vao_ref, vbo_ref, ebo_ref, buf_max_len; GLint att_lc[4], unif_lc[17]; GLfloat *buffer; } rects;
} window_data_t;

typedef void (gfx_draw_t)(void *data, float *rects, int total, long long *err1);

static const WPARAM g2d_REQUEST_EVENT = (WPARAM)"g2dc";
static const WPARAM g2d_QUIT_EVENT    = (WPARAM)"g2dq";
static LPCTSTR const class_name       = TEXT("g2d");
static LPCTSTR const class_name_dummy = TEXT("g2d_dummy");

static const GLfloat default_projection_mat[16] = { 2.0f / 1.0f, 0.0f, 0.0f, 0.0f, 0.0f, -2.0f / 1.0f, 0.0f, 0.0f, 0.0f, 0.0f, -1.0f, 0.0f, -1.0f, 1.0f, 0.0f, 1.0f };

static HINSTANCE instance = NULL;
static BOOL initialized   = FALSE;
static int windows_count  = 0;
static DWORD thread_id    = 0;
static BOOL stop          = FALSE;


static LPCSTR const vs_rect_str = "#version 130\n\
in vec4 in0; \
in vec4 in1; \
in vec4 in2; \
in vec4 in3; \
out vec4 fragementColor; \
out vec3 texCoord; \
uniform float[48] unif; \
void main() { \
  int texIdx = int(in2[0]); \
  gl_Position = mat4(unif[0], unif[1], unif[2], unif[3], unif[4], unif[5], unif[6], unif[7], unif[8], unif[9], unif[10], unif[11], unif[12], unif[13], unif[14], unif[15]) * vec4(in0[0], in0[1], 1.0, 1.0); \
  fragementColor = in1; \
  if (texIdx >= 0) { \
    int offset = 16 + texIdx*2; \
    float texWidth = unif[offset + 0]; \
    float texHeight = unif[offset + 1]; \
    texCoord = vec3(in2[0], in3[0]/texWidth, in3[1]/texHeight); \
  } else { \
    texCoord = vec3(-1.0, 0.0, 0.0); \
  } \
}";
static LPCSTR const fs_rect_str = "#version 130\n\
in vec4 fragementColor; \
in vec3 texCoord; \
out vec4 color; \
uniform sampler2D tex00; uniform sampler2D tex01; uniform sampler2D tex02; uniform sampler2D tex03; \
uniform sampler2D tex04; uniform sampler2D tex05; uniform sampler2D tex06; uniform sampler2D tex07; \
uniform sampler2D tex08; uniform sampler2D tex09; uniform sampler2D tex10; uniform sampler2D tex11; \
uniform sampler2D tex12; uniform sampler2D tex13; uniform sampler2D tex14; uniform sampler2D tex15; \
void main() { \
  int texIdx = int(texCoord[0]); \
  if (texIdx >= 0) { \
    switch (texIdx) { \
      case 0: color = texture(tex00, vec2(texCoord[1], texCoord[2])); break; \
      case 1: color = texture(tex01, vec2(texCoord[1], texCoord[2])); break; \
      case 2: color = texture(tex02, vec2(texCoord[1], texCoord[2])); break; \
      case 3: color = texture(tex03, vec2(texCoord[1], texCoord[2])); break; \
      case 4: color = texture(tex04, vec2(texCoord[1], texCoord[2])); break; \
      case 5: color = texture(tex05, vec2(texCoord[1], texCoord[2])); break; \
      case 6: color = texture(tex06, vec2(texCoord[1], texCoord[2])); break; \
      case 7: color = texture(tex07, vec2(texCoord[1], texCoord[2])); break; \
      case 8: color = texture(tex08, vec2(texCoord[1], texCoord[2])); break; \
      case 9: color = texture(tex09, vec2(texCoord[1], texCoord[2])); break; \
      case 10: color = texture(tex10, vec2(texCoord[1], texCoord[2])); break; \
      case 11: color = texture(tex11, vec2(texCoord[1], texCoord[2])); break; \
      case 12: color = texture(tex12, vec2(texCoord[1], texCoord[2])); break; \
      case 13: color = texture(tex13, vec2(texCoord[1], texCoord[2])); break; \
      case 14: color = texture(tex14, vec2(texCoord[1], texCoord[2])); break; \
      case 15: color = texture(tex15, vec2(texCoord[1], texCoord[2])); break; \
    } \
  } else { \
    color = fragementColor; \
  } \
}";


void g2d_free(void *const data) {
	free(data);
}

#include "win32_debug.h"
#include "win32_keys.h"
#include "win32_init.h"
#include "win32_main_loop.h"
#include "win32_graphics.h"
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

/* #if defined(G2D_WIN32) */
#endif
