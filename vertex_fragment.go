package main

var vertexSource = `
#version 410

uniform mat4 mvp_mat;
uniform mat4 model_mat;

in vec3 position;
in vec3 normal;
out vec4 a_normal;
out vec3 a_position;

void main()
{
	a_normal = model_mat * vec4(normal, 1.0);
	a_position = position;
	gl_Position = mvp_mat * model_mat *	 vec4(position, 1.0);
}
` + "\x00"

var fragmentSource = `
#version 410

uniform vec3 light_position;

in vec3 a_position;
in vec4 a_normal;

out vec4 outColor;

void main()
{
	float light_dot_normal = dot(normalize(light_position - a_position), a_normal.xyz);
	float diffuse_alfa = light_dot_normal > 0 ? light_dot_normal/2+0.5 : 0.3;
	vec3 white = vec3(1.0, 1.0, 1.0);
	outColor = vec4((a_normal.xyz/2+0.5)*diffuse_alfa, 1.0);
}
` + "\x00"
