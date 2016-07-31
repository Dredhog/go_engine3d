#version 410

uniform mat4 mvp_mat;
uniform mat4 model_mat;
uniform mat4 bone_mat[8];

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec2 bones;
layout (location = 4) in vec2 weights;
out vec3 Normal;
out vec3 Position;
out vec2 TexCoord;

void main()
{
	Position = (model_mat * vec4(position, 1.0)).xyz;
	Normal = normalize(normal);
	TexCoord = texCoord;
	gl_Position = mvp_mat * model_mat * vec4(position, 1.0);
}
