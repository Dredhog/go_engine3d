#version 410

uniform vec3 light_position;
uniform sampler2D sampler_diffuse;

in vec3 Position;
in vec3 Normal;
in vec2 TexCoord;
in vec3 Color;

out vec4 outColor;
vec4 white = vec4(1, 1, 1, 1);

void main()
{
	//Light fade effect
	float light_dist = length(light_position - Position);
	float light_dist_norm = (light_dist < 0) ? -light_dist : light_dist;
	float ligth_dist_squared = light_dist_norm  * light_dist_norm;
	float light_inverse_squared = 10/ligth_dist_squared;
	float light_squared_norm = (light_inverse_squared < 1) ? light_inverse_squared : 1;  
	
	//light surface bounce effect
	float light_dot_normal = dot(normalize(light_position - Position), Normal);
	float diffuse_alfa = max(light_dot_normal, 0.2);

	outColor = light_squared_norm * diffuse_alfa * texture(sampler_diffuse, TexCoord);  
}
