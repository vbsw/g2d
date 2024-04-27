/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#if defined(G2D_MODULES_WIN32)

#define OGFL_CDATA_NAME "vbsw.g2d.ogfl"
#define RECTS_CDATA_NAME "vbsw.g2d.rects"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>
#include "modules.h"

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* Exported functions from Go are:                          */
/* g2dProcessMessage                                        */

typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

typedef void* (ogfl_load_t) (void *data, const char *name, long long *err);
typedef struct { ogfl_load_t *load_func; void *data; } oglf_t;

void g2d_mods_rects_init(const int pass, cdata_t *const cdata) {
	cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
	cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
	if (pass == 0) {
		void **const functions = (void**)malloc(sizeof(void*)*37);
		if (functions) {
			long long err; int i;
			oglf_t *const ogfl = (oglf_t*)get(cdata, OGFL_CDATA_NAME);
			ogfl_load_t *const og_load = ogfl[0].load_func;
			/* wgl functions */
			functions[0] = og_load(ogfl[0].data, "wglChoosePixelFormatARB", &err);
			set(cdata, "wglChoosePixelFormatARB", functions[0]);
			functions[1] = og_load(ogfl[0].data, "wglCreateContextAttribsARB", &err);
			set(cdata, "wglCreateContextAttribsARB", functions[1]);
			functions[2] = og_load(ogfl[0].data, "wglSwapIntervalEXT", &err);
			set(cdata, "wglSwapIntervalEXT", functions[2]);
			functions[3] = og_load(ogfl[0].data, "wglGetSwapIntervalEXT", &err);
			set(cdata, "wglGetSwapIntervalEXT", functions[3]);
			/* ogl functions */
			functions[4] = og_load(ogfl[0].data, "glCreateShader", &err);
			set(cdata, "glCreateShader", functions[4]);
			functions[5] = og_load(ogfl[0].data, "glShaderSource", &err);
			set(cdata, "glShaderSource", functions[5]);
			functions[6] = og_load(ogfl[0].data, "glCompileShader", &err);
			set(cdata, "glCompileShader", functions[6]);
			functions[7] = og_load(ogfl[0].data, "glGetShaderiv", &err);
			set(cdata, "glGetShaderiv", functions[7]);
			functions[8] = og_load(ogfl[0].data, "glGetShaderInfoLog", &err);
			set(cdata, "glGetShaderInfoLog", functions[8]);
			functions[9] = og_load(ogfl[0].data, "glCreateProgram", &err);
			set(cdata, "glCreateProgram", functions[9]);
			functions[10] = og_load(ogfl[0].data, "glAttachShader", &err);
			set(cdata, "glAttachShader", functions[10]);
			functions[11] = og_load(ogfl[0].data, "glLinkProgram", &err);
			set(cdata, "glLinkProgram", functions[11]);
			functions[12] = og_load(ogfl[0].data, "glValidateProgram", &err);
			set(cdata, "glValidateProgram", functions[12]);
			functions[13] = og_load(ogfl[0].data, "glGetProgramiv", &err);
			set(cdata, "glGetProgramiv", functions[13]);
			functions[14] = og_load(ogfl[0].data, "glGetProgramInfoLog", &err);
			set(cdata, "glGetProgramInfoLog", functions[14]);
			functions[15] = og_load(ogfl[0].data, "glGenBuffers", &err);
			set(cdata, "glGenBuffers", functions[15]);
			functions[16] = og_load(ogfl[0].data, "glGenVertexArrays", &err);
			set(cdata, "glGenVertexArrays", functions[16]);
			functions[17] = og_load(ogfl[0].data, "glGetAttribLocation", &err);
			set(cdata, "glGetAttribLocation", functions[17]);
			functions[18] = og_load(ogfl[0].data, "glBindVertexArray", &err);
			set(cdata, "glBindVertexArray", functions[18]);
			functions[19] = og_load(ogfl[0].data, "glEnableVertexAttribArray", &err);
			set(cdata, "glEnableVertexAttribArray", functions[19]);
			functions[20] = og_load(ogfl[0].data, "glVertexAttribPointer", &err);
			set(cdata, "glVertexAttribPointer", functions[20]);
			functions[21] = og_load(ogfl[0].data, "glBindBuffer", &err);
			set(cdata, "glBindBuffer", functions[21]);
			functions[22] = og_load(ogfl[0].data, "glBufferData", &err);
			set(cdata, "glBufferData", functions[22]);
			functions[23] = og_load(ogfl[0].data, "glBufferSubData", &err);
			set(cdata, "glBufferSubData", functions[23]);
			functions[24] = og_load(ogfl[0].data, "glGetVertexAttribPointerv", &err);
			set(cdata, "glGetVertexAttribPointerv", functions[24]);
			functions[25] = og_load(ogfl[0].data, "glUseProgram", &err);
			set(cdata, "glUseProgram", functions[25]);
			functions[26] = og_load(ogfl[0].data, "glDeleteVertexArrays", &err);
			set(cdata, "glDeleteVertexArrays", functions[26]);
			functions[27] = og_load(ogfl[0].data, "glDeleteBuffers", &err);
			set(cdata, "glDeleteBuffers", functions[27]);
			functions[28] = og_load(ogfl[0].data, "glDeleteProgram", &err);
			set(cdata, "glDeleteProgram", functions[28]);
			functions[29] = og_load(ogfl[0].data, "glDeleteShader", &err);
			set(cdata, "glDeleteShader", functions[29]);
			functions[30] = og_load(ogfl[0].data, "glGetUniformLocation", &err);
			set(cdata, "glGetUniformLocation", functions[30]);
			functions[31] = og_load(ogfl[0].data, "glUniformMatrix3fv", &err);
			set(cdata, "glUniformMatrix3fv", functions[31]);
			functions[32] = og_load(ogfl[0].data, "glUniform1fv", &err);
			set(cdata, "glUniform1fv", functions[32]);
			functions[33] = og_load(ogfl[0].data, "glUniformMatrix4fv", &err);
			set(cdata, "glUniformMatrix4fv", functions[33]);
			functions[34] = og_load(ogfl[0].data, "glUniformMatrix2x3fv", &err);
			set(cdata, "glUniformMatrix2x3fv", functions[34]);
			functions[35] = og_load(ogfl[0].data, "glGenerateMipmap", &err);
			set(cdata, "glGenerateMipmap", functions[35]);
			functions[36] = og_load(ogfl[0].data, "glActiveTexture", &err);
			set(cdata, "glActiveTexture", functions[36]);
			for (i = 0; i < 2000; i++) {
				/* wgl functions */
				if (functions[0])
					functions[0] = get(cdata, "wglChoosePixelFormatARB");
				if (functions[1])
					functions[1] = get(cdata, "wglCreateContextAttribsARB");
				if (functions[2])
					functions[2] = get(cdata, "wglSwapIntervalEXT");
				if (functions[3])
					functions[3] = get(cdata, "wglGetSwapIntervalEXT");
				/* ogl functions */
				if (functions[4])
					functions[4] = get(cdata, "glCreateShader");
				if (functions[5])
					functions[5] = get(cdata, "glShaderSource");
				if (functions[6])
					functions[6] = get(cdata, "glCompileShader");
				if (functions[7])
					functions[7] = get(cdata, "glGetShaderiv");
				if (functions[8])
					functions[8] = get(cdata, "glGetShaderInfoLog");
				if (functions[9])
					functions[9] = get(cdata, "glCreateProgram");
				if (functions[10])
					functions[10] = get(cdata, "glAttachShader");
				if (functions[11])
					functions[11] = get(cdata, "glLinkProgram");
				if (functions[12])
					functions[12] = get(cdata, "glValidateProgram");
				if (functions[13])
					functions[13] = get(cdata, "glGetProgramiv");
				if (functions[14])
					functions[14] = get(cdata, "glGetProgramInfoLog");
				if (functions[15])
					functions[15] = get(cdata, "glGenBuffers");
				if (functions[16])
					functions[16] = get(cdata, "glGenVertexArrays");
				if (functions[17])
					functions[17] = get(cdata, "glGetAttribLocation");
				if (functions[18])
					functions[18] = get(cdata, "glBindVertexArray");
				if (functions[19])
					functions[19] = get(cdata, "glEnableVertexAttribArray");
				if (functions[20])
					functions[20] = get(cdata, "glVertexAttribPointer");
				if (functions[21])
					functions[21] = get(cdata, "glBindBuffer");
				if (functions[22])
					functions[22] = get(cdata, "glBufferData");
				if (functions[23])
					functions[23] = get(cdata, "glBufferSubData");
				if (functions[24])
					functions[24] = get(cdata, "glGetVertexAttribPointerv");
				if (functions[25])
					functions[25] = get(cdata, "glUseProgram");
				if (functions[26])
					functions[26] = get(cdata, "glDeleteVertexArrays");
				if (functions[27])
					functions[27] = get(cdata, "glDeleteBuffers");
				if (functions[28])
					functions[28] = get(cdata, "glDeleteProgram");
				if (functions[29])
					functions[29] = get(cdata, "glDeleteShader");
				if (functions[30])
					functions[30] = get(cdata, "glGetUniformLocation");
				if (functions[31])
					functions[31] = get(cdata, "glUniformMatrix3fv");
				if (functions[32])
					functions[32] = get(cdata, "glUniform1fv");
				if (functions[33])
					functions[33] = get(cdata, "glUniformMatrix4fv");
				if (functions[34])
					functions[34] = get(cdata, "glUniformMatrix2x3fv");
				if (functions[35])
					functions[35] = get(cdata, "glGenerateMipmap");
				if (functions[36])
					functions[36] = get(cdata, "glActiveTexture");
			}
			set(cdata, RECTS_CDATA_NAME, (void*)functions);
		} else {
			cdata[0].err1 = 20;
		}
	} else if (pass == 1) {
		void *const functions = get(cdata, RECTS_CDATA_NAME);
		if (functions)
			free(functions);
		else
			cdata[0].err1 = 21;
	} else if (pass < 0) {
		void *const functions = get(cdata, RECTS_CDATA_NAME);
		if (functions)
			free(functions);
	}
}

/*
void g2d_mods_rects_init(void **const data, int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
}
*/

/* #if defined(G2D_MODULES_WIN32) */
#endif
