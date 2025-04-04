/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

void g2d_main_loop() {
	MSG msg; BOOL ret_code; stop = FALSE;
	thread_id = GetCurrentThreadId();
	g2dMainLoopStarted();
	while (!stop && (ret_code = GetMessage(&msg, NULL, 0, 0)) > 0) {
		if (msg.message == WM_APP) {
			if (msg.wParam == g2d_REQUEST_EVENT)
				g2dProcessRequest();
			else if (msg.wParam == g2d_QUIT_EVENT)
				break;
		} else {
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
	}
}
