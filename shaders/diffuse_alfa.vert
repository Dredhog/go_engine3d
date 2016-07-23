#version 410

uniform mat4 mvp_mat;
uniform mat4 model_mat;
uniform mat4 bone_mat[8];

layout (location = 0) in vec3 position;
layout (location = 1) in vec2 texCoord;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec2 bones;
layout (location = 4) in vec2 weights;
out vec4 a_normal;
out vec3 a_position;
out vec2 a_texCoord;

void main()
{
	a_position = position;
	a_normal = vec4(normal.x, normal.y, normal.x,  1.0);
	a_texCoord = texCoord;
	gl_Position = mvp_mat * vec4(position, 1.0);
}
