/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static PROC get_proc(LPCSTR const func_name, void **const err) {
	PROC proc = NULL;
	if (err[0] == NULL) {
		// wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions).
		proc = wglGetProcAddress(func_name);
		const DWORD err_win32 = GetLastError();
		if (err_win32) {
			char *const err_str = (char*)malloc(sizeof(char) * 100);
			proc = NULL;
			if (err_str) {
				const size_t length0 = strlen(func_name) + 1;
				memcpy(err_str, func_name, length0);
				err[0] = error_new(2, err_win32, err_str);
			} else {
				err[0] = error_new(1, 0, NULL);
			}
		}
	}
	return proc;
}

static void module_init(void **const err) {
	if (err[0] == NULL) {
		instance = GetModuleHandle(NULL);
		if (!instance)
			err[0] = error_new(2, GetLastError(), NULL);
	}
}

static void dummy_class_init(window_t *const dummy, void **const err) {
	if (err[0] == NULL) {
		dummy[0].cls.cbSize = sizeof(WNDCLASSEX);
		dummy[0].cls.style = CS_OWNDC | CS_HREDRAW | CS_VREDRAW;
		dummy[0].cls.lpfnWndProc = DefWindowProc;
		dummy[0].cls.cbClsExtra = 0;
		dummy[0].cls.cbWndExtra = 0;
		dummy[0].cls.hInstance = instance;
		dummy[0].cls.hIcon = NULL;
		dummy[0].cls.hCursor = NULL;
		dummy[0].cls.hbrBackground = NULL;
		dummy[0].cls.lpszMenuName = NULL;
		dummy[0].cls.lpszClassName = TEXT("g2d_dummy");
		dummy[0].cls.hIconSm = NULL;
		if (RegisterClassEx(&dummy[0].cls) == INVALID_ATOM)
			err[0] = error_new(3, GetLastError(), NULL);
	}
}

static void dummy_window_create(window_t *const dummy, void **const err) {
	if (err[0] == NULL) {
		dummy[0].hndl = CreateWindow(dummy[0].cls.lpszClassName, TEXT("Dummy"), WS_OVERLAPPEDWINDOW, 0, 0, 1, 1, NULL, NULL, dummy[0].cls.hInstance, NULL);
		if (!dummy[0].hndl)
			err[0] = error_new(4, GetLastError(), NULL);
	}
}

static void dummy_context_init(window_t *const dummy, void **const err) {
	if (err[0] == NULL) {
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
					if (!dummy[0].ctx.rc) {
						err[0] = error_new(8, GetLastError(), NULL);
					}
				} else {
					err[0] = error_new(7, GetLastError(), NULL);
				}
			} else {
				err[0] = error_new(6, GetLastError(), NULL);
			}
		} else {
			err[0] = error_new(5, GetLastError(), NULL);
		}
	}
}

static void dummy_make_context_current(window_t *const dummy, void **const err) {
	if (err[0] == NULL)
		if (!wglMakeCurrent(dummy[0].ctx.dc, dummy[0].ctx.rc))
			err[0] = error_new(9, GetLastError(), NULL);
}

static void wgl_functions_init(void **const err) {
	wglChoosePixelFormatARB    = (PFNWGLCHOOSEPIXELFORMATARBPROC)    get_proc("wglChoosePixelFormatARB",    err);
	wglCreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC) get_proc("wglCreateContextAttribsARB", err);
	wglSwapIntervalEXT         = (PFNWGLSWAPINTERVALEXTPROC)         get_proc("wglSwapIntervalEXT",         err);
	wglGetSwapIntervalEXT      = (PFNWGLGETSWAPINTERVALEXTPROC)      get_proc("wglGetSwapIntervalEXT",      err);
}

static void ogl_functions_init(void **const err) {
	glCreateShader            = (PFNGLCREATESHADERPROC)            get_proc("glCreateShader",            err);
	glShaderSource            = (PFNGLSHADERSOURCEPROC)            get_proc("glShaderSource",            err);
	glCompileShader           = (PFNGLCOMPILESHADERPROC)           get_proc("glCompileShader",           err);
	glGetShaderiv             = (PFNGLGETSHADERIVPROC)             get_proc("glGetShaderiv",             err);
	glGetShaderInfoLog        = (PFNGLGETSHADERINFOLOGPROC)        get_proc("glGetShaderInfoLog",        err);
	glCreateProgram           = (PFNGLCREATEPROGRAMPROC)           get_proc("glCreateProgram",           err);
	glAttachShader            = (PFNGLATTACHSHADERPROC)            get_proc("glAttachShader",            err);
	glLinkProgram             = (PFNGLLINKPROGRAMPROC)             get_proc("glLinkProgram",             err);
	glValidateProgram         = (PFNGLVALIDATEPROGRAMPROC)         get_proc("glValidateProgram",         err);
	glGetProgramiv            = (PFNGLGETPROGRAMIVPROC)            get_proc("glGetProgramiv",            err);
	glGetProgramInfoLog       = (PFNGLGETPROGRAMINFOLOGPROC)       get_proc("glGetProgramInfoLog",       err);
	glGenBuffers              = (PFNGLGENBUFFERSPROC)              get_proc("glGenBuffers",              err);
	glGenVertexArrays         = (PFNGLGENVERTEXARRAYSPROC)         get_proc("glGenVertexArrays",         err);
	glGetAttribLocation       = (PFNGLGETATTRIBLOCATIONPROC)       get_proc("glGetAttribLocation",       err);
	glBindVertexArray         = (PFNGLBINDVERTEXARRAYPROC)         get_proc("glBindVertexArray",         err);
	glEnableVertexAttribArray = (PFNGLENABLEVERTEXATTRIBARRAYPROC) get_proc("glEnableVertexAttribArray", err);
	glVertexAttribPointer     = (PFNGLVERTEXATTRIBPOINTERPROC)     get_proc("glVertexAttribPointer",     err);
	glBindBuffer              = (PFNGLBINDBUFFERPROC)              get_proc("glBindBuffer",              err);
	glBufferData              = (PFNGLBUFFERDATAPROC)              get_proc("glBufferData",              err);
	glGetVertexAttribPointerv = (PFNGLGETVERTEXATTRIBPOINTERVPROC) get_proc("glGetVertexAttribPointerv", err);
	glUseProgram              = (PFNGLUSEPROGRAMPROC)              get_proc("glUseProgram",              err);
	glDeleteVertexArrays      = (PFNGLDELETEVERTEXARRAYSPROC)      get_proc("glDeleteVertexArrays",      err);
	glDeleteBuffers           = (PFNGLDELETEBUFFERSPROC)           get_proc("glDeleteBuffers",           err);
	glDeleteProgram           = (PFNGLDELETEPROGRAMPROC)           get_proc("glDeleteProgram",           err);
	glDeleteShader            = (PFNGLDELETESHADERPROC)            get_proc("glDeleteShader",            err);
	glGetUniformLocation      = (PFNGLGETUNIFORMLOCATIONPROC)      get_proc("glGetUniformLocation",      err);
	glUniformMatrix3fv        = (PFNGLUNIFORMMATRIX3FVPROC)        get_proc("glUniformMatrix3fv",        err);
	glUniformMatrix4fv        = (PFNGLUNIFORMMATRIX4FVPROC)        get_proc("glUniformMatrix4fv",        err);
	glUniformMatrix2x3fv      = (PFNGLUNIFORMMATRIX2X3FVPROC)      get_proc("glUniformMatrix2x3fv",      err);
	glGenerateMipmap          = (PFNGLGENERATEMIPMAPPROC)          get_proc("glGenerateMipmap",          err);
	glActiveTexture           = (PFNGLACTIVETEXTUREPROC)           get_proc("glActiveTexture",           err);
}

static void dummy_destroy(window_t *const dummy, void **const err) {
	if (dummy[0].ctx.rc) {
		if (!wglMakeCurrent(NULL, NULL) && err[0] == NULL)
			err[0] = error_new(10, GetLastError(), NULL);
		wglDeleteContext(dummy[0].ctx.rc);
	}
	if (dummy[0].ctx.dc) {
		ReleaseDC(dummy[0].hndl, dummy[0].ctx.dc);
	}
	if (dummy[0].hndl) {
		if (!DestroyWindow(dummy[0].hndl) && err[0] == NULL)
			err[0] = error_new(11, GetLastError(), NULL);
	}
	if (dummy[0].cls.lpszClassName) {
		if (!UnregisterClass(dummy[0].cls.lpszClassName, dummy[0].cls.hInstance) && err[0] == NULL)
			err[0] = error_new(12, GetLastError(), NULL);
	}
}

void *g2d_init() {
	void *err = NULL;
	if (!initialized) {
		window_t dummy;
		ZeroMemory((void*)&dummy, sizeof(window_t));
		module_init(&err);
		dummy_class_init(&dummy, &err);
		dummy_window_create(&dummy, &err);
		dummy_context_init(&dummy, &err);
		dummy_make_context_current(&dummy, &err);
		wgl_functions_init(&err);
		ogl_functions_init(&err);
		dummy_destroy(&dummy, &err);
		initialized = (BOOL)(err == NULL);
	}
	return err;
}
