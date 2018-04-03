#include <stdio.h>
#include <stdint.h>
#include <string.h> 

uint32_t jenkins_one_at_a_time_hash(const uint8_t* key, size_t length);

uint32_t jenkins_one_at_a_time_hash(const uint8_t* key, size_t length) {
    size_t i = 0;
    uint32_t hash = 0;
    while (i != length) {
        hash += key[i++];
        hash += hash << 10;
        hash ^= hash >> 6;
    }
    hash += hash << 3;
    hash ^= hash >> 11;
    hash += hash << 15;
    return hash;
}

int main(int argc,char *argv[])  
{
    char* keys[]= {"key1","key2","key3","key4","key5","key6"}; 
    uint32_t n;
    int i;
    for(i=0;i<6;i++){
        n = jenkins_one_at_a_time_hash(keys[i],strlen(keys[i]));
        printf("jenkins hash: k=%s v=%zu\n", keys[i], n);
    }
    return 0;  
}  
/*
jenkins hash: k=key1 v=-1091249108
jenkins hash: k=key2 v=-1389086549
jenkins hash: k=key3 v=-612199097
jenkins hash: k=key4 v=-908791316
jenkins hash: k=key5 v=-193411277
jenkins hash: k=key6 v=-492756092
*/