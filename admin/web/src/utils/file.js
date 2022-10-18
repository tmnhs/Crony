// 文件处理

// 文件大小
const FileSizes = {
	K: 1024,
	M: 1048576,
	G: 1073741824,
	T: 1099511627776,
}

// 计算文件大小
export function calcFileSize(fileByte) {
	const KB = FileSizes.K
	const MB = FileSizes.M
	const GB = FileSizes.G
	const TB = FileSizes.T
	const FIXED_TWO_POINT = 2
	let fileSizeMsg = ''
	if (fileByte < KB) {
		fileSizeMsg = '小于1K'
	} else if (fileByte > KB && fileByte < MB) {
		fileSizeMsg = (fileByte / KB).toFixed(FIXED_TWO_POINT) + 'K'
	} else if (fileByte === MB) {
		fileSizeMsg = '1M'
	} else if (fileByte > MB && fileByte < GB) {
		fileSizeMsg = (fileByte / (KB * KB)).toFixed(FIXED_TWO_POINT) + 'M'
	} else if (fileByte > MB && fileByte === GB) {
		fileSizeMsg = '1G'
	} else if (fileByte > GB && fileByte < TB) {
		fileSizeMsg = (fileByte / (KB * KB * KB)).toFixed(FIXED_TWO_POINT) + 'G'
	} else {
		fileSizeMsg = '超过1T'
	}
	return fileSizeMsg
}

// 获取文件后缀
export function getFileSuffix(fileName) {
	const pointIndex = fileName.lastIndexOf('.')
	let suffix
	if (pointIndex > -1) {
		suffix = fileName.slice(pointIndex + 1)
	} else {
		suffix = 'file'
	}
	return suffix
}

// 根据二进制数据流下载文件，
// 必须手动指定http请求的responseType为‘blob’或‘'arraybuffer'’
export function downloadByBlob(binaryString, resHeader) {
	// 从响应头的content-disposition中获取文件名称
	const contentDisposition = resHeader['content-disposition']
	const patt = new RegExp('filename=([^;]+\\.[^\\.;]+);*')
	const result = patt.exec(contentDisposition)
	let fileName = ''
	if (result) {
		const reg = /^["](.*)["]$/g
		fileName = decodeURI(result[1].replace(reg, '$1'))
	}

	const blob = new Blob([binaryString], { type: 'application/octet-stream' })
	const link = document.createElement('a')
	link.href = window.URL.createObjectURL(blob)
	link.download = fileName
	link.click()
	// 延时释放
	setTimeout(() => {
		window.URL.revokeObjectURL(link.href)
	}, 100)
}

/**
 * 将base64的图片转换为File
 */
export const base64toFile = (dataURL, filename = 'file') => {
	const arr = dataURL.split(',')
	const mime = arr[0].match(/:(.*?);/)[1]
	const suffix = mime.split('/')[1]
	const bstr = window.atob(arr[1])
	let n = bstr.length
	const u8arr = new Uint8Array(n)
	while (n--) {
		u8arr[n] = bstr.charCodeAt(n)
	}
	return new File([u8arr], `${filename}.${suffix}`, {
		type: mime,
	})
}

/**
 * 将base64的图片转换为blob对象
 */
export const base64ToBlob = dataURL => {
	const arr = dataURL.split(',')
	const mime = arr[0].match(/:(.*?);/)[1] || 'image/jpeg'
	// 去掉url的头，并转化为byte
	const bytes = window.atob(arr[1])
	// 处理异常,将ascii码小于0的转换为大于0
	const ab = new ArrayBuffer(bytes.length)
	// 生成视图（直接针对内存）：8位无符号整数，长度1个字节
	const ia = new Uint8Array(ab)
	for (let i = 0; i < bytes.length; i++) {
		ia[i] = bytes.charCodeAt(i)
	}
	return new Blob([ab], { type: mime })
}
