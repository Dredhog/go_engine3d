#version 330

uniform mat4 vp_mat;
uniform mat4 model_mat;
uniform vec3 var_color;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec3 bones;
layout (location = 5) in vec3 weights;

out vec3 Position;
out vec3 Normal;
out vec3 Color;

void main(){
	Position = (model_mat * vec4(position, 1)).xyz;
	Normal = normal;
	Color = var_color;
	gl_Position = vp_mat * model_mat * vec4(position, 1);
}
