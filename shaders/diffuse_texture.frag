#version 330

uniform vec3 light_position;
uniform sampler2D texture_diffuse0;

in vec3 Position;
in vec3 Normal;
in vec3 Color;
in vec2 TexCoord;

out vec4 outColor;
vec4 white = vec4(1, 1, 1, 1);

void main(){
	outColor = vec4(Color, 1) * texture(texture_diffuse0, TexCoord);
}
