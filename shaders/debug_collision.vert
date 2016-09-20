#version 330

uniform mat4 vp_mat;
uniform vec4 dynamic_color;
uniform vec3 dynamic_positions[64];
uniform int dynamic_count;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoord;
layout (location = 3) in vec3 color;
layout (location = 4) in vec3 bones;
layout (location = 5) in vec3 weights;

out vec4 Color;

void main(){
	Color = dynamic_color;
	vec3 Position;  
	if(gl_VertexID >= dynamic_count){
		Position = dynamic_positions[dynamic_count-1];
	}else{
		Position = dynamic_positions[gl_VertexID];
	}
	gl_Position = vp_mat * vec4(Position, 1);
}
