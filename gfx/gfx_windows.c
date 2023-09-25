/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#if defined(G2D_GFX_WIN32)

#define OGFL_CDATA_ID "vbsw.g2d.ogfl"
#define RECTS_CDATA_ID "vbsw.g2d.gfx.rects"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>
#include "gfx.h"

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* Exported functions from Go are:                          */
/* g2dProcessMessage                                        */

typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

typedef void* (ogfl_load_func_t) (void *obj, const char *name, long long *err);
typedef struct { ogfl_load_func_t *load_func; void *obj; } oglf_t;

void g2d_gfx_rects_init(const int pass, cdata_t *const cdata) {
	cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
	cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
	if (pass == 0) {
		void **const functions = (void**)malloc(sizeof(void*)*37);
		if (functions) {
			long long err; int i;
			oglf_t *const ogfl = (oglf_t*)get(cdata, OGFL_CDATA_ID);
			void *const og_obj = ogfl[0].obj;
			ogfl_load_func_t *const og_load = ogfl[0].load_func;
			/* wgl functions */
			functions[0] = og_load(og_obj, "wglChoosePixelFormatARB", &err);
			functions[1] = og_load(og_obj, "wglCreateContextAttribsARB", &err);
			functions[2] = og_load(og_obj, "wglSwapIntervalEXT", &err);
			functions[3] = og_load(og_obj, "wglGetSwapIntervalEXT", &err);
			/* ogl functions */
			functions[4] = og_load(og_obj, "glCreateShader", &err);
			functions[5] = og_load(og_obj, "glShaderSource", &err);
			functions[6] = og_load(og_obj, "glCompileShader", &err);
			functions[7] = og_load(og_obj, "glGetShaderiv", &err);
			functions[8] = og_load(og_obj, "glGetShaderInfoLog", &err);
			functions[9] = og_load(og_obj, "glCreateProgram", &err);
			functions[10] = og_load(og_obj, "glAttachShader", &err);
			functions[11] = og_load(og_obj, "glLinkProgram", &err);
			functions[12] = og_load(og_obj, "glValidateProgram", &err);
			functions[13] = og_load(og_obj, "glGetProgramiv", &err);
			functions[14] = og_load(og_obj, "glGetProgramInfoLog", &err);
			functions[15] = og_load(og_obj, "glGenBuffers", &err);
			functions[16] = og_load(og_obj, "glGenVertexArrays", &err);
			functions[17] = og_load(og_obj, "glGetAttribLocation", &err);
			functions[18] = og_load(og_obj, "glBindVertexArray", &err);
			functions[19] = og_load(og_obj, "glEnableVertexAttribArray", &err);
			functions[20] = og_load(og_obj, "glVertexAttribPointer", &err);
			functions[21] = og_load(og_obj, "glBindBuffer", &err);
			functions[22] = og_load(og_obj, "glBufferData", &err);
			functions[23] = og_load(og_obj, "glBufferSubData", &err);
			functions[24] = og_load(og_obj, "glGetVertexAttribPointerv", &err);
			functions[25] = og_load(og_obj, "glUseProgram", &err);
			functions[26] = og_load(og_obj, "glDeleteVertexArrays", &err);
			functions[27] = og_load(og_obj, "glDeleteBuffers", &err);
			functions[28] = og_load(og_obj, "glDeleteProgram", &err);
			functions[29] = og_load(og_obj, "glDeleteShader", &err);
			functions[30] = og_load(og_obj, "glGetUniformLocation", &err);
			functions[31] = og_load(og_obj, "glUniformMatrix3fv", &err);
			functions[32] = og_load(og_obj, "glUniform1fv", &err);
			functions[33] = og_load(og_obj, "glUniformMatrix4fv", &err);
			functions[34] = og_load(og_obj, "glUniformMatrix2x3fv", &err);
			functions[35] = og_load(og_obj, "glGenerateMipmap", &err);
			functions[36] = og_load(og_obj, "glActiveTexture", &err);

			set(cdata, RECTS_CDATA_ID, (void*)functions);
		} else {
			cdata[0].err1 = 20;
		}
	} else if (pass == 1) {
		void *const functions = get(cdata, RECTS_CDATA_ID);
		if (functions)
			free(functions);
		else
			cdata[0].err1 = 21;
	} else if (pass < 0) {
		void *const functions = get(cdata, RECTS_CDATA_ID);
		if (functions)
			free(functions);
	}
}

/*
void g2d_gfx_rects_init(void **const data, int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
}
*/

/* #if defined(G2D_GFX_WIN32) */
#endif
