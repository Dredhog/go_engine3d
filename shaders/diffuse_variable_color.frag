#version 330

in vec3 Position;
in vec3 Normal;
in vec3 Color;

out vec4 outColor;
vec3 light_position = vec3(0, 0, 0);

void main(){
	//Phong light diffuse and alfa effect
	float light_dot_normal = dot(normalize(light_position - Position), Normal);
	float diffuse_alpha = max(light_dot_normal, 0.2);

	outColor = vec4(Color, 1) * diffuse_alpha;
	outColor.a = 0.6;
}
