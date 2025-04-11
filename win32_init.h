/*
 *        Copyright 2023, 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

/* wglGetProcAddress could return -1, 1, 2 or 3 on failure (https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions). */
#define LOAD_WGL(t, n) if (err1[0] == 0) { PROC const proc = wglGetProcAddress(#n); const DWORD last_err = GetLastError(); if (last_err == 0) n = (t) proc; else { err1[0] = G2D_ERR_1000101; err2[0] = (long long)last_err; err_nfo[0] = #n; }}
#define LOAD_OGL(t, n) if (err1[0] == 0) { PROC const proc = wglGetProcAddress(#n); const DWORD last_err = GetLastError(); if (last_err == 0) n = (t) proc; else { err1[0] = G2D_ERR_1000102; err2[0] = (long long)last_err; err_nfo[0] = #n; }}

void g2d_init(int *const numbers, long long *const err1, long long *const err2, char **const err_nfo) {
	if (!initialized) {
		/* module */
		instance = GetModuleHandle(NULL);
		if (instance) {
			/* dummy class */
			WNDCLASSEX cls;
			ZeroMemory(&cls, sizeof(WNDCLASSEX));
			cls.cbSize = sizeof(WNDCLASSEX);
			cls.style = CS_OWNDC;
			cls.lpfnWndProc = DefWindowProc;
			cls.hInstance = instance;
			cls.lpszClassName = class_name_dummy;
			if (RegisterClassEx(&cls) != INVALID_ATOM) {
				/* dummy window */
				HWND const dummy_hndl = CreateWindow(class_name_dummy, TEXT("Dummy"), WS_OVERLAPPEDWINDOW, 0, 0, 1, 1, NULL, NULL, instance, NULL);
				if (dummy_hndl) {
					/* dummy context */
					HDC const dummy_dc = GetDC(dummy_hndl);
					if (dummy_dc) {
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
						pixelFormat = ChoosePixelFormat(dummy_dc, &pixelFormatDesc);
						if (pixelFormat) {
							if (SetPixelFormat(dummy_dc, pixelFormat, &pixelFormatDesc)) {
								HGLRC const dummy_rc = wglCreateContext(dummy_dc);
								if (dummy_rc) {
									if (wglMakeCurrent(dummy_dc, dummy_rc)) {
										glGetIntegerv(GL_MAX_TEXTURE_SIZE, &numbers[0]);
										glGetIntegerv(GL_MAX_TEXTURE_IMAGE_UNITS, &numbers[1]);
										LOAD_WGL(PFNWGLCHOOSEPIXELFORMATARBPROC,    wglChoosePixelFormatARB)
										LOAD_WGL(PFNWGLCREATECONTEXTATTRIBSARBPROC, wglCreateContextAttribsARB)
										LOAD_WGL(PFNWGLGETEXTENSIONSSTRINGARBPROC,  wglGetExtensionsStringARB)
										if (err1[0] == 0) {
											int begin = 0, end = 0, i;
											LPCSTR const extensions = (LPCSTR)wglGetExtensionsStringARB(dummy_dc);
											while (extensions[end] && (!&numbers[2] || !&numbers[3])) {
												for (end = begin; extensions[end] && extensions[end] != ' '; end++);
												for (i = begin; i < end && extensions[i] == "WGL_EXT_swap_control"[i-begin]; i++);
												if (i == end && "WGL_EXT_swap_control"[end-begin] == 0) {
													LOAD_WGL(PFNWGLSWAPINTERVALEXTPROC,    wglSwapIntervalEXT)
													LOAD_WGL(PFNWGLGETSWAPINTERVALEXTPROC, wglGetSwapIntervalEXT)
													numbers[2] = 1;
												}
												for (i = begin; i < end && extensions[i] == "WGL_EXT_swap_control_tear"[i-begin]; i++);
												if (i == end && "WGL_EXT_swap_control_tear"[end-begin] == 0) {
													numbers[3] = 1;
												}
												begin = end + 1;
											}
										}
										LOAD_OGL(PFNGLCREATESHADERPROC,             glCreateShader)
										LOAD_OGL(PFNGLSHADERSOURCEPROC,             glShaderSource)
										LOAD_OGL(PFNGLCOMPILESHADERPROC,            glCompileShader)
										LOAD_OGL(PFNGLGETSHADERIVPROC,              glGetShaderiv)
										LOAD_OGL(PFNGLGETSHADERINFOLOGPROC,         glGetShaderInfoLog)
										LOAD_OGL(PFNGLCREATEPROGRAMPROC,            glCreateProgram)
										LOAD_OGL(PFNGLATTACHSHADERPROC,             glAttachShader)
										LOAD_OGL(PFNGLLINKPROGRAMPROC,              glLinkProgram)
										LOAD_OGL(PFNGLVALIDATEPROGRAMPROC,          glValidateProgram)
										LOAD_OGL(PFNGLGETPROGRAMIVPROC,             glGetProgramiv)
										LOAD_OGL(PFNGLGETPROGRAMINFOLOGPROC,        glGetProgramInfoLog)
										LOAD_OGL(PFNGLGENBUFFERSPROC,               glGenBuffers)
										LOAD_OGL(PFNGLGENVERTEXARRAYSPROC,          glGenVertexArrays)
										LOAD_OGL(PFNGLGETATTRIBLOCATIONPROC,        glGetAttribLocation)
										LOAD_OGL(PFNGLBINDVERTEXARRAYPROC,          glBindVertexArray)
										LOAD_OGL(PFNGLENABLEVERTEXATTRIBARRAYPROC,  glEnableVertexAttribArray)
										LOAD_OGL(PFNGLVERTEXATTRIBPOINTERPROC,      glVertexAttribPointer)
										LOAD_OGL(PFNGLBINDBUFFERPROC,               glBindBuffer)
										LOAD_OGL(PFNGLBUFFERDATAPROC,               glBufferData)
										LOAD_OGL(PFNGLBUFFERSUBDATAPROC,            glBufferSubData)
										LOAD_OGL(PFNGLGETVERTEXATTRIBPOINTERVPROC,  glGetVertexAttribPointerv)
										LOAD_OGL(PFNGLUSEPROGRAMPROC,               glUseProgram)
										LOAD_OGL(PFNGLDELETEVERTEXARRAYSPROC,       glDeleteVertexArrays)
										LOAD_OGL(PFNGLDELETEBUFFERSPROC,            glDeleteBuffers)
										LOAD_OGL(PFNGLDELETEPROGRAMPROC,            glDeleteProgram)
										LOAD_OGL(PFNGLDELETESHADERPROC,             glDeleteShader)
										LOAD_OGL(PFNGLGETUNIFORMLOCATIONPROC,       glGetUniformLocation)
										LOAD_OGL(PFNGLUNIFORMMATRIX3FVPROC,         glUniformMatrix3fv)
										LOAD_OGL(PFNGLUNIFORM1FVPROC,               glUniform1fv)
										LOAD_OGL(PFNGLUNIFORMMATRIX4FVPROC,         glUniformMatrix4fv)
										LOAD_OGL(PFNGLUNIFORMMATRIX2X3FVPROC,       glUniformMatrix2x3fv)
										LOAD_OGL(PFNGLGENERATEMIPMAPPROC,           glGenerateMipmap)
										LOAD_OGL(PFNGLACTIVETEXTUREPROC,            glActiveTexture)
										/* destroy dummy */
										if (wglGetCurrentContext() == dummy_rc && !wglMakeCurrent(NULL, NULL) && err1[0] == 0) {
											err1[0] = G2D_ERR_1000009; err2[0] = (long long)GetLastError();
										}
										if (!wglDeleteContext(dummy_rc) && err1[0] == 0) {
											err1[0] = G2D_ERR_1000010; err2[0] = (long long)GetLastError();
										}
										ReleaseDC(dummy_hndl, dummy_dc);
										if (!DestroyWindow(dummy_hndl) && err1[0] == 0) {
											err1[0] = G2D_ERR_1000011; err2[0] = (long long)GetLastError();
										}
										if (!UnregisterClass(class_name_dummy, instance) && err1[0] == 0) {
											err1[0] = G2D_ERR_1000012; err2[0] = (long long)GetLastError();
										}
										initialized = (BOOL)(err1[0] == 0);
									} else {
										err1[0] = G2D_ERR_1000008; err2[0] = (long long)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
									}
								} else {
									err1[0] = G2D_ERR_1000007; err2[0] = (long long)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
								}
							} else {
								err1[0] = G2D_ERR_1000006; err2[0] = (long long)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
							}
						} else {
							err1[0] = G2D_ERR_1000005; err2[0] = (long long)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
						}
					} else {
						err1[0] = G2D_ERR_1000004;
						DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
					}
				} else {
					err1[0] = G2D_ERR_1000003; err2[0] = (long long)GetLastError();
					UnregisterClass(class_name_dummy, instance);
				}
			} else {
				err1[0] = G2D_ERR_1000002; err2[0] = (long long)GetLastError();
			}
		} else {
			err1[0] = G2D_ERR_1000001; err2[0] = (long long)GetLastError();
		}
	}
}
