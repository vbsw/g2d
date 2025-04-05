/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

void g2d_gfx_make_current(void *const data, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (!wglMakeCurrent(wnd_data[0].wnd.dc, wnd_data[0].wnd.rc)) {
		err1[0] = 220; err2[0] = (long long)GetLastError();
	}
}

void g2d_gfx_release(void *const data, long long *const err1, long long *const err2) {
	if (!wglMakeCurrent(NULL, NULL)) {
		err1[0] = 220; err2[0] = (long long)GetLastError();
	}
}

void g2d_gfx_draw(void *const data, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (!SwapBuffers(wnd_data[0].wnd.dc)) {
		err1[0] = 220; err2[0] = (long long)GetLastError();
	}
}
