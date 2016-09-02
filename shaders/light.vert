#version 330

uniform vec3 light_position;
uniform mat4 vp_mat;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec3 bones;
layout (location = 5) in vec3 weights;

void main(){
	gl_Position = vp_mat * vec4(position + light_position, 1);
}
