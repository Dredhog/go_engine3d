#version 330

uniform mat4 vp_mat;
uniform mat4 model_mat;
uniform mat4 model_rotation_mat;
uniform mat4 bone_mat[15];
uniform int  edit_bone;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec3 bones;
layout (location = 5) in vec3 weights;

out vec3 Normal;
out vec3 Position;
out vec2 TexCoord;
out vec4 Color;

vec4 white = vec4(1, 1, 1, 1);
vec4 red   = vec4(1, 0, 0, 1);

void main()
{
	vec3 Weights = weights / ((weights.x + weights.y + weights.z != 0) ? weights.x+weights.y+weights.z : 1);
	Position = (model_mat * (Weights.x*bone_mat[int(bones.x)] + Weights.y*bone_mat[int(bones.y)] + Weights.z*bone_mat[int(bones.z)]) * vec4(position, 1.0)).xyz;
	Normal = normalize((model_rotation_mat * (transpose(inverse(Weights.x*bone_mat[int(bones.x)] + Weights.y*bone_mat[int(bones.y)]))) * vec4(normal, 0.0)).xyz);
	TexCoord = texCoord;
	Color = white;
	if(int(bones.x) == edit_bone) {
		Color = Weights.x*red + (1-Weights.x)*white;
	}
	if(int(bones.y) == edit_bone) {
		Color = Weights.y*red + (1-Weights.y)*white;
	}
	if(int(bones.z) == edit_bone) {
		Color = Weights.z*red + (1-Weights.z)*white;
	}
	gl_Position = vp_mat * model_mat * (Weights.x*bone_mat[int(bones.x)] + Weights.y*bone_mat[int(bones.y)] + Weights.z*bone_mat[int(bones.z)]) * vec4(position, 1);
}
