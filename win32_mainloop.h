/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

void g2d_mainloop_process_messages() {
	MSG msg; BOOL ret_code; int skip = 0; stop = FALSE;
	thread_id = GetCurrentThreadId();
	g2dMainLoopInit();
	while (!stop && (ret_code = GetMessage(&msg, NULL, 0, 0)) > 0) {
		if (msg.message == WM_APP) {
			if (msg.wParam == g2d_CUSTOM_EVENT)
				if (skip == 0)
					g2dMainLoopProcessCustomEvents(&skip);
				else
					skip--;
			else if (msg.wParam == g2d_QUIT_EVENT)
				break;
		} else {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
}

void g2d_mainloop_post_custom(long long *const err2) {
	if (!PostThreadMessage(thread_id, WM_APP, g2d_CUSTOM_EVENT, 0))
		err2[0] = (long long)GetLastError();
}

void g2d_mainloop_post_quit(long long *const err2) {
	stop = TRUE;
	if (!PostThreadMessage(thread_id, WM_APP, g2d_QUIT_EVENT, 0))
		err2[0] = (long long)GetLastError();
}

void g2d_mainloop_clean_up() {
	MSG msg;
	while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE));
}
