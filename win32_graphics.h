/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

#define RECT1K_COUNT 1000

/* orthographic projection */
static const float mat_projection[4*4] = { 2.0f / 1.0f, 0.0f, 0.0f, 0.0f, 0.0f, -2.0f / 1.0f, 0.0f, 0.0f, 0.0f, 0.0f, -1.0f, 0.0f, -1.0f, 1.0f, 0.0f, 1.0f };
static LPCSTR const vs_rect_str = "#version 130\nin vec2 positionIn; out vec4 fragementColor; uniform float data[9]; uniform mat4 projection = mat4(1.0); vec4 pos = vec4((data[0] + positionIn.x * data[2]) * data[4], (data[1] + positionIn.y * data[3]) * data[4], 1.0, 1.0); void main() { gl_Position = projection * pos; fragementColor = vec4(data[5], data[6], data[7], data[8]); }";
static LPCSTR const vs_rect1k_str = "#version 130\nin vec2 positionIn; out vec4 fragementColor; uniform float data[1000*9]; uniform mat4 projection = mat4(1.0); void main() {int vba_offset = int(positionIn.y) / 2; float y = positionIn.y - float(vba_offset * 2); int index = vba_offset * 9; vec4 pos = vec4((data[index] + positionIn.x * data[index+2]) * data[index+4], (data[index+1] + y * data[index+3]) * data[index+4], 1.0, 1.0); gl_Position = projection * pos; fragementColor = vec4(data[index+5], data[index+6], data[index+7], data[index+8]); }";
static LPCSTR const fs_str = "#version 130\nin vec4 fragementColor; out vec4 color; void main() { color = fragementColor; }";

static GLuint shader_create(const GLenum shader_type, LPCSTR shader, const int err, int *const err_num, char **const err_str) {
	const GLuint id = glCreateShader(shader_type);
	if (id) {
		GLint compiled;
		glShaderSource(id, 1, &shader, NULL);
		glCompileShader(id);
		glGetShaderiv(id, GL_COMPILE_STATUS, &compiled);
		if (compiled == GL_FALSE) {
			GLsizei err_len;
			err_num[0] = err;
			glGetShaderiv(id, GL_INFO_LOG_LENGTH, &err_len);
			if (err_len > 0) {
				err_str[0] = malloc(err_len + 1);
				glGetShaderInfoLog(id, err_len, &err_len, (GLchar*)err_str[0]);
				err_str[0][err_len + 1] = 0;
			}
			glDeleteShader(id);
		}
	} else {
		err_num[0] = 1030;
	}
	return id;
}

static void program_check(const GLuint prog_id, const GLenum status, int *const err_num, char **const err_str) {
	GLint success;
	glGetProgramiv(prog_id, status, &success);
	if (success == GL_FALSE) {
		GLsizei err_len;
		glGetProgramiv(prog_id, GL_INFO_LOG_LENGTH, &err_len);
		err_num[0] = 1037;
		if (err_len > 0) {
			err_str[0] = malloc(err_len);
			glGetProgramInfoLog(prog_id, err_len, &err_len, err_str[0]);
		}
	}
}

static void create_program(program_t *const prog, LPCSTR const vs_str, const GLuint fs_id, const int err_a, const int err_b, int *const err_num, char **const err_str) {
	prog[0].vs_id = shader_create(GL_VERTEX_SHADER, vs_str, err_b, err_num, err_str);
	if (err_num[0] == 0) {
		prog[0].fs_id = fs_id;
		prog[0].id = glCreateProgram();
		if (prog[0].id) {
			glAttachShader(prog[0].id, prog[0].vs_id);
			glAttachShader(prog[0].id, fs_id);
			glLinkProgram(prog[0].id);
			program_check(prog[0].id, GL_LINK_STATUS, err_num, err_str);
			if (err_num[0] == 0) {
				glValidateProgram(prog[0].id);
				program_check(prog[0].id, GL_VALIDATE_STATUS, err_num, err_str);
				if (err_num[0] == 0) {
					prog[0].position_att = glGetAttribLocation(prog[0].id, "positionIn");
					prog[0].projection_unif = glGetUniformLocation(prog[0].id, "projection");
					prog[0].data_unif = glGetUniformLocation(prog[0].id, "data");
					const GLenum err_enum = glGetError();
					if (err_enum == GL_INVALID_OPERATION) {
						goDebug(1, 1, 1, 1);
					} else if (err_enum == GL_INVALID_VALUE) {
						goDebug(1, 1, 1, 2);
					}
				}
			}
		} else {
			err_num[0] = err_a;
		}
	}
}

static void gen_vertex_array(GLuint *const vao, const GLint attr, int *const err_num) {
	glGenVertexArrays(1, vao);
	glBindVertexArray(vao[0]);
	glEnableVertexAttribArray(attr);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_OPERATION) {
		err_num[0] = 1043;
	} else if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = 1044;
	}
}

static void gen_array_buffer(GLuint *const vbo, int *const err_num) {
	glGenBuffers(1, vbo);
	glBindBuffer(GL_ARRAY_BUFFER, vbo[0]);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_ENUM) {
		err_num[0] = 1045;
	} else if (err_enum == GL_INVALID_VALUE) {
		err_num[0] = 10346;
	}
}

static void prog_use(const GLuint id, int *const err_num) {
	glUseProgram(id);
	const GLenum err_enum = glGetError();
	if (err_enum == GL_INVALID_VALUE) {
		// invalid parameter
		err_num[0] = 1035;
	} else if (err_enum == GL_INVALID_OPERATION) {
		// invalid operation
		err_num[0] = 1036;
	}
}

static float fbuffer[RECT1K_COUNT * 12];

void g2d_gfx_init(void *const data, int *const err_num, char **const err_str) {
	window_data_t *const wnd_data = (window_data_t*)data;
	wnd_data[0].fs_id = shader_create(GL_FRAGMENT_SHADER, fs_str, 1051, err_num, err_str);
	if (err_num[0] == 0) {
		create_program(&wnd_data[0].prog_rect, vs_rect_str, wnd_data[0].fs_id, 1052, 1053, err_num, err_str);
		if (wnd_data[0].prog_rect.position_att >= 0 && wnd_data[0].prog_rect.projection_unif >= 0 && wnd_data[0].prog_rect.data_unif >= 0) {
			// top left, top right, bottom left, bottom right
			const float vertices[] = { 0.0f, 0.0f, 1.0f, 0.0f, 0.0f, 1.0f, 1.0f, 1.0f };
			gen_vertex_array(&wnd_data[0].prog_rect.vao, wnd_data[0].prog_rect.position_att, err_num);
			gen_array_buffer(&wnd_data[0].prog_rect.vbo, err_num);
			glBufferData(GL_ARRAY_BUFFER, sizeof(float) * 8, vertices, GL_STATIC_DRAW);
			glVertexAttribPointer(wnd_data[0].prog_rect.position_att, 2, GL_FLOAT, GL_FALSE, sizeof(float) * 2, (void*)(0));
			create_program(&wnd_data[0].prog_rect1k, vs_rect1k_str, wnd_data[0].fs_id, 1054, 1055, err_num, err_str);
			if (err_num[0] == 0 && wnd_data[0].prog_rect1k.position_att >= 0 && wnd_data[0].prog_rect1k.projection_unif >= 0 && wnd_data[0].prog_rect1k.data_unif >= 0) {
				int i;
				for (i = 0; i < RECT1K_COUNT; i++) {
					const int index = i*12;
					const float offset = (float)(i*2);
					fbuffer[index] = 0.0f; fbuffer[index+1] = offset + 0.0f;
					fbuffer[index+2] = 1.0f; fbuffer[index+3] = offset + 0.0f;
					fbuffer[index+4] = 0.0f; fbuffer[index+5] = offset + 1.0f;
					fbuffer[index+6] = 0.0f; fbuffer[index+7] = offset + 1.0f;
					fbuffer[index+8] = 1.0f; fbuffer[index+9] = offset + 0.0f;
					fbuffer[index+10] = 1.0f; fbuffer[index+11] = offset + 1.0f;
				}
				gen_vertex_array(&wnd_data[0].prog_rect1k.vao, wnd_data[0].prog_rect1k.position_att, err_num);
				gen_array_buffer(&wnd_data[0].prog_rect1k.vbo, err_num);
				glBufferData(GL_ARRAY_BUFFER, sizeof(float) * RECT1K_COUNT * 12, fbuffer, GL_STATIC_DRAW);
				glVertexAttribPointer(wnd_data[0].prog_rect1k.position_att, 2, GL_FLOAT, GL_FALSE, sizeof(float) * 2, (void*)(0));
				glEnable(GL_BLEND);
				glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
				memcpy(&wnd_data[0].mat_projection, mat_projection, sizeof(float)*(4*4));
			} else if (err_num[0] == 0) {
				err_num[0] = 1040;
			}
		} else {
			err_num[0] = 1033;
		}
	} else {
		err_num[0] = 1038;
	}
}

void g2d_gfx_clear_bg(const float r, const float g, const float b) {
	glClearColor((GLclampf)r, (GLclampf)g, (GLclampf)b, 0.0);
	glClear(GL_COLOR_BUFFER_BIT);
}

void g2d_gfx_swap_buffers(void *const data, int *const err_num, g2d_ul_t *const err_win32) {
	//glFlush();
	//glFinish();
	if (!SwapBuffers(((window_data_t*)data)[0].wnd.ctx.dc)) {
		err_num[0] = 61; err_win32[0] = (g2d_ul_t)GetLastError();
	}
}

void g2d_gfx_set_swap_interval(const int interval) {
	wglSwapIntervalEXT(interval);
}

void g2d_gfx_draw_rect(void *const data, const char *const enabled, const g2d_rect_t *const rects, const int length, const int active, int *const err_num, char **const err_str) {
	int i = 0;
	window_data_t *const wnd_data = (window_data_t*)data;
	if (active >= RECT1K_COUNT) {
		prog_use(wnd_data[0].prog_rect1k.id, err_num);
		if (err_num[0] == 0) {
			int j, k;
			const int limit = (active / RECT1K_COUNT) * RECT1K_COUNT;
			glBindVertexArray(wnd_data[0].prog_rect1k.vao);
			glUniformMatrix4fv(wnd_data[0].prog_rect1k.projection_unif, 1, GL_FALSE, wnd_data[0].mat_projection);
			for (j = 0, k = 0; j < limit; i++) {
				if (enabled[i]) {
					const g2d_rect_t rect = rects[i];
					const int offs = k*9;
					fbuffer[offs] = rect.x; fbuffer[offs+1] = rect.y; fbuffer[offs+2] = rect.w; fbuffer[offs+3] = rect.h; fbuffer[offs+4] = 1.0; fbuffer[offs+5] = rect.r; fbuffer[offs+6] = rect.g; fbuffer[offs+7] = rect.b; fbuffer[offs+8] = rect.a;
					if (++k == RECT1K_COUNT) {
						glUniform1fv(wnd_data[0].prog_rect1k.data_unif, RECT1K_COUNT * 9, fbuffer);
						glDrawArrays(GL_TRIANGLES, 0, RECT1K_COUNT * 6);
						k = 0;
					}
					j++;
				}
			}
		}
	}
	if (i < length && err_num[0] == 0) {
		prog_use(wnd_data[0].prog_rect.id, err_num);
		if (err_num[0] == 0) {
			glBindVertexArray(wnd_data[0].prog_rect.vao);
			glUniformMatrix4fv(wnd_data[0].prog_rect.projection_unif, 1, GL_FALSE, wnd_data[0].mat_projection);
			for (; i < length; i++) {
				if (enabled[i]) {
					const g2d_rect_t rect = rects[i];
					const float data[9] = { rect.x, rect.y, rect.w, rect.h, 1.0, rect.r, rect.g, rect.b, rect.a };
					glUniform1fv(wnd_data[0].prog_rect.data_unif, 9, data);
					glDrawArrays(GL_TRIANGLE_STRIP, 0, 4);
				}
			}
		}
	}
}

void g2d_gfx_set_view_size(void *const data, const int w, const int h) {
	window_data_t *const wnd_data = (window_data_t*)data;
	glViewport((WORD)0, (WORD)0, (WORD)w, (WORD)h);
	wnd_data[0].mat_projection[0] = 2.0f / (float)w;
	wnd_data[0].mat_projection[5] = -2.0f / (float)h;
}
