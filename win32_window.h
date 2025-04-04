/*
 *       Copyright 2024, 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static LPTSTR to_tstr(void *const go_cstr, const size_t length, long long *err1) {
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
	return str_new;
}

static void window_metrics(window_data_t *const wnd_data, int *const x, int *const y, int *const w, int *const h) {
	RECT rect = { wnd_data[0].client.x, wnd_data[0].client.y, wnd_data[0].client.x + wnd_data[0].client.width, wnd_data[0].client.y + wnd_data[0].client.height };
	AdjustWindowRect(&rect, wnd_data[0].config.style, FALSE);
	x[0] = rect.left; y[0] = rect.top; w[0] = rect.right - rect.left; h[0] = rect.bottom - rect.top;
}

static void monitor_metrics(HMONITOR const monitor, int *const x, int *const y, int *const w, int *const h) {
	MONITORINFO mi = { sizeof(mi) }; GetMonitorInfo(monitor, &mi);
	x[0] = mi.rcMonitor.left; y[0] = mi.rcMonitor.top;
	w[0] = mi.rcMonitor.right - mi.rcMonitor.left; h[0] = mi.rcMonitor.bottom - mi.rcMonitor.top;
}

static void style_update(window_data_t *const wnd_data) {
	if (wnd_data[0].config.borderless)
		if (wnd_data[0].config.resizable)
			wnd_data[0].config.style = WS_POPUP;
		else
			wnd_data[0].config.style = WS_POPUP;
	else
		if (wnd_data[0].config.resizable)
			wnd_data[0].config.style = WS_OVERLAPPEDWINDOW;
		else
			wnd_data[0].config.style = WS_OVERLAPPEDWINDOW & ~WS_THICKFRAME & ~WS_MAXIMIZEBOX;
}

static LRESULT CALLBACK windowProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam) {
	LRESULT result = 0;
	if (message == WM_NCCREATE) {
		window_data_t *const wnd_data = (window_data_t*)(((CREATESTRUCT*)lParam)->lpCreateParams);
		if (wnd_data)
			SetWindowLongPtr(hWnd, GWLP_USERDATA, (LONG_PTR)wnd_data);
		result = DefWindowProc(hWnd, message, wParam, lParam);
	} else {
		window_data_t *const wnd_data = (window_data_t*)GetWindowLongPtr(hWnd, GWLP_USERDATA);
		if (wnd_data) {
			if (!wnd_data[0].state.minimized) {
				switch (message) {
/*
				case WM_MOVE:
					client_update(wnd_data);
					result = DefWindowProc(hWnd, message, wParam, lParam);
					g2dWindowMove(wnd_data[0].cb_id);
					break;
				case WM_SIZE:
					// avoid resize on show
					if (wnd_data[0].state.shown != 0) {
						client_update(wnd_data);
						g2dWindowResize(wnd_data[0].cb_id);
					} else {
						wnd_data[0].state.shown = 1;
					}
					result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
				case WM_CLOSE:
					g2dClose(wnd_data[0].cb_id);
					break;
*/
/*
				case WM_KEYDOWN:
					if (!key_down_process(wnd_data, message, wParam, lParam))
						result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
				case WM_KEYUP:
					if (!key_up_process(wnd_data, message, wParam, lParam))
						result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
*/
/*
				case WM_SYSCOMMAND:
					if (wParam == SC_MINIMIZE) {
						wnd_data[0].state.minimized = 1;
						g2dWindowMinimize(wnd_data[0].cb_id);
					}
					result = DefWindowProc(hWnd, message, wParam, lParam);
					break;
				case WM_MOUSEMOVE:
					wnd_data[0].mouse.x = ((int)(short)LOWORD(lParam));
					wnd_data[0].mouse.y = ((int)(short)HIWORD(lParam));
					g2dMouseMove(wnd_data[0].cb_id);
					result = DefWindowProc(hWnd, message, wParam, lParam);
*/


/*
					if (state.dragging_cust && !state.maximized) {
						move_window(client.x + (int)(short)LOWORD(lParam) - mouse.x, client.y + (int)(short)HIWORD(lParam) - mouse.y, client.width, client.height);
					} else {
						mouse.x = ((int)(short)LOWORD(lParam));
						mouse.y = ((int)(short)HIWORD(lParam));
						result = DefWindowProc(hWnd, message, wParam, lParam);
					}
					if (config.locked && !state.locked && state.focus)
						update_clip_cursor();
*/

/*
					break;
				case WM_LBUTTONDOWN:
					button_down(wnd_data, 0, 0);
					break;
				case WM_LBUTTONUP:
					button_up(wnd_data, 0);
					break;
				case WM_LBUTTONDBLCLK:
					button_down(wnd_data, 0, 1);
					break;
				case WM_RBUTTONDOWN:
					button_down(wnd_data, 1, 0);
					break;
				case WM_RBUTTONUP:
					button_up(wnd_data, 1);
					break;
				case WM_RBUTTONDBLCLK:
					button_down(wnd_data, 1, 1);
					break;
				case WM_MBUTTONDOWN:
					button_down(wnd_data, 2, 0);
					break;
				case WM_MBUTTONUP:
					button_up(wnd_data, 2);
					break;
				case WM_MBUTTONDBLCLK:
					button_down(wnd_data, 2, 1);
					break;
				case WM_MOUSEWHEEL:
					g2dWheel(wnd_data[0].cb_id, (float)GET_WHEEL_DELTA_WPARAM(wParam) / (float)WHEEL_DELTA);
					break;
				case WM_XBUTTONDOWN:
					if (HIWORD(wParam) == XBUTTON1)
						button_down(wnd_data, 3, 0);
					else if (HIWORD(wParam) == XBUTTON2)
						button_down(wnd_data, 4, 0);
					break;
				case WM_XBUTTONUP:
					if (HIWORD(wParam) == XBUTTON1)
						button_up(wnd_data, 3);
					else if (HIWORD(wParam) == XBUTTON2)
						button_up(wnd_data, 4);
					break;
				case WM_XBUTTONDBLCLK:
					if (HIWORD(wParam) == XBUTTON1)
						button_down(wnd_data, 3, 1);
					else if (HIWORD(wParam) == XBUTTON2)
						button_down(wnd_data, 4, 1);
					break;
*/
				default:
					result = DefWindowProc(hWnd, message, wParam, lParam);
				}
			} else {
				result = DefWindowProc(hWnd, message, wParam, lParam);
				if (message == WM_SETFOCUS) {
					// restore from minimized and avoid move/resize events
/*
					if (wnd_data[0].state.minimized) {
						wnd_data[0].state.minimized = 0;
						g2dWindowRestore(wnd_data[0].cb_id);
					}
*/
				}
			}
		} else {
			result = DefWindowProc(hWnd, message, wParam, lParam);
		}
	}
	return result;
}

void g2d_window_create(void **const data, const int cb_id, const int x, const int y, const int w, const int h, const int wn, const int hn, const int wx, const int hx,
	const int b, const int d, const int r, const int f, const int l, const int c, void *const t, const size_t ts, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)malloc(sizeof(window_data_t));
	if (wnd_data) {
		LPCTSTR const title = to_tstr(t, ts, err1);
		if (err1[0] == 0) {
			ZeroMemory(wnd_data, sizeof(window_data_t));
			wnd_data[0].cb_id = cb_id;
			wnd_data[0].client.x = x;
			wnd_data[0].client.y = y;
			wnd_data[0].client.width = w;
			wnd_data[0].client.height = h;
			wnd_data[0].config.width_min = wn;
			wnd_data[0].config.height_min = hn;
			wnd_data[0].config.width_max = wx;
			wnd_data[0].config.height_max = hx;
			wnd_data[0].config.borderless = b;
			wnd_data[0].config.dragable = d;
			wnd_data[0].config.fullscreen = f;
			wnd_data[0].config.resizable = r;
			wnd_data[0].config.locked = l;
			style_update(wnd_data);
			if (c) {
				int wx, wy, ww, wh, mx, my, mw, mh;
				window_metrics(wnd_data, &wx, &wy, &ww, &wh);
				monitor_metrics(MonitorFromWindow(NULL, MONITOR_DEFAULTTOPRIMARY), &mx, &my, &mw, &mh);
				wnd_data[0].client.x = mx + (mw - ww) / 2 + (wnd_data[0].client.x - wx);
				wnd_data[0].client.y = my + (mh - wh) / 2 + (wnd_data[0].client.y - wy);
			}
			if (windows_count == 0) {
				WNDCLASSEX cls;
				ZeroMemory(&cls, sizeof(WNDCLASSEX));
				cls.cbSize = sizeof(WNDCLASSEX);
				cls.style = CS_OWNDC | CS_HREDRAW | CS_VREDRAW | CS_DBLCLKS;
				cls.lpfnWndProc = windowProc;
				cls.hInstance = instance;
				cls.hIcon = LoadIcon(NULL, IDI_WINLOGO);
				cls.hCursor = LoadCursor(NULL, IDC_ARROW);
				cls.lpszClassName = class_name;
				if (RegisterClassEx(&cls) != INVALID_ATOM) {
					windows_count++;
				} else {
					err1[0] = 13; err2[0] = (long long)GetLastError(); free(wnd_data);
				}
			} else {
				windows_count++;
			}
			if (err1[0] == 0) {
				int x, y, w, h; window_metrics(wnd_data, &x, &y, &w, &h);
				const DWORD style = wnd_data[0].config.style;
				wnd_data[0].wnd.hndl = CreateWindow(class_name, title, style, x, y, w, h, NULL, NULL, instance, (LPVOID)wnd_data);
				if (wnd_data[0].wnd.hndl) {
					wnd_data[0].wnd.dc = GetDC(wnd_data[0].wnd.hndl);
					if (wnd_data[0].wnd.dc) {
						int pixelFormat;
						BOOL status = FALSE;
						UINT numFormats = 0;
						const int pixelAttribs[] = {
							WGL_DRAW_TO_WINDOW_ARB, GL_TRUE,
							WGL_SUPPORT_OPENGL_ARB, GL_TRUE,
							WGL_DOUBLE_BUFFER_ARB, GL_TRUE,
							/* WGL_SWAP_COPY_ARB might have update problems in fullscreen */
							/* WGL_SWAP_EXCHANGE_ARB might have problems with start menu in fullscreen */
							WGL_SWAP_METHOD_ARB, WGL_SWAP_EXCHANGE_ARB,
							WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
							WGL_ACCELERATION_ARB, WGL_FULL_ACCELERATION_ARB,
							WGL_COLOR_BITS_ARB, 32,
							WGL_ALPHA_BITS_ARB, 8,
							WGL_DEPTH_BITS_ARB, 24,
							/* anti aliasing */
							//WGL_SAMPLE_BUFFERS_ARB, 1,
							//WGL_SAMPLES_ARB, 4,
							0
						};
						const int contextAttributes[] = {
							WGL_CONTEXT_MAJOR_VERSION_ARB, 3,
							WGL_CONTEXT_MINOR_VERSION_ARB, 0,
							WGL_CONTEXT_PROFILE_MASK_ARB, WGL_CONTEXT_CORE_PROFILE_BIT_ARB,
							0
						};
						status = wglChoosePixelFormatARB(wnd_data[0].wnd.dc, pixelAttribs, NULL, 1, &pixelFormat, &numFormats);
						if (status && numFormats) {
							PIXELFORMATDESCRIPTOR pfd;
							ZeroMemory(&pfd, sizeof(PIXELFORMATDESCRIPTOR));
							DescribePixelFormat(wnd_data[0].wnd.dc, pixelFormat, sizeof(PIXELFORMATDESCRIPTOR), &pfd);
							if (SetPixelFormat(wnd_data[0].wnd.dc, pixelFormat, &pfd)) {
								wnd_data[0].wnd.rc = wglCreateContextAttribsARB(wnd_data[0].wnd.dc, 0, contextAttributes);
								if (wnd_data[0].wnd.rc) {
									data[0] = (void*)wnd_data;
								} else {
									err1[0] = 18; err2[0] = (long long)GetLastError(); windows_count--;
									ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.dc); DestroyWindow(wnd_data[0].wnd.hndl);
									free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
								}
							} else {
								err1[0] = 17; err2[0] = (long long)GetLastError(); windows_count--;
								ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.dc); DestroyWindow(wnd_data[0].wnd.hndl);
								free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
							}
						} else {
							err1[0] = 16; err2[0] = (long long)GetLastError(); windows_count--;
							ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.dc); DestroyWindow(wnd_data[0].wnd.hndl);
							free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
						}
					} else {
						err1[0] = 15; windows_count--; DestroyWindow(wnd_data[0].wnd.hndl);
						free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
					}
				} else {
					err1[0] = 14; err2[0] = (long long)GetLastError(); windows_count--;
					free(wnd_data); if (windows_count <= 0) { UnregisterClass(class_name, instance); }
				}
			}
			free((void*)title);
		} else {
			err1[0] = 122;
		}
	} else {
		err1[0] = 121;
	}
}

void g2d_window_show(void *const data, long long *const err1, long long *const err2) {
	if (data) {
		window_data_t *const wnd_data = (window_data_t*)data;
		ShowWindow(wnd_data[0].wnd.hndl, SW_SHOWDEFAULT);
	}
}
