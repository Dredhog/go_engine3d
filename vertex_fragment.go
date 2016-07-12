package main

var vertexSource = `
#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform float time;
uniform float _weightX;
uniform float _weightY;
uniform float _weightZ;
uniform float _period;
uniform float _amplitude;

in vec3 position;
in vec3 color;

out vec3 Color;

void main()
{
	
	Color = color;
	//vec4 waved_position = vec4(position.x, _amplitude * sin( position.x*_weightX + position.y*_weightY + position.z*_weightZ + time/_period), position.z, 1.0);
	
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
