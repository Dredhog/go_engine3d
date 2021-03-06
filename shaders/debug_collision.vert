#version 330

uniform mat4 vp_mat;
uniform vec4 var_color;
uniform vec3 var_positions[100];
uniform int var_count;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec3 bones;
layout (location = 5) in vec3 weights;

vec4 blue = vec4(0, 0, 1, 1);
out vec4 Color;

void main(){
	Color = var_color;
	vec3 Position;  
	if(gl_VertexID >= var_count){
		Position = var_positions[var_count-1];
	}else{
		Position = var_positions[gl_VertexID];
	}
	if ( gl_VertexID == 3){
		Color = blue;
	}
	gl_Position = vp_mat * vec4(Position, 1);
}
