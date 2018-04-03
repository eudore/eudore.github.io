#include "stdio.h"
#include "stdlib.h"
#include "string.h"
#include "libmemcached/memcached.h"

int main(int argc, char *argv[]) 
{
	memcached_st *memc;
	memcached_return_t rc;
	memcached_server_st *servers;

	memc = memcached_create(NULL);
	servers = memcached_server_list_append(NULL, "127.0.0.1", 12001, &rc);
	servers = memcached_server_list_append(servers, "127.0.0.1", 12002, &rc);
	servers = memcached_server_list_append(servers, "127.0.0.1", 12003, &rc);
	rc = memcached_server_push(memc, servers);
	memcached_server_free(servers);

	printf("the server count is %d\n",memcached_server_count(memc));  

	char* key="key1";
	char* value = "string";  
	rc = memcached_set(memc,key,strlen(key),value,strlen(value),0,0); 
	//printf("rc = %d\n", rc);


	char return_key[MEMCACHED_MAX_KEY];
	size_t return_key_length; 
	char *return_value;
	size_t return_value_length;

	const char* keys[]= {"key1","key2","key3","key4","key5","key6"}; 
	size_t key_length[]= {4,4,4,4,4,4};
	uint32_t flags;

	return_value = memcached_get(memc,key,strlen(key),&return_value_length,&flags,&rc);
	printf("memcached_get key:%s data:%s\n", key, return_value); 

	rc = memcached_mget(memc, keys, key_length, 6);
	//printf("rc = %d\n", rc);
	while(return_value = memcached_fetch(memc, return_key,&return_key_length, &return_value_length, &flags, &rc)){
		if (rc == MEMCACHED_SUCCESS) {
			printf("Fetch key:%s data:%s\n", return_key, return_value); 
			//printf("return_key_length = %d\n", return_key_length);
			//printf("return_value_length = %d\n", return_value_length);
		}else{
			printf("else rc = %d\n", rc);
		}
	}
	memcached_free(memc);
	return 0;
}
// yum -y install libmemcached libmemcached-devel
// gcc -o test-mem c-memcache.c -I/usr/include/ -L/usr/lib64/ -lmemcached;./test-mem