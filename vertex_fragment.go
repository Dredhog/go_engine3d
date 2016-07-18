package main

var vertexSource = `
#version 410

uniform mat4 mvp_mat;
uniform mat4 model_mat;
uniform mat4 bone_mat[3];

in vec3 position;
in vec3 normal;
in float bones[2];
in float weights[2];
out vec4 a_normal;
out vec3 a_position;

void main()
{
	a_position = position;
	a_normal = vec4(normal, 1.0);
	gl_Position = (weights[0]*bone_mat[int(bones[0])] + weights[1]*bone_mat[int(bones[1])])*  vec4(position, 1.0);
}
` + "\x00"

var fragmentSource = `
#version 410

uniform vec3 light_position;

in vec3 a_position;
in vec4 a_normal;

out vec4 outColor;

vec3 white = vec3(1.0, 1.0, 1.0);
vec3 black = vec3(0.0, 0.0, 0.0);

void main()
{
	//float light_dot_normal = dot(normalize(light_position - a_position), a_normal.xyz);
	//float diffuse = light_dot_normal/2 + 0.5;
	//float diffuse_alfa = diffuse >= 0.2 ? diffuse : 0.2;
	//outColor = vec4((white/2+0.5)*diffuse_alfa, 1.0);
	outColor = vec4(black, 1.0);
}
` + "\x00"
