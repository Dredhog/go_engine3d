#version 330

uniform mat4 vp_mat;
uniform float time;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec3 bones;
layout (location = 5) in vec3 weights;

float weightX = 0.2;
float weightY = 0.2;
float weightZ = 0.2;

out vec3 Position;
out vec3 Normal;
out vec3 Color;
out vec2 TexCoord;

void main(){
	Position = position;
	Normal = normal;
	Color = color;
	TexCoord = texCoord;
	vec4 waved_position = vec4(position.x, 0.5 * sin(position.x*weightX + position.y*weightY + position.z*weightZ + time/2) + position.y, position.z, 1);
	gl_Position = vp_mat * waved_position;
}
