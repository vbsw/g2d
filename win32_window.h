/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static void monitor_metrics(HMONITOR const monitor, int *const x, int *const y, int *const w, int *const h) {
	MONITORINFO mi = { sizeof(mi) };
	GetMonitorInfo(monitor, &mi);
	x[0] = mi.rcMonitor.left;
	y[0] = mi.rcMonitor.top;
	w[0] = mi.rcMonitor.right - mi.rcMonitor.left;
	h[0] = mi.rcMonitor.bottom - mi.rcMonitor.top;
}

static void window_metrics(window_data_t *const wnd_data, int *const x, int *const y, int *const w, int *const h) {
	RECT rect = { wnd_data[0].client.x, wnd_data[0].client.y, wnd_data[0].client.x + wnd_data[0].client.width, wnd_data[0].client.y + wnd_data[0].client.height };
	AdjustWindowRect(&rect, wnd_data[0].config.style, FALSE);
	x[0] = rect.left;
	y[0] = rect.top;
	w[0] = rect.right - rect.left;
	h[0] = rect.bottom - rect.top;
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

static LPCTSTR ensure_title(void *const title) {
	if (title)
		return (LPCTSTR)title;
	return TEXT("OpenGL");
}

static window_data_t *window_allocate(int *const err_num) {
	window_data_t *const wnd_data = (window_data_t*)malloc(sizeof(window_data_t));
	if (wnd_data)
		ZeroMemory(wnd_data, sizeof(window_data_t));
	else
		ERR_NEW1(-19)
	return wnd_data;
}

static void window_config(window_data_t *const wnd_data, const int x, const int y, const int w, const int h, const int wn, const int hn, const int wx, const int hx,
	const int b, const int d, const int r, const int f, const int l, const int c, int *const err_num) {
	if (err_num[0] == 0) {
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
	}
}

static void client_update(window_data_t *const wnd_data) {
	POINT point = { 0, 0 };
	RECT rect = {0, 0, 0, 0};
	ClientToScreen(wnd_data[0].wnd.hndl, &point);
	GetClientRect(wnd_data[0].wnd.hndl, &rect);
	wnd_data[0].client.x = point.x;
	wnd_data[0].client.y = point.y;
	wnd_data[0].client.width = (int)(rect.right - rect.left);
	wnd_data[0].client.height = (int)(rect.bottom - rect.top);
}

static BOOL dispatch_msg_go(WPARAM const wParam, const int go_obj_id) {
	if (wParam == MSG_SHOW) {
		g2dShow(go_obj_id);
		return TRUE;
	} else if (MSG_UPDATE) {
		g2dUpdate(go_obj_id);
		return TRUE;
	} else if (MSG_PROPS) {
		g2dProps(go_obj_id);
		return TRUE;
	} else if (MSG_ERROR) {
		g2dError(go_obj_id);
		return TRUE;
	}
	return FALSE;
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
			//print_message(message);
			switch (message) {
			case WM_MOVE:
				client_update(wnd_data);
				result = DefWindowProc(hWnd, message, wParam, lParam);
				//goOnMove();
				break;
			case WM_SIZE:
				client_update(wnd_data);
/*
				if (state.dragging) {
					if (!state.maximized)
						maximize_begin();
					else
						maximize_end();
				}
				goOnResize();
*/
				result = DefWindowProc(hWnd, message, wParam, lParam);
				break;
			case WM_CLOSE:
				g2dClose(wnd_data[0].go_obj_id);
				break;
			case WM_KEYDOWN:
				if (!key_down_process(wnd_data, message, wParam, lParam))
					result = DefWindowProc(hWnd, message, wParam, lParam);
				break;
			case WM_KEYUP:
				if (!key_up_process(wnd_data, message, wParam, lParam))
					result = DefWindowProc(hWnd, message, wParam, lParam);
				break;
			case WM_APP:
				if (!dispatch_msg_go(wParam, wnd_data[0].go_obj_id))
					result = DefWindowProc(hWnd, message, wParam, lParam);
			default:
				result = DefWindowProc(hWnd, message, wParam, lParam);
			}
		} else {
			result = DefWindowProc(hWnd, message, wParam, lParam);
		}
	}
	return result;

/*
	if (running && !state.minimized) {
		switch (message) {
		case WM_MOVE:
			update_client_props(client.width, client.height);
			result = DefWindowProc(hWnd, message, wParam, lParam);
			goOnMove();
			break;
		case WM_SIZE:
			update_client_props((int)LOWORD(lParam), (int)HIWORD(lParam));
			if (state.dragging) {
				if (!state.maximized)
					maximize_begin();
				else
					maximize_end();
			}
			goOnResize();
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_SETFOCUS:
			state.focus = 1;
			goOnFocusGain();
			break;
		case WM_KILLFOCUS:
			state.focus = 0;
			clear_keys();
			clear_clip_cursor();
			goOnFocusLoose();
			break;
		case WM_CLOSE:
			goOnClose();
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_SETCURSOR:
			if (LOWORD(lParam) == HTCLIENT) {
				SetCursor(mouse.cursor);
				result = TRUE;
			} else {
				result = DefWindowProc(hWnd, message, wParam, lParam);
			}
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_GETMINMAXINFO:
			get_window_min_max((LPMINMAXINFO)lParam);
			break;
		case WM_NCHITTEST:
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_NCMOUSEMOVE:
			drag_end();
			resize_end();
			break;
		case WM_NCLBUTTONDOWN:
			result = DefWindowProc(window.hndl, WM_NCHITTEST, wParam, lParam);
			if (result == HTCAPTION)
				drag_begin();
			else if (result == HTTOPLEFT || result == HTTOP || result == HTTOPRIGHT || result == HTRIGHT || result == HTBOTTOMRIGHT || result == HTBOTTOM || result == HTBOTTOMLEFT || result == HTLEFT)
				resize_begin();
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_NCLBUTTONUP:
			result = DefWindowProc(hWnd, message, wParam, lParam);
			drag_end();
			resize_end();
			break;
		case WM_NCLBUTTONDBLCLK:
			drag_end();
			resize_end();
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_KEYDOWN:
			if (!process_key_down(message, wParam, lParam))
				result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_KEYUP:
			if (!process_key_up(message, wParam, lParam))
				result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_SYSKEYDOWN:
			if (!process_key_down(message, wParam, lParam))
				result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_SYSKEYUP:
			if (!process_key_up(message, wParam, lParam))
				result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_SYSCOMMAND:
			if (wParam == SC_MINIMIZE) {
				state.minimized = 1;
				goOnMinimize();
			} else if (wParam == SC_MAXIMIZE) {
				maximize_begin();
			} else if (wParam == SC_RESTORE && state.maximized) {
				state.maximized = 0;
				goOnRestore();
			} else if (wParam == SC_MOVE) {
				drag_begin();
			} else if (wParam == SC_SIZE) {
				resize_begin();
			}
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_MOUSEMOVE:
			if (state.dragging_cust && !state.maximized) {
				move_window(client.x + (int)(short)LOWORD(lParam) - mouse.x, client.y + (int)(short)HIWORD(lParam) - mouse.y, client.width, client.height);
			} else {
				mouse.x = ((int)(short)LOWORD(lParam));
				mouse.y = ((int)(short)HIWORD(lParam));
				result = DefWindowProc(hWnd, message, wParam, lParam);
			}
			if (config.locked && !state.locked && state.focus)
				update_clip_cursor();
			break;
		case WM_LBUTTONDOWN:
			process_lb_down(message, wParam, lParam, 0);
			break;
		case WM_LBUTTONUP:
			ReleaseCapture();
			if (state.dragging_cust) {
				state.dragging_cust = 0;
				goOnDragCustEnd();
			}
			goOnButtonUp(1);
			break;
		case WM_LBUTTONDBLCLK:
			process_lb_down(message, wParam, lParam, 1);
			break;
		case WM_RBUTTONDOWN:
			goOnButtonDown(2, 0);
			break;
		case WM_RBUTTONUP:
			goOnButtonUp(2);
			break;
		case WM_RBUTTONDBLCLK:
			goOnButtonDown(2, 1);
			break;
		case WM_MBUTTONDOWN:
			goOnButtonDown(3, 0);
			break;
		case WM_MBUTTONUP:
			goOnButtonUp(3);
			break;
		case WM_MBUTTONDBLCLK:
			goOnButtonDown(3, 1);
			break;
		case WM_MOUSEWHEEL:
			goOnWheel((float)GET_WHEEL_DELTA_WPARAM(wParam) / (float)WHEEL_DELTA);
			break;
		case WM_XBUTTONDOWN:
			if (HIWORD(wParam) == XBUTTON1)
				goOnButtonDown(4, 0);
			else if (HIWORD(wParam) == XBUTTON2)
				goOnButtonDown(5, 0);
			break;
		case WM_XBUTTONUP:
			if (HIWORD(wParam) == XBUTTON1)
				goOnButtonUp(4);
			else if (HIWORD(wParam) == XBUTTON2)
				goOnButtonUp(5);
			break;
		case WM_XBUTTONDBLCLK:
			if (HIWORD(wParam) == XBUTTON1)
				goOnButtonDown(4, 1);
			else if (HIWORD(wParam) == XBUTTON2)
				goOnButtonDown(5, 1);
			break;
		case WM_ENTERMENULOOP:
			goOnMenuEnter();
			result = DefWindowProc(hWnd, message, wParam, lParam);
			break;
		case WM_EXITMENULOOP:
			result = DefWindowProc(hWnd, message, wParam, lParam);
			goOnMenuLeave();
			drag_end();
			resize_end();
			update_mouse_pos();
			break;
		case WM_EXITSIZEMOVE:
			result = DefWindowProc(hWnd, message, wParam, lParam);
			drag_end();
			resize_end();
			break;
		default:
			result = DefWindowProc(hWnd, message, wParam, lParam);
		}
	} else {
		if (message == WM_DESTROY)
			// stop event queue thread
			PostQuitMessage(0);
		result = DefWindowProc(hWnd, message, wParam, lParam);
		if (message == WM_SETFOCUS) {
			state.focus = 1;
			// restore from minimized and avoid move/resize events
			if (state.minimized) {
				state.minimized = 0;
				goOnRestore();
				drag_end();
				resize_end();
			}
		}
	}
*/
}

static void class_unregister(int *const err_num, g2d_ul_t *const err_win32) {
	if (active_windows < 2 && class_registered() && !UnregisterClass(class_name, instance) && err_num[0] == 0)
		ERR_NEW2(60, GetLastError())
}

static void class_register(int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0 && active_windows == 0) {
		WNDCLASSEX cls;
		ZeroMemory(&cls, sizeof(WNDCLASSEX));
		cls.cbSize = sizeof(WNDCLASSEX);
		cls.style = CS_OWNDC | CS_HREDRAW | CS_VREDRAW | CS_DBLCLKS;
		cls.lpfnWndProc = windowProc;
		cls.hInstance = instance;
		cls.hIcon = LoadIcon(NULL, IDI_WINLOGO);
		cls.hCursor = LoadCursor(NULL, IDC_ARROW);
		cls.lpszClassName = class_name;
		if (RegisterClassEx(&cls) == INVALID_ATOM)
			ERR_NEW2(50, GetLastError())
	}
}

static void window_create(window_data_t *const wnd_data, LPCTSTR const title, int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0) {
		int x, y, w, h; window_metrics(wnd_data, &x, &y, &w, &h);
		const DWORD style = wnd_data[0].config.style;
		wnd_data[0].wnd.hndl = CreateWindow(class_name, title, style, x, y, w, h, NULL, NULL, instance, (LPVOID)wnd_data);
		if (!wnd_data[0].wnd.hndl)
			ERR_NEW2(51, GetLastError())
	}
}

static void context_create(window_data_t *const wnd_data, int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0) {
		wnd_data[0].wnd.ctx.dc = GetDC(wnd_data[0].wnd.hndl);
		if (wnd_data[0].wnd.ctx.dc) {
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
				0
			};
			const int contextAttributes[] = {
				WGL_CONTEXT_MAJOR_VERSION_ARB, 3,
				WGL_CONTEXT_MINOR_VERSION_ARB, 0,
				WGL_CONTEXT_PROFILE_MASK_ARB, WGL_CONTEXT_CORE_PROFILE_BIT_ARB,
				0
			};
			status = wglChoosePixelFormatARB(wnd_data[0].wnd.ctx.dc, pixelAttribs, NULL, 1, &pixelFormat, &numFormats);
			if (status && numFormats) {
				PIXELFORMATDESCRIPTOR pfd;
				ZeroMemory(&pfd, sizeof(PIXELFORMATDESCRIPTOR));
				DescribePixelFormat(wnd_data[0].wnd.ctx.dc, pixelFormat, sizeof(PIXELFORMATDESCRIPTOR), &pfd);
				if (SetPixelFormat(wnd_data[0].wnd.ctx.dc, pixelFormat, &pfd)) {
					wnd_data[0].wnd.ctx.rc = wglCreateContextAttribsARB(wnd_data[0].wnd.ctx.dc, 0, contextAttributes);
					if (!wnd_data[0].wnd.ctx.rc)
						ERR_NEW2(55, GetLastError())
				} else {
					ERR_NEW2(54, GetLastError())
				}
			} else {
				ERR_NEW2(53, GetLastError())
			}
		} else {
			ERR_NEW1(52)
		}
	}
}

static void window_destroy(window_data_t *const wnd_data, int *const err_num, g2d_ul_t *const err_win32) {
	g2dDestroyBegin(wnd_data[0].go_obj_id);
	if (wnd_data[0].wnd.ctx.rc) {
		if (wnd_data[0].wnd.ctx.rc == wglGetCurrentContext() && !wglMakeCurrent(NULL, NULL) && err_num[0] == 0)
			ERR_NEW2(57, GetLastError())
		if (!wglDeleteContext(wnd_data[0].wnd.ctx.rc) && err_num[0] == 0) {
			ERR_NEW2(58, GetLastError())
		}
		wnd_data[0].wnd.ctx.rc = NULL;
	}
	if (wnd_data[0].wnd.ctx.dc) {
		ReleaseDC(wnd_data[0].wnd.hndl, wnd_data[0].wnd.ctx.dc);
		wnd_data[0].wnd.ctx.dc = NULL;
	}
	if (wnd_data[0].wnd.hndl) {
		if (!DestroyWindow(wnd_data[0].wnd.hndl) && err_num[0] == 0)
			ERR_NEW2(59, GetLastError())
		wnd_data[0].wnd.hndl = NULL;
	}
	class_unregister(err_num, err_win32);
	g2dDestroyEnd(wnd_data[0].go_obj_id);
	free(wnd_data);
}

static void client_props_update(window_data_t *const wnd_data) {
	RECT rect;
	POINT point = {0, 0};
	GetClientRect(wnd_data[0].wnd.hndl, &rect);
	ClientToScreen(wnd_data[0].wnd.hndl, &point);
	wnd_data[0].client.x = point.x;
	wnd_data[0].client.y = point.y;
	wnd_data[0].client.width = (int)(rect.right - rect.left);
	wnd_data[0].client.height = (int)(rect.bottom - rect.top);
}

static void cursor_clip_update(window_data_t *const wnd_data) {
	if (wnd_data[0].config.locked) {
		const RECT rect = { wnd_data[0].client.x, wnd_data[0].client.y, wnd_data[0].client.x + wnd_data[0].client.width, wnd_data[0].client.y + wnd_data[0].client.height };
		ClipCursor(&rect);
		wnd_data[0].state.locked = 1;
	}
}

void g2d_window_create(void **const data, const int go_obj, const int x, const int y, const int w, const int h, const int wn, const int hn, const int wx, const int hx,
	const int b, const int d, const int r, const int f, const int l, const int c, void *t, int *const err_num, g2d_ul_t *const err_win32) {
	if (initialized) {
		window_data_t *const wnd_data = window_allocate(err_num);
		window_config(wnd_data, x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, err_num);
		class_register(err_num, err_win32);
		window_create(wnd_data, ensure_title(t), err_num, err_win32);
		context_create(wnd_data, err_num, err_win32);
		if (err_num[0] == 0)
			data[0] = (void*)wnd_data;
		else if (wnd_data)
			window_destroy(wnd_data, err_num, err_win32);
	} else {
		ERR_NEW1(64)
	}
}

void g2d_window_show(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	if (initialized) {
		window_data_t *const wnd_data = (window_data_t*)data;
		ShowWindow(wnd_data[0].wnd.hndl, SW_SHOWDEFAULT);
		if (wnd_data[0].config.fullscreen)
			g2d_window_fullscreen_set(wnd_data, err_num, err_win32);
		if (err_num[0] == 0) {
			client_props_update(wnd_data);
			cursor_clip_update(wnd_data);
			if (PostMessage(wnd_data[0].wnd.hndl, WM_APP, MSG_SHOW, 0))
				active_windows++;
			else
				ERR_NEW1(65)
		}
	} else {
		ERR_NEW1(64)
	}
}

void g2d_window_destroy(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	if (data) {
		window_data_t *const wnd_data = (window_data_t*)data;
		window_destroy(wnd_data, err_num, err_win32);
		active_windows--;
	}
}

void g2d_window_props(void *const data, int *const x, int *const y, int *const w, int *const h, int *const wn, int *const hn,
	int *const wx, int *const hx, int *const b, int *const d, int *const r, int *const f, int *const l) {
	window_data_t *const wnd_data = (window_data_t*)data;
	*x = wnd_data[0].client.x;
	*y = wnd_data[0].client.y;
	*w = wnd_data[0].client.width;
	*h = wnd_data[0].client.height;
	*wn = wnd_data[0].config.width_min;
	*hn = wnd_data[0].config.height_min;
	*wx = wnd_data[0].config.width_max;
	*hx = wnd_data[0].config.height_max;
	*b = wnd_data[0].config.borderless;
	*d = wnd_data[0].config.dragable;
	*r = wnd_data[0].config.resizable;
	*f = wnd_data[0].config.fullscreen;
	*l = wnd_data[0].config.locked;
}

void g2d_window_fullscreen_set(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	window_data_t *const wnd_data = (window_data_t*)data;
	int mx, my, mw, mh; monitor_metrics(MonitorFromWindow(wnd_data[0].wnd.hndl, MONITOR_DEFAULTTONEAREST), &mx, &my, &mw, &mh);
	wnd_data[0].client_bak = wnd_data[0].client;
	const LONG_PTR style = SetWindowLongPtr(wnd_data[0].wnd.hndl, GWL_STYLE, 0);
	if (style && SetWindowPos(wnd_data[0].wnd.hndl, HWND_TOP, mx, my, mw, mh, SWP_NOOWNERZORDER | SWP_FRAMECHANGED | SWP_SHOWWINDOW))
		wnd_data[0].config.fullscreen = 1;
	else
		ERR_NEW2(70, GetLastError())
}

void g2d_client_pos_set(void *const data, const int x, const int y) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wnd_data[0].config.fullscreen) {
		wnd_data[0].client_bak.x = x;
		wnd_data[0].client_bak.y = y;
	} else {
		wnd_data[0].client.x = x;
		wnd_data[0].client.y = y;
	}
}

void g2d_client_size_set(void *const data, const int width, const int height) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wnd_data[0].config.fullscreen) {
		wnd_data[0].client_bak.width = width;
		wnd_data[0].client_bak.height = height;
	} else {
		wnd_data[0].client.width = width;
		wnd_data[0].client.height = height;
	}
}

void g2d_client_restore_bak(void *const data) {
	window_data_t *const wnd_data = (window_data_t*)data;
	wnd_data[0].client = wnd_data[0].client_bak;
}

/*
void *g2d_client_pos_apply(void *const data) {
	window_data_t *const wnd_data = (window_data_t*)data;
	int wx, wy, ww, wh; window_metrics(wnd_data, &wx, &wy, &ww, &wh);
	const LONG_PTR style = GetWindowLongPtr(wnd_data[0].wnd.hndl, GWL_STYLE);
	if (style == SetWindowLongPtr(wnd_data[0].wnd.hndl, GWL_STYLE, wnd_data[0].config.style)) {
		if (SetWindowPos(wnd_data[0].wnd.hndl, HWND_NOTOPMOST, wx, wy, ww, wh, SWP_NOOWNERZORDER | SWP_FRAMECHANGED | SWP_SHOWWINDOW)) {
			wnd_data[0].config.fullscreen = 0;
			cursor_clip_update(wnd_data);
			return NULL;
		}
	}
	return error_new(71, GetLastError(), NULL);
}

void *g2d_client_move(void *const data) {
	window_data_t *const wnd_data = (window_data_t*)data;
	int wx, wy, ww, wh; window_metrics(wnd_data, &wx, &wy, &ww, &wh);
	if (MoveWindow(wnd_data[0].wnd.hndl, wx, wy, ww, wh, FALSE))
		return NULL;
	return error_new(72, GetLastError(), NULL);
}
*/

void g2d_window_style_set(void *const data, const int wn, const int hn, const int wx, const int hx, const int b, const int d, const int r, const int l) {
	window_data_t *const wnd_data = (window_data_t*)data;
	wnd_data[0].config.width_min = wn;
	wnd_data[0].config.height_min = hn;
	wnd_data[0].config.width_max = wx;
	wnd_data[0].config.height_max = hx;
	wnd_data[0].config.borderless = b;
	wnd_data[0].config.dragable = d;
	wnd_data[0].config.resizable = r;
}

/*
void *g2d_window_title_set(void *const data, void *const title) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (!SetWindowText(wnd_data[0].wnd.hndl, (LPCTSTR)title))
		return error_new(68, GetLastError(), NULL);
	return NULL;
}

void *g2d_mouse_pos_set(void *const data, const int x, const int y) {
	window_data_t *const wnd_data = (window_data_t*)data;
	POINT point = {0, 0};
	ClientToScreen(wnd_data[0].wnd.hndl, &point);
	if (!SetCursorPos(point.x + x, point.y + y))
		return error_new(69, 0, NULL);
	return NULL;
}
*/
