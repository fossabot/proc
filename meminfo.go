/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package proc

import (
	"github.com/ennoo/rivet/utils/file"
	"github.com/ennoo/rivet/utils/log"
	str "github.com/ennoo/rivet/utils/string"
	"strings"
)

// MemInfo 存储器使用信息，包括物理内存和交换内存
type MemInfo struct {
	MemTotal          string // 所有可用RAM大小 （即物理内存减去一些预留位和内核的二进制代码大小）
	MemFree           string // LowFree与HighFree的总和，被系统留着未使用的内存
	MemAvailable      string // 有些应用程序会根据系统的可用内存大小自动调整内存申请的多少，所以需要一个记录当前可用内存数量的统计值，MemFree并不适用，因为MemFree不能代表全部可用的内存，系统中有些内存虽然已被使用但是可以回收的，比如cache/buffer、slab都有一部分可以回收，所以这部分可回收的内存加上MemFree才是系统可用的内存，即MemAvailable。/proc/meminfo中的MemAvailable是内核使用特定的算法估算出来的，要注意这是一个估计值，并不精确。
	Buffers           string // 用来给文件做缓冲大小
	Cached            string // 被高速缓冲存储器（cache memory）用的内存的大小（等于 diskcache minus SwapCache ）
	SwapCached        string // 被高速缓冲存储器（cache memory）用的交换空间的大小。已经被交换出来的内存，但仍然被存放在swap file中。用来在需要的时候很快的被替换而不需要再次打开I/O端口
	Active            string // 在活跃使用中的缓冲或高速缓冲存储器页面文件的大小，除非非常必要否则不会被移作他用
	Inactive          string // 在不经常使用中的缓冲或高速缓冲存储器页面文件的大小，可能被用于其他途径
	ActiveAnon        string //
	InactiveAnon      string //
	ActiveFile        string //
	InactiveFile      string //
	Unevictable       string //
	MLocked           string //
	SwapTotal         string // 交换空间的总大小
	SwapFree          string // 未被使用交换空间的大小
	Dirty             string // 等待被写回到磁盘的内存大小
	WriteBack         string // 正在被写回到磁盘的内存大小
	AnonPages         string // 未映射页的内存大小
	Mapped            string // 设备和文件等映射的大小
	Shmem             string //
	Slab              string // 内核数据结构缓存的大小，可以减少申请和释放内存带来的消耗
	SReclaimable      string // 可收回Slab的大小
	SUnreclaim        string // 不可收回Slab的大小（SUnreclaim+SReclaimable＝Slab）
	KernelStack       string // 每一个用户线程都会分配一个kernel stack（内核栈），内核栈虽然属于线程，但用户态的代码不能访问，只有通过系统调用(syscall)、自陷(trap)或异常(exception)进入内核态的时候才会用到，也就是说内核栈是给kernel code使用的。在x86系统上Linux的内核栈大小是固定的8K或16K
	PageTables        string // 管理内存分页页面的索引表的大小
	NFSUnstable       string // 不稳定页表的大小
	Bounce            string // 有些老设备只能访问低端内存，比如16M以下的内存，当应用程序发出一个I/O 请求，DMA的目的地址却是高端内存时（比如在16M以上），内核将在低端内存中分配一个临时buffer作为跳转，把位于高端内存的缓存数据复制到此处。这种额外的数据拷贝被称为“bounce buffering”，会降低I/O 性能。大量分配的bounce buffers 也会占用额外的内存。
	WriteBackTmp      string //
	CommitLimit       string //
	CommittedAS       string //
	VMAllocTotal      string // 可以vmalloc虚拟内存大小
	VMAllocUsed       string // 已经被使用的虚拟内存大小
	VMAllocChunk      string //
	HardwareCorrupted string // 当系统检测到内存的硬件故障时，会把有问题的页面删除掉，不再使用，/proc/meminfo中的HardwareCorrupted统计了删除掉的内存页的总大小。
	AnonHugePages     string //
	CmaTotal          string //
	CmaFree           string //
	HugePagesTotal    string // 对应内核参数 vm.nr_hugepages，也可以在运行中的系统上直接修改 /proc/sys/vm/nr_hugepages，修改的结果会立即影响空闲内存 MemFree的大小，因为HugePages在内核中独立管理，只要一经定义，无论是否被使用，都不再属于free memory。
	HugePagesFree     string //
	HugePagesRsvd     string //
	HugePagesSurp     string //
	HugePageSize      string //
	DirectMap4k       string //
	DirectMap2M       string //
	DirectMap1G       string //
}

// FormatMemInfo 将文件内容转为 MemInfo 对象
func (m *MemInfo) FormatMemInfo(filePath string) {
	data, err := file.ReadFileByLine(filePath)
	if nil != err {
		log.Self.Error("read mem info error", log.Error(err))
	} else {
		for index := range data {
			m.formatMemInfo(data[index])
		}
	}
}

func (m *MemInfo) formatMemInfo(lineStr string) {
	if strings.HasPrefix(lineStr, "MemTotal") {
		m.MemTotal = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "MemFree") {
		m.MemFree = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "MemAvailable") {
		m.MemAvailable = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Buffers") {
		m.Buffers = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Cached") {
		m.Cached = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "SwapCached") {
		m.SwapCached = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Active") {
		m.Active = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Inactive") {
		m.Inactive = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Active(anon)") {
		m.ActiveAnon = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Inactive(anon)") {
		m.InactiveAnon = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Active(file)") {
		m.ActiveFile = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Inactive(file)") {
		m.InactiveFile = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Unevictable") {
		m.Unevictable = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Mlocked") {
		m.MLocked = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "SwapTotal") {
		m.SwapTotal = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "SwapFree") {
		m.SwapFree = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Dirty") {
		m.Dirty = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Writeback") {
		m.WriteBack = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "AnonPages") {
		m.AnonPages = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Mapped") {
		m.Mapped = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Shmem") {
		m.Shmem = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Slab") {
		m.Slab = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "SReclaimable") {
		m.SReclaimable = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "SUnreclaim") {
		m.SUnreclaim = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "KernelStack") {
		m.KernelStack = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "PageTables") {
		m.PageTables = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "NFS_Unstable") {
		m.NFSUnstable = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Bounce") {
		m.Bounce = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "WritebackTmp") {
		m.WriteBackTmp = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "CommitLimit") {
		m.CommitLimit = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Committed_AS") {
		m.CommittedAS = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "VmallocTotal") {
		m.VMAllocTotal = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "VmallocUsed") {
		m.VMAllocUsed = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "VmallocChunk") {
		m.VMAllocChunk = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "HardwareCorrupted") {
		m.HardwareCorrupted = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "AnonHugePages") {
		m.AnonHugePages = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "CmaTotal") {
		m.CmaTotal = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "CmaFree") {
		m.CmaFree = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "HugePages_Total") {
		m.HugePagesTotal = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "HugePages_Free") {
		m.HugePagesFree = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "HugePages_Rsvd") {
		m.HugePagesRsvd = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "HugePages_Surp") {
		m.HugePagesSurp = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "Hugepagesize") {
		m.HugePageSize = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "DirectMap4k") {
		m.DirectMap4k = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "DirectMap2M") {
		m.DirectMap2M = str.Trim(strings.Split(lineStr, ":")[1])
	} else if strings.HasPrefix(lineStr, "DirectMap1G") {
		m.DirectMap1G = str.Trim(strings.Split(lineStr, ":")[1])
	}
}
