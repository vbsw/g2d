/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

void g2d_init(void **const data, int *const xts, long long *const err1, long long *const err2) {
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
										glGetIntegerv(GL_MAX_TEXTURE_SIZE, xts);
										/* wgl functions */
										LOAD_FUNC(PFNWGLCHOOSEPIXELFORMATARBPROC, wglChoosePixelFormatARB, 1100)
										LOAD_FUNC(PFNWGLCREATECONTEXTATTRIBSARBPROC, wglCreateContextAttribsARB, 1101)
										LOAD_FUNC(PFNWGLSWAPINTERVALEXTPROC, wglSwapIntervalEXT, 1102)
										LOAD_FUNC(PFNWGLGETSWAPINTERVALEXTPROC, wglGetSwapIntervalEXT, 1103)
										/* ogl functions */
										LOAD_FUNC(PFNGLCREATESHADERPROC, glCreateShader, 1104)
										LOAD_FUNC(PFNGLSHADERSOURCEPROC, glShaderSource, 1105)
										LOAD_FUNC(PFNGLCOMPILESHADERPROC, glCompileShader, 1106)
										LOAD_FUNC(PFNGLGETSHADERIVPROC, glGetShaderiv, 1107)
										LOAD_FUNC(PFNGLGETSHADERINFOLOGPROC, glGetShaderInfoLog, 1108)
										LOAD_FUNC(PFNGLCREATEPROGRAMPROC, glCreateProgram, 1109)
										LOAD_FUNC(PFNGLATTACHSHADERPROC, glAttachShader, 1110)
										LOAD_FUNC(PFNGLLINKPROGRAMPROC, glLinkProgram, 1111)
										LOAD_FUNC(PFNGLVALIDATEPROGRAMPROC, glValidateProgram, 1112)
										LOAD_FUNC(PFNGLGETPROGRAMIVPROC, glGetProgramiv, 1113)
										LOAD_FUNC(PFNGLGETPROGRAMINFOLOGPROC, glGetProgramInfoLog, 1114)
										LOAD_FUNC(PFNGLGENBUFFERSPROC, glGenBuffers, 1115)
										LOAD_FUNC(PFNGLGENVERTEXARRAYSPROC, glGenVertexArrays, 1116)
										LOAD_FUNC(PFNGLGETATTRIBLOCATIONPROC, glGetAttribLocation, 1117)
										LOAD_FUNC(PFNGLBINDVERTEXARRAYPROC, glBindVertexArray, 1118)
										LOAD_FUNC(PFNGLENABLEVERTEXATTRIBARRAYPROC, glEnableVertexAttribArray, 1119)
										LOAD_FUNC(PFNGLVERTEXATTRIBPOINTERPROC, glVertexAttribPointer, 1120)
										LOAD_FUNC(PFNGLBINDBUFFERPROC, glBindBuffer, 1121)
										LOAD_FUNC(PFNGLBUFFERDATAPROC, glBufferData, 1122)
										LOAD_FUNC(PFNGLBUFFERSUBDATAPROC, glBufferSubData, 1123)
										LOAD_FUNC(PFNGLGETVERTEXATTRIBPOINTERVPROC, glGetVertexAttribPointerv, 1124)
										LOAD_FUNC(PFNGLUSEPROGRAMPROC, glUseProgram, 1125)
										LOAD_FUNC(PFNGLDELETEVERTEXARRAYSPROC, glDeleteVertexArrays, 1126)
										LOAD_FUNC(PFNGLDELETEBUFFERSPROC, glDeleteBuffers, 1127)
										LOAD_FUNC(PFNGLDELETEPROGRAMPROC, glDeleteProgram, 1128)
										LOAD_FUNC(PFNGLDELETESHADERPROC, glDeleteShader, 1129)
										LOAD_FUNC(PFNGLGETUNIFORMLOCATIONPROC, glGetUniformLocation, 1130)
										LOAD_FUNC(PFNGLUNIFORMMATRIX3FVPROC, glUniformMatrix3fv, 1131)
										LOAD_FUNC(PFNGLUNIFORM1FVPROC, glUniform1fv, 1132)
										LOAD_FUNC(PFNGLUNIFORMMATRIX4FVPROC, glUniformMatrix4fv, 1133)
										LOAD_FUNC(PFNGLUNIFORMMATRIX2X3FVPROC, glUniformMatrix2x3fv, 1134)
										LOAD_FUNC(PFNGLGENERATEMIPMAPPROC, glGenerateMipmap, 1135)
										LOAD_FUNC(PFNGLACTIVETEXTUREPROC, glActiveTexture, 1136)
										/* destroy dummy */
										if (!wglMakeCurrent(NULL, NULL) && err1[0] == 0) {
											err1[0] = 1008; err2[0] = (long long)GetLastError();
										}
										if (!wglDeleteContext(dummy_rc) && err1[0] == 0) {
											err1[0] = 1009; err2[0] = (long long)GetLastError();
										}
										ReleaseDC(dummy_hndl, dummy_dc);
										if (!DestroyWindow(dummy_hndl) && err1[0] == 0) {
											err1[0] = 1010; err2[0] = (long long)GetLastError();
										}
										if (!UnregisterClass(class_name_dummy, instance) && err1[0] == 0) {
											err1[0] = 1011; err2[0] = (long long)GetLastError();
										}
										if (err1[0] == 0) {
											initialized = TRUE;
										} else {
											free(engine);
										}
									} else {
										err1[0] = 1007; err2[0] = (long long)GetLastError();
										wglDeleteContext(dummy_rc); ReleaseDC(dummy_hndl, dummy_dc);
										DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
									}
								} else {
									err1[0] = 1006; err2[0] = (long long)GetLastError();
									ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
								}
							} else {
								err1[0] = 1005; err2[0] = (long long)GetLastError();
								ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
							}
						} else {
							err1[0] = 1004; err2[0] = (long long)GetLastError();
							ReleaseDC(dummy_hndl, dummy_dc); DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
						}
					} else {
						err1[0] = 1003;
						DestroyWindow(dummy_hndl); UnregisterClass(class_name_dummy, instance);
					}
				} else {
					err1[0] = 1002; err2[0] = (long long)GetLastError();
					UnregisterClass(class_name_dummy, instance);
				}
			} else {
				err1[0] = 1001; err2[0] = (long long)GetLastError();
			}
		} else {
			err1[0] = 1000; err2[0] = (long long)GetLastError();
		}
	}
}
