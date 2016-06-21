package main

var vertexSource = `
#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 position;
in vec3 color;

out vec3 Color;

void main()
{
	Color = color;
	gl_Position = projection * camera * model * vec4(position, 1.0);
}
` + "\x00"

var fragmentSource = `
#version 410

in vec3 Color;

out vec4 outColor;

void main()
{
	outColor = vec4(Color, 0.5);
}
` + "\x00"
