/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static PROC get_proc(LPCSTR const func_name, int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
	PROC proc = NULL;
	if (err_num[0] == 0) {
		// wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions).
		proc = wglGetProcAddress(func_name);
		const DWORD last_error = GetLastError();
		if (last_error) {
			char *const func_name_copy = str_copy(func_name);
			proc = NULL;
			if (func_name_copy)
				ERR_NEW3(17, last_error, func_name_copy)
			else
				ERR_NEW2(-17, last_error)
		}
	}
	return proc;
}

static void module_init(int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0) {
		instance = GetModuleHandle(NULL);
		if (!instance)
			ERR_NEW2(1, GetLastError())
	}
}

static BOOL class_register_dummy(int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0) {
		WNDCLASSEX cls;
		ZeroMemory(&cls, sizeof(WNDCLASSEX));
		cls.cbSize = sizeof(WNDCLASSEX);
		cls.style = CS_OWNDC | CS_HREDRAW | CS_VREDRAW;
		cls.lpfnWndProc = DefWindowProc;
		cls.hInstance = instance;
		cls.lpszClassName = class_name_dummy;
		if (RegisterClassEx(&cls) == INVALID_ATOM)
			ERR_NEW2(10, GetLastError())
		else
			return TRUE;
	}
	return FALSE;
}

static void window_create_dummy(window_t *const dummy, int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0) {
		dummy[0].hndl = CreateWindow(class_name_dummy, TEXT("Dummy"), WS_OVERLAPPEDWINDOW, 0, 0, 1, 1, NULL, NULL, instance, NULL);
		if (!dummy[0].hndl)
			ERR_NEW2(11, GetLastError())
	}
}

static void context_init_dummy(window_t *const dummy, int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0) {
		dummy[0].ctx.dc = GetDC(dummy[0].hndl);
		if (dummy[0].ctx.dc) {
			int pixelFormat;
			PIXELFORMATDESCRIPTOR pixelFormatDesc;
			ZeroMemory(&pixelFormatDesc, sizeof(PIXELFORMATDESCRIPTOR));
			pixelFormatDesc.nSize = sizeof(PIXELFORMATDESCRIPTOR);
			pixelFormatDesc.nVersion = 1;
			pixelFormatDesc.dwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL;
			pixelFormatDesc.iPixelType = PFD_TYPE_RGBA;
			pixelFormatDesc.cColorBits = 32;
			pixelFormatDesc.cAlphaBits = 8;
			pixelFormatDesc.cDepthBits = 24;
			pixelFormat = ChoosePixelFormat(dummy[0].ctx.dc, &pixelFormatDesc);
			if (pixelFormat) {
				if (SetPixelFormat(dummy[0].ctx.dc, pixelFormat, &pixelFormatDesc)) {
					dummy[0].ctx.rc = wglCreateContext(dummy[0].ctx.dc);
					if (!dummy[0].ctx.rc)
						ERR_NEW2(15, GetLastError())
				} else {
					ERR_NEW2(14, GetLastError())
				}
			} else {
				ERR_NEW2(13, GetLastError())
			}
		} else {
			ERR_NEW2(12, GetLastError())
		}
	}
}

static void context_make_current_dummy(window_t *const dummy, int *const err_num, g2d_ul_t *const err_win32) {
	if (err_num[0] == 0)
		if (!wglMakeCurrent(dummy[0].ctx.dc, dummy[0].ctx.rc))
			ERR_NEW2(16, GetLastError())
}

static void wgl_functions_init(int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
	wglChoosePixelFormatARB    = (PFNWGLCHOOSEPIXELFORMATARBPROC)    get_proc("wglChoosePixelFormatARB",    err_num, err_win32, err_str);
	wglCreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC) get_proc("wglCreateContextAttribsARB", err_num, err_win32, err_str);
	wglSwapIntervalEXT         = (PFNWGLSWAPINTERVALEXTPROC)         get_proc("wglSwapIntervalEXT",         err_num, err_win32, err_str);
	wglGetSwapIntervalEXT      = (PFNWGLGETSWAPINTERVALEXTPROC)      get_proc("wglGetSwapIntervalEXT",      err_num, err_win32, err_str);
}

static void ogl_functions_init(int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
	glCreateShader            = (PFNGLCREATESHADERPROC)            get_proc("glCreateShader",            err_num, err_win32, err_str);
	glShaderSource            = (PFNGLSHADERSOURCEPROC)            get_proc("glShaderSource",            err_num, err_win32, err_str);
	glCompileShader           = (PFNGLCOMPILESHADERPROC)           get_proc("glCompileShader",           err_num, err_win32, err_str);
	glGetShaderiv             = (PFNGLGETSHADERIVPROC)             get_proc("glGetShaderiv",             err_num, err_win32, err_str);
	glGetShaderInfoLog        = (PFNGLGETSHADERINFOLOGPROC)        get_proc("glGetShaderInfoLog",        err_num, err_win32, err_str);
	glCreateProgram           = (PFNGLCREATEPROGRAMPROC)           get_proc("glCreateProgram",           err_num, err_win32, err_str);
	glAttachShader            = (PFNGLATTACHSHADERPROC)            get_proc("glAttachShader",            err_num, err_win32, err_str);
	glLinkProgram             = (PFNGLLINKPROGRAMPROC)             get_proc("glLinkProgram",             err_num, err_win32, err_str);
	glValidateProgram         = (PFNGLVALIDATEPROGRAMPROC)         get_proc("glValidateProgram",         err_num, err_win32, err_str);
	glGetProgramiv            = (PFNGLGETPROGRAMIVPROC)            get_proc("glGetProgramiv",            err_num, err_win32, err_str);
	glGetProgramInfoLog       = (PFNGLGETPROGRAMINFOLOGPROC)       get_proc("glGetProgramInfoLog",       err_num, err_win32, err_str);
	glGenBuffers              = (PFNGLGENBUFFERSPROC)              get_proc("glGenBuffers",              err_num, err_win32, err_str);
	glGenVertexArrays         = (PFNGLGENVERTEXARRAYSPROC)         get_proc("glGenVertexArrays",         err_num, err_win32, err_str);
	glGetAttribLocation       = (PFNGLGETATTRIBLOCATIONPROC)       get_proc("glGetAttribLocation",       err_num, err_win32, err_str);
	glBindVertexArray         = (PFNGLBINDVERTEXARRAYPROC)         get_proc("glBindVertexArray",         err_num, err_win32, err_str);
	glEnableVertexAttribArray = (PFNGLENABLEVERTEXATTRIBARRAYPROC) get_proc("glEnableVertexAttribArray", err_num, err_win32, err_str);
	glVertexAttribPointer     = (PFNGLVERTEXATTRIBPOINTERPROC)     get_proc("glVertexAttribPointer",     err_num, err_win32, err_str);
	glBindBuffer              = (PFNGLBINDBUFFERPROC)              get_proc("glBindBuffer",              err_num, err_win32, err_str);
	glBufferData              = (PFNGLBUFFERDATAPROC)              get_proc("glBufferData",              err_num, err_win32, err_str);
	glGetVertexAttribPointerv = (PFNGLGETVERTEXATTRIBPOINTERVPROC) get_proc("glGetVertexAttribPointerv", err_num, err_win32, err_str);
	glUseProgram              = (PFNGLUSEPROGRAMPROC)              get_proc("glUseProgram",              err_num, err_win32, err_str);
	glDeleteVertexArrays      = (PFNGLDELETEVERTEXARRAYSPROC)      get_proc("glDeleteVertexArrays",      err_num, err_win32, err_str);
	glDeleteBuffers           = (PFNGLDELETEBUFFERSPROC)           get_proc("glDeleteBuffers",           err_num, err_win32, err_str);
	glDeleteProgram           = (PFNGLDELETEPROGRAMPROC)           get_proc("glDeleteProgram",           err_num, err_win32, err_str);
	glDeleteShader            = (PFNGLDELETESHADERPROC)            get_proc("glDeleteShader",            err_num, err_win32, err_str);
	glGetUniformLocation      = (PFNGLGETUNIFORMLOCATIONPROC)      get_proc("glGetUniformLocation",      err_num, err_win32, err_str);
	glUniformMatrix3fv        = (PFNGLUNIFORMMATRIX3FVPROC)        get_proc("glUniformMatrix3fv",        err_num, err_win32, err_str);
	glUniformMatrix4fv        = (PFNGLUNIFORMMATRIX4FVPROC)        get_proc("glUniformMatrix4fv",        err_num, err_win32, err_str);
	glUniformMatrix2x3fv      = (PFNGLUNIFORMMATRIX2X3FVPROC)      get_proc("glUniformMatrix2x3fv",      err_num, err_win32, err_str);
	glGenerateMipmap          = (PFNGLGENERATEMIPMAPPROC)          get_proc("glGenerateMipmap",          err_num, err_win32, err_str);
	glActiveTexture           = (PFNGLACTIVETEXTUREPROC)           get_proc("glActiveTexture",           err_num, err_win32, err_str);
}

static void window_destroy_dummy(window_t *const dummy, int *const err_num, g2d_ul_t *const err_win32) {
	if (dummy[0].ctx.rc) {
		if (dummy[0].ctx.rc == wglGetCurrentContext() && !wglMakeCurrent(NULL, NULL) && err_num[0] == 0)
			ERR_NEW2(18, GetLastError())
		if (!wglDeleteContext(dummy[0].ctx.rc) && err_num[0] == 0) {
			ERR_NEW2(19, GetLastError())
		}
	}
	if (dummy[0].ctx.dc) {
		ReleaseDC(dummy[0].hndl, dummy[0].ctx.dc);
	}
	if (dummy[0].hndl) {
		if (!DestroyWindow(dummy[0].hndl) && err_num[0] == 0)
			ERR_NEW2(20, GetLastError())
	}
}

static void class_unregister_dummy(int *const err_num, g2d_ul_t *const err_win32) {
	if (!UnregisterClass(class_name_dummy, instance) && err_num[0] == 0)
		ERR_NEW2(21, GetLastError())
}

void g2d_init(int *const err_num, g2d_ul_t *const err_win32, char **const err_str) {
	if (!initialized) {
		window_t dummy;
		ZeroMemory((void*)&dummy, sizeof(window_t));
		module_init(err_num, err_win32);
		const BOOL registered = class_register_dummy(err_num, err_win32);
		window_create_dummy(&dummy, err_num, err_win32);
		context_init_dummy(&dummy, err_num, err_win32);
		context_make_current_dummy(&dummy, err_num, err_win32);
		wgl_functions_init(err_num, err_win32, err_str);
		ogl_functions_init(err_num, err_win32, err_str);
		window_destroy_dummy(&dummy, err_num, err_win32);
		if (registered)
			class_unregister_dummy(err_num, err_win32);
		initialized = (BOOL)(err_num[0] == 0);
	}
}
