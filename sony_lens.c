#include <stdio.h>
#include <stdint.h>
#include <string.h>

typedef struct {
	uint8_t unknown_data0[7]; // 01 00 01 00 01 ea 00
	uint8_t lens_name_len; // max 64(0x40)
	uint8_t lens_name[64];
	uint8_t padding0[320]; // 0
	uint8_t unknown_data1; // 1, not related to focal_length field?
	uint8_t focal_length_msb; // 256mm -> 01, 123mm -> 00
	uint8_t focal_length_lsb; // 256mm -> 00, 123mm -> 7B
	uint8_t use_max_aparture_value; // 1 -> use max_aparture field value
	uint8_t max_aparture; // f7.1 -> 7
	uint8_t max_aparture_decimal; // f7.1 -> 10, f6.3 -> 30
	int8_t visnetting_brightness;
	int8_t visnetting_red;
	int8_t visnetting_blue;
	int8_t ca_red; // ca = Chromatic aberration
	int8_t ca_blue;
	int8_t distortion; // -15~15
	uint8_t padding1[92]; // 0
} LensProfile;

uint8_t unknown_header[] = {0x01, 0x00, 0x01, 0x00, 0x01, 0xea, 0x00};

void main(int argc, char* argv[]){
	FILE *csv = fopen("lenses.csv", "r");
	int i = 0;
	
	while(1){
		LensProfile l1 = {0};
		uint16_t focal_length;
		float max_aparture;
		if (fscanf(csv, "%[^,],%d,%f\n", l1.lens_name, &focal_length, &max_aparture) == EOF || strlen(l1.lens_name) == 0)
		{
			break;
		}
		printf("LENS%04d.BIN: %s %dmm f/%.2f\n", i, l1.lens_name, focal_length, max_aparture);
			
		memcpy(l1.unknown_data0, unknown_header, 7);
			
		l1.lens_name_len = strlen(l1.lens_name);

		
		l1.focal_length_msb = (focal_length >> 8);
		l1.focal_length_lsb = (focal_length & 0xff);

		l1.use_max_aparture_value = 1;
		l1.max_aparture = (uint8_t)max_aparture;
		l1.max_aparture_decimal = (uint16_t)(max_aparture * 100) - (uint16_t)((uint8_t)max_aparture * 100);
		
		//l1.visnetting_brightness = 1;
		//l1.visnetting_red = -3;
		//l1.visnetting_blue = 5;
		//l1.ca_red = 7;
		//l1.ca_blue = -4;
		//l1.distortion = -12;
		char fn[16];
		sprintf(fn, "LENS%04d.BIN", i);
		FILE *fp = fopen(fn, "wb");
		fwrite((void *)&l1, sizeof(LensProfile), 1, fp);
		fclose(fp);
		
		i++;
	}
}
