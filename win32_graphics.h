/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

static GLuint shader_create(const GLenum shader_type, LPCSTR shader, const int err_a, const int err_b, long long *const err1, char **const err_nfo) {
	const GLuint id = glCreateShader(shader_type);
	if (id) {
		GLint compiled; glShaderSource(id, 1, &shader, NULL); glCompileShader(id);
		glGetShaderiv(id, GL_COMPILE_STATUS, &compiled);
		if (compiled == GL_FALSE) {
			GLsizei err_len; err1[0] = err_b; glGetShaderiv(id, GL_INFO_LOG_LENGTH, &err_len);
			if (err_len > 0) {
				err_nfo[0] = (char*)malloc(err_len);
				if (err_nfo[0])
					glGetShaderInfoLog(id, err_len, &err_len, (GLchar*)err_nfo[0]);
			}
			glDeleteShader(id);
		}
	} else {
		err1[0] = err_a;
	}
	return id;
}

static void shader_attach(const GLuint prog_id, const GLuint shader_id, const int err_a, const int err_b, long long *const err1) {
	glAttachShader(prog_id, shader_id);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_b;
	}
}

static void program_check(const GLuint prog_id, const GLenum status, const int err, long long *const err1, char **const err_nfo) {
	GLint success; glGetProgramiv(prog_id, status, &success);
	if (success == GL_FALSE) {
		GLsizei err_len; glGetProgramiv(prog_id, GL_INFO_LOG_LENGTH, &err_len); err1[0] = err;
		if (err_len > 0) {
			err_nfo[0] = (char*)malloc(err_len);
			if (err_nfo[0])
				glGetProgramInfoLog(prog_id, err_len, &err_len, err_nfo[0]);
		}
	}
}

static GLuint rect_prog_create(const GLuint vs_id, const GLuint fs_id, long long *const err1, char **const err_nfo) {
	if (err1[0] == 0) {
		const GLuint id = glCreateProgram();
		if (id) {
			shader_attach(id, vs_id, 1005, 1006, err1);
			if (err1[0] == 0) {
				shader_attach(id, fs_id, 1007, 1008, err1);
				if (err1[0] == 0) {
					glLinkProgram(id);
					program_check(id, GL_LINK_STATUS, 1009, err1, err_nfo);
				}
			}
		} else {
			err1[0] = 1004;
		}
		return id;
	}
	return 0;
}

static void bind_vao(const GLuint vao, const int err, long long *const err1) {
	if (err1[0] == 0) {
		glBindVertexArray(vao);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err;
		}
	}
}

static void bind_vbo(const GLuint vbo, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		glBindBuffer(GL_ARRAY_BUFFER, vbo);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
	}
}

static void bind_ebo(const GLuint ebo, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, ebo);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
	}
}

static void enable_attr(const GLint attr, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		glEnableVertexAttribArray(attr);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
	}
}

static GLint att_location(const GLuint prog_id, LPCSTR const name, const int err, long long *const err1) {
	if (err1[0] == 0) {
		const GLint att_id = glGetAttribLocation(prog_id, name);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err;
		}
		return att_id;
	}
	return -1;
}

static GLint unf_location(const GLuint prog_id, LPCSTR const name, const int err_a, const int err_b, long long *const err1) {
	if (err1[0] == 0) {
		const GLint att_id = glGetUniformLocation(prog_id, name);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		}
		return att_id;
	}
	return -1;
}

static void buffer_data(const GLenum target, const GLsizeiptr size, const void *const data, const GLenum usage, const int err_a, const int err_b, const int err_c, const int err_d, long long *const err1) {
	if (err1[0] == 0) {
		glBufferData(target, size, data, usage);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_b;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_c;
		} else if (err_enum == GL_OUT_OF_MEMORY) {
			err1[0] = err_d;
		}
	}
}

static void vertex_att_pointer(const GLuint index, const GLint size, const GLsizei stride, const void *const pointer, const int err_a, const int err_b, const int err_c, long long *const err1) {
	if (err1[0] == 0) {
		glVertexAttribPointer(index, size, GL_FLOAT, GL_FALSE, stride, pointer);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_INVALID_VALUE) {
			err1[0] = err_a;
		} else if (err_enum == GL_INVALID_ENUM) {
			err1[0] = err_b;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = err_c;
		}
	}
}

static void draw_elements(const GLsizei count, const int err_a, const int err_b, const int err_c, long long *const err1) {
	glDrawElements(GL_TRIANGLES, count, GL_UNSIGNED_INT, 0);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_b;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_c;
	}
}

void g2d_gfx_init(void *const data, long long *const err1, long long *const err2, char **const err_nfo) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wglMakeCurrent(wnd_data[0].wnd.dc, wnd_data[0].wnd.rc)) {
		const GLuint vs_id = shader_create(GL_VERTEX_SHADER, vs_rect_str, 1000, 1001, err1, err_nfo);
		if (err1[0] == 0) {
			const GLuint fs_id = shader_create(GL_FRAGMENT_SHADER, fs_rect_str, 1002, 1003, err1, err_nfo);
			if (err1[0] == 0) {
				const size_t size = 16000;
				wnd_data[0].rects.max_size = (GLuint)size;
				wnd_data[0].rects.id = rect_prog_create(vs_id, fs_id, err1, err_nfo);
				wnd_data[0].rects.pos_att = att_location(wnd_data[0].rects.id, "positionIn", 1010, err1);
				wnd_data[0].rects.col_att = att_location(wnd_data[0].rects.id, "colorIn", 1011, err1);
				wnd_data[0].rects.proj_unif = unf_location(wnd_data[0].rects.id, "projection", 1012, 1013, err1);
				if (err1[0] == 0) {
					GLuint objs[3]; glGenVertexArrays(1, objs); glGenBuffers(2, &objs[1]);
					wnd_data[0].rects.vao = objs[0]; wnd_data[0].rects.vbo = objs[1]; wnd_data[0].rects.ebo = objs[2];
					bind_vao(wnd_data[0].rects.vao, 1014, err1);
					enable_attr(wnd_data[0].rects.pos_att, 1015, 1016, err1);
					enable_attr(wnd_data[0].rects.col_att, 1017, 1018, err1);
					bind_vbo(wnd_data[0].rects.vbo, 1019, 1020, err1);
					buffer_data(GL_ARRAY_BUFFER, sizeof(float) * size * 4 * (2+4), NULL, GL_DYNAMIC_DRAW, 1021, 1022, 1023, 1024, err1);
					vertex_att_pointer(wnd_data[0].rects.pos_att, 2, sizeof(float) * (2+4), (void*)(sizeof(float) * 0), 1025, 1026, 1027, err1);
					vertex_att_pointer(wnd_data[0].rects.col_att, 4, sizeof(float) * (2+4), (void*)(sizeof(float) * 2), 1028, 1029, 1030, err1);
					bind_ebo(wnd_data[0].rects.ebo, 1031, 1032, err1);
					if (err1[0] == 0) {
						unsigned int *indices = (unsigned int*)malloc(sizeof(unsigned int) * size * (2+4));
						if (indices) {
							wnd_data[0].rects.buffer = (float*)malloc(sizeof(float) * size * 4 * (2+4));
							if (wnd_data[0].rects.buffer) {
								size_t i;
								for (i = 0; i < size; i++) {
									const size_t offs = i * (3+3);
									const size_t index = i * 4;
									indices[offs] = index; indices[offs+1] = index+1; indices[offs+2] = index+2; indices[offs+3] = index+2; indices[offs+4] = index+1; indices[offs+5] = index+3;
								}
								buffer_data(GL_ELEMENT_ARRAY_BUFFER, sizeof(unsigned int) * size * (3+3), indices, GL_STATIC_DRAW, 1035, 1036, 1037, 1038, err1);
							} else {
								err1[0] = 1034;
							}
							free((void*)indices);
						} else {
							err1[0] = 1033;
						}
					}
				}
				glDeleteShader(fs_id);
			}
			glDeleteShader(vs_id);
		}
		if (err1[0] == 0) {
			glEnable(GL_BLEND);
			glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
		}
	} else {
		err1[0] = 220, err2[0] = (long long)GetLastError();
	}
}

void g2d_gfx_release(void *const data, long long *const err1, long long *const err2) {
	if (!wglMakeCurrent(NULL, NULL))
		err1[0] = 220, err2[0] = (long long)GetLastError();
}

void g2d_gfx_draw(void *const data, const int w, const int h, const int i, const float r, const float g, const float b, long long *const err1, long long *const err2) {
	window_data_t *const wnd_data = (window_data_t*)data;
	if (wnd_data[0].gfx.w != w || wnd_data[0].gfx.g != h) {
		wnd_data[0].gfx.w = w; wnd_data[0].gfx.h = h;
		wnd_data[0].gfx.projection_mat[0] = 2.0f / (float)w;
		wnd_data[0].gfx.projection_mat[5] = -2.0f / (float)h;
		glViewport((WORD)0, (WORD)0, (WORD)w, (WORD)h);
	}
	if (wnd_data[0].gfx.r != r || wnd_data[0].gfx.g != g || wnd_data[0].gfx.b != b) {
		wnd_data[0].gfx.r = r; wnd_data[0].gfx.g = g; wnd_data[0].gfx.b = b;
		glClearColor((GLfloat)r, (GLfloat)g, (GLfloat)b, 0.0);
	}
	if (wnd_data[0].gfx.i != i && wglSwapIntervalEXT) {
		wnd_data[0].gfx.i = i;
		wglSwapIntervalEXT(i);
	}
	glClear(GL_COLOR_BUFFER_BIT);
	if (!SwapBuffers(wnd_data[0].wnd.dc))
		err1[0] = 220, err2[0] = (long long)GetLastError();
}
