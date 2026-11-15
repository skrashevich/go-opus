/* Independent packet-file decoder driver; does not share the Go/Python harness. */
#include <opus.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char **argv) {
  if (argc != 8) {fprintf(stderr,"usage: %s rate channels int16|float packets output ranges frame_size\n",argv[0]);return 2;}
  int frame=atoi(argv[7]); int rate=atoi(argv[1]), channels=atoi(argv[2]), floating=!strcmp(argv[3],"float"), err=0;
  if (frame<1 || frame>5760 || channels<1 || channels>2) return 2;
  FILE *input=fopen(argv[4],"rb"), *output=fopen(argv[5],"wb"), *ranges=fopen(argv[6],"w");
  if (!input || !output || !ranges) return 3;
  OpusDecoder *decoder=opus_decoder_create(rate,channels,&err);
  if (!decoder || err) {fprintf(stderr,"create error %d\n",err);return 4;}
  opus_int16 pcm[5760*2]; float fp[5760*2]; unsigned char packet[4000], size[4]; int frames=0;
  for (;;) {
    size_t got=fread(size,1,4,input);if (!got && feof(input)) break;
    if (got!=4) return 5;
    uint32_t len=(uint32_t)size[0] | (uint32_t)size[1]<<8 | (uint32_t)size[2]<<16 | (uint32_t)size[3]<<24;
    if (!len || len>sizeof(packet) || fread(packet,1,len,input)!=len) return 6;
    int n=floating ? opus_decode_float(decoder,packet,(opus_int32)len,fp,frame,0) : opus_decode(decoder,packet,(opus_int32)len,pcm,frame,0);
    if (n!=frame) {fprintf(stderr,"frame %d: returned %d\n",frames,n);return 7;}
    if (floating) {if (fwrite(fp,sizeof(float),n*channels,output)!=(size_t)(n*channels))return 8;}
    else {if (fwrite(pcm,sizeof(opus_int16),n*channels,output)!=(size_t)(n*channels))return 8;}
    opus_uint32 final_range=0;
    if (opus_decoder_ctl(decoder,OPUS_GET_FINAL_RANGE(&final_range)))return 9;
    fprintf(ranges,"%u\n",final_range); frames++;
  }
  fprintf(stderr,"%s; rate=%d; channels=%d; float=%d; frames=%d\n",opus_get_version_string(),rate,channels,floating,frames);
  opus_decoder_destroy(decoder);fclose(input);fclose(output);fclose(ranges);return 0;
}
