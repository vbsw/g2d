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

void g2d_gfx_draw(void *const data, const int w, const int h, const int r, const int g, const int b, const int i, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wnd_data[0].gfx.w != w || wnd_data[0].gfx.g != h) {
		wnd_data[0].gfx.w = w; wnd_data[0].gfx.h = h;
		glViewport((WORD)0, (WORD)0, (WORD)w, (WORD)h);
	}
	if (wnd_data[0].gfx.r != r || wnd_data[0].gfx.g != g || wnd_data[0].gfx.b != b) {
		wnd_data[0].gfx.r = r; wnd_data[0].gfx.g = g; wnd_data[0].gfx.b = b;
		glClearColor((GLclampf)r, (GLclampf)g, (GLclampf)b, 0.0);
	}
	if (wnd_data[0].gfx.i != i) {
		wnd_data[0].gfx.i = i;
		wglSwapIntervalEXT(i);
	}
	glClear(GL_COLOR_BUFFER_BIT);
	if (!SwapBuffers(wnd_data[0].wnd.dc)) {
		err1[0] = 220; err2[0] = (long long)GetLastError();
	}
}
