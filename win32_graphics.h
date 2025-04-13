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

static void prog_use(const GLuint id, const int err_a, const int err_b, long long *const err1) {
	glUseProgram(id);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_b;
	}
}

static GLuint rects_create(const GLuint vs_id, const GLuint fs_id, long long *const err1, char **const err_nfo) {
	if (err1[0] == 0) {
		const GLuint id = glCreateProgram();
		if (id) {
			shader_attach(id, vs_id, G2D_ERR_1002006, G2D_ERR_1002007, err1);
			if (err1[0] == 0) {
				shader_attach(id, fs_id, G2D_ERR_1002008, G2D_ERR_1002009, err1);
				if (err1[0] == 0) {
					glLinkProgram(id);
					program_check(id, GL_LINK_STATUS, G2D_ERR_1002010, err1, err_nfo);
				}
			}
		} else {
			err1[0] = G2D_ERR_1002011;
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

static void rects_enable(const GLuint prog_id, const GLint proj_unif, const GLuint vao, const GLuint vbo, const float *const projection_mat, long long *const err1) {
	prog_use(prog_id, G2D_ERR_1002012, G2D_ERR_1002013, err1);
	if (err1[0] == 0) {
		bind_vao(vao, G2D_ERR_1002014, err1);
		if (err1[0] == 0) {
			glUniformMatrix4fv(proj_unif, 1, GL_FALSE, projection_mat);
			bind_vbo(vbo, G2D_ERR_1002015, G2D_ERR_1002016, err1);
		}
	}
}

static void buffer_sub_data(const GLsizeiptr size, const void *const data, const int err_a, const int err_b, const int err_c, long long *const err1) {
	glBufferSubData(GL_ARRAY_BUFFER, 0, size, data);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err1[0] = err_a;
	} else if (err_enum == GL_INVALID_OPERATION) {
		err1[0] = err_b;
	} else if (err_enum == GL_INVALID_VALUE) {
		err1[0] = err_c;
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

static void bind_texture(const GLuint texture, const int err_a, const int err_b, const int err_c, long long *const err1) {
	glBindTexture(GL_TEXTURE_2D, texture);
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
		const GLuint vs_id = shader_create(GL_VERTEX_SHADER, vs_rect_str, G2D_ERR_1002002, G2D_ERR_1002003, err1, err_nfo);
		if (err1[0] == 0) {
			const GLuint fs_id = shader_create(GL_FRAGMENT_SHADER, fs_rect_str, G2D_ERR_1002004, G2D_ERR_1002005, err1, err_nfo);
			if (err1[0] == 0) {
				const size_t length = 16000;
				wnd_data[0].rects.max_length = (GLuint)length;
				wnd_data[0].rects.id = rects_create(vs_id, fs_id, err1, err_nfo);
				wnd_data[0].rects.pos_att = att_location(wnd_data[0].rects.id, "positionIn", G2D_ERR_1002023, err1);
				wnd_data[0].rects.col_att = att_location(wnd_data[0].rects.id, "colorIn", G2D_ERR_1002024, err1);
				wnd_data[0].rects.proj_unif = unf_location(wnd_data[0].rects.id, "projection", G2D_ERR_1002025, G2D_ERR_1002026, err1);
				if (err1[0] == 0) {
					GLuint objs[3]; glGenVertexArrays(1, objs); glGenBuffers(2, &objs[1]);
					wnd_data[0].rects.vao = objs[0]; wnd_data[0].rects.vbo = objs[1]; wnd_data[0].rects.ebo = objs[2];
					bind_vao(wnd_data[0].rects.vao, G2D_ERR_1002027, err1);
					enable_attr(wnd_data[0].rects.pos_att, G2D_ERR_1002028, G2D_ERR_1002029, err1);
					enable_attr(wnd_data[0].rects.col_att, G2D_ERR_1002030, G2D_ERR_1002031, err1);
					bind_vbo(wnd_data[0].rects.vbo, G2D_ERR_1002032, G2D_ERR_1002033, err1);
					buffer_data(GL_ARRAY_BUFFER, sizeof(float) * length * 4 * (2+4), NULL, GL_DYNAMIC_DRAW, G2D_ERR_1002034, G2D_ERR_1002035, G2D_ERR_1002036, G2D_ERR_1002037, err1);
					vertex_att_pointer(wnd_data[0].rects.pos_att, 2, sizeof(float) * (2+4), (void*)(sizeof(float) * 0), G2D_ERR_1002042, G2D_ERR_1002043, G2D_ERR_1002044, err1);
					vertex_att_pointer(wnd_data[0].rects.col_att, 4, sizeof(float) * (2+4), (void*)(sizeof(float) * 2), G2D_ERR_1002045, G2D_ERR_1002046, G2D_ERR_1002047, err1);
					bind_ebo(wnd_data[0].rects.ebo, G2D_ERR_1002048, G2D_ERR_1002049, err1);
					if (err1[0] == 0) {
						unsigned int *indices = (unsigned int*)malloc(sizeof(unsigned int) * length * (3+3));
						if (indices) {
							wnd_data[0].rects.buffer = (float*)malloc(sizeof(float) * length * 4 * (2+4));
							if (wnd_data[0].rects.buffer) {
								size_t i;
								for (i = 0; i < length; i++) {
									const size_t offs = i * (3+3);
									const size_t index = i * 4;
									indices[offs] = index; indices[offs+1] = index+1; indices[offs+2] = index+2; indices[offs+3] = index+2; indices[offs+4] = index+1; indices[offs+5] = index+3;
								}
								buffer_data(GL_ELEMENT_ARRAY_BUFFER, sizeof(unsigned int) * length * (3+3), indices, GL_STATIC_DRAW, G2D_ERR_1002038, G2D_ERR_1002039, G2D_ERR_1002040, G2D_ERR_1002041, err1);
							} else {
								err1[0] = G2D_ERR_0000018;
							}
							free((void*)indices);
						} else {
							err1[0] = G2D_ERR_0000017;
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
		err1[0] = G2D_ERR_1002001, err2[0] = (long long)GetLastError();
	}
}

void g2d_gfx_release(void *const data, long long *const err1, long long *const err2) {
	if (!wglMakeCurrent(NULL, NULL))
		err1[0] = G2D_ERR_1002051, err2[0] = (long long)GetLastError();
}

void g2d_gfx_draw(void *const data, const int w, const int h, const int i, const float r, const float g, const float b,
	float **const buffs, const int *const bs, void **const procs, const int l, long long *const err1, long long *const err2) {
	int k;
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
	for (k = 0; k < l && err1[0] == 0; k++) {
		gfx_draw_t *const draw = (gfx_draw_t*) procs[k];
		draw(data, buffs[k], bs[k], err1);
	}
	if (err1[0] == 0)
		if (!SwapBuffers(wnd_data[0].wnd.dc))
			err1[0] = G2D_ERR_1002050, err2[0] = (long long)GetLastError();
}

void g2d_gfx_draw_rectangles(void *const data, float *const rects, const int total, long long *const err1) {
	int i, drawn;
	window_data_t *const wnd_data = (window_data_t*)data;
	const int length = (int)wnd_data[0].rects.max_length;
	float *const buffer = wnd_data[0].rects.buffer;
	rects_enable(wnd_data[0].rects.id, wnd_data[0].rects.proj_unif, wnd_data[0].rects.vao, wnd_data[0].rects.vbo, wnd_data[0].gfx.projection_mat, err1);
	for (i = 0, drawn = 0; err1[0] == 0 && drawn < total; drawn += length) {
		int k;
		const int limit = drawn + length > total ? total - drawn : length;
		for (k = 0; k < limit; i++) {
			const int offs = k * 4 * (2+4); const int index = i * 8;
			const float x = rects[index], y = rects[index+1], w = rects[index+2], h = rects[index+3], r = rects[index+4], g = rects[index+5], b = rects[index+6], a = rects[index+7];
			buffer[offs+0] = x;
			buffer[offs+1] = y;
			buffer[offs+2] = r;
			buffer[offs+3] = g;
			buffer[offs+4] = b;
			buffer[offs+5] = a;
			buffer[offs+6] = x + w;
			buffer[offs+7] = y;
			buffer[offs+8] = r;
			buffer[offs+9] = g;
			buffer[offs+10] = b;
			buffer[offs+11] = a;
			buffer[offs+12] = x;
			buffer[offs+13] = y + h;
			buffer[offs+14] = r;
			buffer[offs+15] = g;
			buffer[offs+16] = b;
			buffer[offs+17] = a;
			buffer[offs+18] = x + w;
			buffer[offs+19] = y + h;
			buffer[offs+20] = r;
			buffer[offs+21] = g;
			buffer[offs+22] = b;
			buffer[offs+23] = a;
			k++;
		}
		buffer_sub_data(sizeof(float) * limit * 4 * (2+4), buffer, G2D_ERR_1002020, G2D_ERR_1002021, G2D_ERR_1002022, err1);
		draw_elements(limit * 6, G2D_ERR_1002017, G2D_ERR_1002018, G2D_ERR_1002019, err1);
	}
}

void g2d_gfx_gen_tex(void *const data, const void *const tex, const int w, const int h, int *const tex_id, long long *const err1) {
	GLuint texture; glGenTextures(1, &texture);
	bind_texture(texture, G2D_ERR_1002052, G2D_ERR_1002053, G2D_ERR_1002054, err1);
	if (err1[0] == 0) {
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_BORDER);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_BORDER);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
		glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
		glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, (GLsizei)w, (GLsizei)h, 0, GL_RGBA, GL_UNSIGNED_BYTE, tex);
		const GLenum err_enum = glGetError();
		if (err_enum == GL_NO_ERROR) {
			tex_id[0] = (int)texture;
		} else if (err_enum == GL_INVALID_ENUM) {
			err1[0] = G2D_ERR_1002055;
		} else if (err_enum == GL_INVALID_VALUE) {
			err1[0] = G2D_ERR_1002056;
		} else if (err_enum == GL_INVALID_OPERATION) {
			err1[0] = G2D_ERR_1002057;
		}
	}
}
