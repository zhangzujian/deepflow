/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef DF_COMMON_H
#define DF_COMMON_H
#include <stdbool.h>
#include <stdint.h>
#include <fcntl.h>
#include <time.h>
#include <string.h>
#include <stdlib.h>
#include <pthread.h>
#include "types.h"
#include "clib.h"
#include "log.h"

#define __unused __attribute__((__unused__))

#define PORT_NUM_MAX	65536
#define NS_IN_SEC       1000000000ULL
#define NS_IN_MSEC      1000000ULL
#define NS_IN_USEC      1000ULL
#define US_IN_SEC	1000000ULL
#define MS_IN_SEC       1000ULL
#define TIME_TYPE_NAN   1
#define TIME_TYPE_SEC   0

#define OPEN_FILES_MAX 65536
#define KPROBE_EVENTS_FILE "/sys/kernel/debug/tracing/kprobe_events"
#define UPROBE_EVENTS_FILE "/sys/kernel/debug/tracing/uprobe_events"

#ifndef NELEMS
#define NELEMS(a) (sizeof(a) / sizeof((a)[0]))
#endif

#define MAX_PATH_LENGTH 1024
#define CONTAINER_ID_SIZE 65

struct sysinfo {
	long uptime;
	unsigned long loads[3];
	unsigned long totalram;
	unsigned long freeram;
	unsigned long sharedram;
	unsigned long bufferram;
	unsigned long totalswap;
	unsigned long freeswap;
	uint16_t procs;
	uint16_t pad;
	unsigned long totalhigh;
	unsigned long freehigh;
	uint32_t mem_unit;
	char _f[20 - 2 * sizeof(unsigned long) - sizeof(uint32_t)];
};

extern int sysinfo(struct sysinfo *__info);

#ifndef likely
#define likely(x)  __builtin_expect((x),1)
#endif

#ifndef unlikely
#define unlikely(x)  __builtin_expect((x),0)
#endif

#ifndef offsetof
#define offsetof(TYPE, MEMBER)  __builtin_offsetof (TYPE, MEMBER)
#endif

#define CACHE_LINE_SIZE   64
#define CACHE_LINE_MASK (CACHE_LINE_SIZE-1)

#define zfree(P)		\
do {				\
	free((void *)(P));      \
	P = NULL;	        \
} while(0)

static_always_inline void safe_buf_copy(void *dst, int dst_len,
					void *src, int src_len)
{
	if (dst == NULL || src == NULL) {
		ebpf_error("dst:%p, src:%p\n", dst, src);
		return;
	}

	int copy_count = clib_min(dst_len, src_len);
	if (copy_count <= 0) {
		ebpf_error("dst_len:%d, src_len:%d\n", dst_len, src_len);
		return;
	}

	memset(dst, 0, dst_len);
	memcpy(dst, src, copy_count);
}

#if defined(__x86_64__)
#include <emmintrin.h>
static inline void __pause(void)
{
	_mm_pause();
}
#elif defined(__aarch64__)
static inline void __pause(void)
{
	asm volatile("yield" ::: "memory");
}
#else
_Pragma("GCC error \"__pause()\"");
#endif

#define CACHE_LINE_ROUNDUP(size) \
   (CACHE_LINE_SIZE * ((size + CACHE_LINE_SIZE - 1) / CACHE_LINE_SIZE))

#ifndef container_of
#define container_of(ptr, type, member) ({                \
	const typeof(((type *)0)->member) *__p = (ptr);   \
	(type *)( (void *)__p - offsetof(type, member) );})
#endif

#define __cache_aligned __aligned(CACHE_LINE_SIZE)

enum {
	ETR_OK = 0,
	ETR_INVAL = -1,		/* invalid parameter */
	ETR_NOMEM = -2,		/* no memory */
	ETR_EXIST = -3,		/* already exist */
	ETR_NOTEXIST = -4,	/* not exist */
	ETR_NOTGOELF = -5,	/* not go elf */
	ETR_NOPROT = -7,	/* no protocol */
	ETR_NOSYMBOL = -8,	/* no uprobe symbols */
	ETR_UPDATE_MAP_FAILD = -9,	/* update map failed */
	ETR_IDLE = -12,		/* nothing to do */
	ETR_BUSY = -13,		/* resource busy */
	ETR_NOTSUPP = -14,	/* not support */
	ETR_NORESOURCE = -15,	/* no resource */
	ETR_OVERLOAD = -16,	/* overloaded */
	ETR_NOSERV = -17,	/* no service */
	ETR_DISABLED = -18,	/* disabled */
	ETR_NOROOM = -19,	/* no room */
	ETR_NONEALCORE = -20,	/* non-eal thread lcore */
	ETR_CALLBACKFAIL = -21,	/* callbacks fail */
	ETR_IO = -22,		/* I/O error */
	ETR_MSG_FAIL = -23,	/* msg callback failed */
	ETR_MSG_DROP = -24,	/* msg callback dropped */
	ETR_SYSCALL = -26,	/* system call failed */
	ETR_PROC_FAIL = -27,	/* procfs failed */
	ETR_NOHANDLE = -28,	/* not find event handle */
	ETR_LOAD = -29,		/* bpf programe load failed */
	ETR_EPOLL = -30,         /* epoll error */

	/* positive code for non-error */
	ETR_INPROGRESS = 2,	/* in progress */
	ETR_CONTINUE = 4,
	ETR_NEWBUF = 5,
};

struct trace_err_tab {
	int errcode;
	const char *errmsg;
};

static struct trace_err_tab err_tab[] = {
	{ETR_OK, "OK"},
	{ETR_INVAL, "invalid parameter"},
	{ETR_NOMEM, "no memory"},
	{ETR_EXIST, "already exist"},
	{ETR_NOTEXIST, "not exist"},
	{ETR_NOTGOELF, "not go elf"},
	{ETR_NOPROT, "no protocol"},
	{ETR_NOSYMBOL, "no uprobe symbols"},
	{ETR_UPDATE_MAP_FAILD, "update map failed"},
	{ETR_IDLE, "nothing to do"},
	{ETR_BUSY, "resource busy"},
	{ETR_NOTSUPP, "not support"},
	{ETR_NORESOURCE, "no resource"},
	{ETR_OVERLOAD, "overloaded"},
	{ETR_NOSERV, "no service"},
	{ETR_DISABLED, "disabled"},
	{ETR_NOROOM, "no room"},
	{ETR_NONEALCORE, "non-EAL thread lcore"},
	{ETR_CALLBACKFAIL, "callback failed"},
	{ETR_IO, "I/O error"},
	{ETR_MSG_FAIL, "msg callback failed"},
	{ETR_MSG_DROP, "msg dropped"},
	{ETR_SYSCALL, "system call failed"},
	{ETR_PROC_FAIL, "procfs failed"},
	{ETR_NOHANDLE, "not find event handle"},
	{ETR_LOAD, "bpf programe load failed"},
	{ETR_EPOLL, "epoll error"},

	{ETR_INPROGRESS, "in progress"},
};

static inline const char *trace_strerror(int err)
{
	int i;

	for (i = 0; i < NELEMS(err_tab); i++) {
		if (err == err_tab[i].errcode)
			return err_tab[i].errmsg;
	}

	return "<unknow>";
}

#define RUN_ONCE(condition, f, arg) ({    \
  int __ret_warn_once = !!(condition);    \
                                          \
  if (unlikely(__ret_warn_once)) {        \
          f(arg);                         \
          condition = !(__ret_warn_once); \
  }                                       \
})

static inline int is_power_of_2(uint32_t n)
{
	return n && !(n & (n - 1));
}

// Aligns input parameter to the next power of 2
static inline uint32_t align32pow2(uint32_t x)
{
	x--;
	x |= x >> 1;
	x |= x >> 2;
	x |= x >> 4;
	x |= x >> 8;
	x |= x >> 16;

	return x + 1;
}

uint64_t gettime(clockid_t clk_id, int flag);
static inline int64_t get_sysboot_time_ns(void)
{
	int64_t real_time, monotonic_time;
	real_time = gettime(CLOCK_REALTIME, TIME_TYPE_NAN);
	monotonic_time = gettime(CLOCK_MONOTONIC, TIME_TYPE_NAN);
	return (real_time - monotonic_time);
}

bool is_core_kernel(void);
int get_cpus_count(bool **mask);
void clear_residual_probes();
int max_locked_memory_set_unlimited(void);
int sysfs_write(const char *file_name, char *v);
int sysfs_read_num(const char *file_name);
uint32_t get_sys_uptime(void);
u64 get_sys_btime_msecs(void);
u64 get_process_starttime(pid_t pid);
int max_rlim_open_files_set(int num);
int fetch_kernel_version(int *major, int *minor, int *rev, int *num);
unsigned int fetch_kernel_version_code(void);
int get_num_possible_cpus(void);

// Check if task is the main thread based on pid.
// Ignore threads other than the main thread in uprobe to avoid repeating hooks
bool is_user_process(int pid);
bool is_process(int pid);
char *gen_file_name_by_datetime(void);
char *gen_timestamp_prefix(void);
char *gen_timestamp_str(u64 ns);
int fetch_system_type(const char *sys_type, int type_len);
void fetch_linux_release(const char *buf, int buf_len);
u64 get_process_starttime_and_comm(pid_t pid,
				   char *name_base,
				   int len);
int fetch_process_name_from_proc(pid_t pid, char *name, int n_size);
u32 legacy_fetch_log2_page_size(void);
u64 get_netns_id_from_pid(pid_t pid);
bool check_netns_enabled(void);
int get_nspid(int pid);
int get_target_uid_and_gid(int target_pid, int *uid, int *gid);
int copy_file(const char *src_file, const char *dest_file);
int df_enter_ns(int pid, const char *type, int *self_fd);
void df_exit_ns(int fd);
int gen_file_from_mem(const char *mem_ptr, int write_bytes, const char *path);
int exec_command(const char *cmd, const char *args, char *ret_buf, int ret_buf_size);
u64 current_sys_time_secs(void);
int fetch_container_id_from_str(char *buff, char *id, int copy_bytes);
int fetch_container_id_from_proc(pid_t pid, char *id, int copy_bytes);
int parse_num_range(const char *config_str, int bytes_count,
		    bool **mask, int *count);
int parse_num_range_disorder(const char *config_str,
			     int bytes_count, bool ** mask);
int generate_random_integer(int max_value);
bool is_same_netns(int pid);
bool is_same_mntns(int pid);
int is_file_opened_by_other_processes(const char *filename);
/**
 * @brief Find the address through kernel symbols.
 *
 * @param[in] name Kernel symbol name
 * @return 0 indicates that the kernel symbol name was not found, while
 * a non-zero value represents the address of the kernel symbol.
 */
u64 kallsyms_lookup_name(const char *name);
bool substring_starts_with(const char *haystack, const char *needle);
char *get_timestamp_from_us(u64 microseconds);
int find_pid_by_name(const char *process_name, int exclude_pid);
u32 djb2_32bit(const char *str);

/**
 * @brief Format a list of ports into a string, showing ranges of consecutive ports.
 * 
 * This function takes an array of ports, sorts it, and formats it into a string
 * where consecutive port numbers are grouped into ranges. Non-consecutive ports 
 * are shown individually.
 * 
 * @param ports Pointer to an array of uint16_t port numbers.
 * @param size The number of ports in the array.
 * @param ret_str The buffer to store the formatted string.
 * @param str_sz The size of the buffer.
 * 
 * @return None. The result is written to `ret_str`.
 */
void format_port_ranges(uint16_t *ports, size_t size, char *ret_str, int str_sz);

/**
 * @brief Compute 32-bit MurmurHash3 of the given data buffer.
 *
 * MurmurHash3 is a non-cryptographic hash function known for good distribution and speed.
 *
 * @param[in] key   Pointer to the input data buffer to hash.
 * @param[in] len   Length in bytes of the input data.
 * @param[in] seed  Initial seed value for the hash; can be zero or any 32-bit integer.
 *
 * @return 32-bit hash value computed over the input data.
 */
uint32_t murmurhash(const void *key, size_t len, uint32_t seed);
#if !defined(AARCH64_MUSL) && !defined(JAVA_AGENT_ATTACH_TOOL)
int create_work_thread(const char *name, pthread_t *t, void *fn, void *arg);
#endif /* !defined(AARCH64_MUSL) && !defined(JAVA_AGENT_ATTACH_TOOL) */
#endif /* DF_COMMON_H */
