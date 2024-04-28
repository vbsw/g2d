/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#if defined(VBSW_G2D_WIN32)

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <stdio.h>
#include <gl/GL.h>
#include "g2d.h"

/* Go functions can not be passed to c directly.            */
/* They can only be called from c.                          */
/* This code is an indirection to call Go callbacks.        */
/* _cgo_export.h is generated automatically by cgo.         */
#include "_cgo_export.h"

/* Exported functions from Go are:                          */
/* g2dStartWindows                                          */
/* g2dProcessMessage                                        */
/* g2dResize                                                */
/* g2dKeyDown                                               */
/* g2dKeyUp                                                 */
/* g2dClose                                                 */

void g2d_free(void *const data) {
	free(data);
}

void g2d_to_tstr(void **const str, void *const go_cstr, const size_t length, long long *err1) {
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

/* #if defined(VBSW_G2D_WIN32) */
#endif
