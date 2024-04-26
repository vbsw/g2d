/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#include "analytics.h"

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <gl/GL.h>

/* from glcorearb.h */
#define GL_MAX_TEXTURE_IMAGE_UNITS        0x8872

#define NLTX_ID "vbsw.g2d.analytics"

/* from github.com/vbsw/golib/cdata/cdata.c */
typedef void (*cdata_set_func_t)(cdata_t *cdata, void *data, const char *id);
typedef void* (*cdata_get_func_t)(cdata_t *cdata, const char *id);

/* for internal usage */
typedef struct { int max_tex_size; int max_tex_units; } nltx_t;

void vbsw_nltx_init(const int pass, cdata_t *const cdata) {
	if (pass == 0) {
		nltx_t *const nltx = (nltx_t*)malloc(sizeof(nltx_t));
		if (nltx) {
			cdata_set_func_t const set = (cdata_set_func_t)cdata[0].set_func;
			GLint i_val = 0;
			glGetIntegerv(GL_MAX_TEXTURE_SIZE, &i_val); nltx[0].max_tex_size = (int)i_val;
			glGetIntegerv(GL_MAX_TEXTURE_IMAGE_UNITS, &i_val); nltx[0].max_tex_units = (int)i_val;
			set(cdata, (void*)nltx, NLTX_ID);
		} else {
			cdata[0].err1 = 20;
		}
	} else if (pass < 0) {
		cdata_get_func_t const get = (cdata_get_func_t)cdata[0].get_func;
		nltx_t *const nltx = (nltx_t*)get(cdata, NLTX_ID);
		if (nltx)
			free(nltx);
	}
}

void vbsw_nltx_result_and_free(void *const data, int *const mts, int *const mtu) {
	if (data) {
		nltx_t *const nltx = (nltx_t*)data;
		mts[0] = nltx[0].max_tex_size;
		mtu[0] = nltx[0].max_tex_units;
		free(data);
	}
}
