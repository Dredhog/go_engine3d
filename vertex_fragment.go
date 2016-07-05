package main

var vertexSource = `
#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform float time;

float _amplitude = 0.3f;

float _weightX = 1f;
float _weightY = 1f;
float _weightZ = 1f;
float _period = 1000000000.0;

in vec3 position;
in vec3 color;

out vec3 Color;

void main()
{
	
	Color = color;
	vec4 waved_position = vec4(position.x,_amplitude * sin( position.x*_weightX + position.y*_weightY + position.z*_weightZ + time/_period), position.z, 1.0);
	
	gl_Position = projection * camera * model * waved_position;
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
