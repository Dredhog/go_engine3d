#version 330

uniform mat4 mvp_mat;
uniform mat4 model_mat;
uniform mat4 bone_mat[8];

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec2 bones;
layout (location = 5) in vec2 weights;
out vec3 Normal;
out vec3 Position;
out vec2 TexCoord;
out vec3 Color;

void main()
{
	Position = ((weights.x*bone_mat[int(bones.x)] + weights.y*bone_mat[int(bones.y)]) * vec4(position, 1.0)).xyz;
	Normal = normalize(((transpose(inverse(weights.x*bone_mat[int(bones.x)]))) * vec4(normal, 0.0)).xyz);
	TexCoord = texCoord;
	Color = color;
	gl_Position = mvp_mat * (weights.x*bone_mat[int(bones.x)] + weights.y*bone_mat[int(bones.y)]) * vec4(position, 1.0);
}
