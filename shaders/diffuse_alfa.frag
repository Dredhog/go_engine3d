#version 410

uniform vec3 light_position;
uniform sampler2D diffuseTexture;

in vec3 a_position;
in vec4 a_normal;
in vec2 a_texCoord;

out vec4 outColor;

void main()
{
	float light_dot_normal = dot(normalize(light_position - a_position), a_normal.xyz);
	float diffuse = light_dot_normal/2 + 0.5;
	float diffuse_alfa = diffuse >= 0.05 ? diffuse : 0.05;
	outColor = diffuse_alfa * texture(diffuseTexture, a_texCoord); 
}
