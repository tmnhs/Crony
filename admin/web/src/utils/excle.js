import XLSX from 'xlsx'

/**
 * 导入/导出excel
 */
class ExcelHandler {
	constructor() {
		this.exportExcel = this.exportExcel.bind(this)
		this.readExcelData = this.readExcelData.bind(this)
	}
	/**
	 * 将json数据导出为escel文件
	 * @param {Array} header         表头          如：['姓名', '年龄', '性别', '电话', '电子邮箱']
	 * @param {Array} dataSource     json数据     如：[{name:'wly',age:12, gender:'男'，mobile:'54321',email:'123@qq.com}]
	 * @param {String} fileName      excel文件名称
	 * @param {Boolean} autoWidth   excel中的每个格子是否自动被内容撑开
	 */
	exportExcel(header, dataSource, fileName, autoWidth = true) {
		let sheet
		// 自动设置表格宽度所需要的数据格式与不设置宽度不一样。
		if (autoWidth) {
			const data = dataSource.map(item => Object.values(item))
			data.unshift(header)
			sheet = this.sheet_from_array_of_arrays(data)
			this.adaptWidth(data, sheet)
		} else {
			const data = dataSource.map(i => {
				const values = Object.values(i)
				const newItem = {}
				header.forEach((item, index) => {
					newItem[header[index]] = values[index]
				})
				return newItem
			})
			sheet = XLSX.utils.json_to_sheet(data)
		}

		const sheetName = 'Sheet1'
		const wb = {
			SheetNames: [],
			Sheets: {},
			Props: {},
		}
		wb.SheetNames.push(sheetName)
		wb.Sheets[sheetName] = sheet
		const wbout = XLSX.write(wb, {
			bookType: 'xlsx',
			bookSST: false,
			type: 'binary',
		})
		if (window.btoa) {
			this.downloadURI(wbout, fileName)
		} else {
			this.downloadBlob(wbout, fileName)
		}
	}

	sheet_from_array_of_arrays(data) {
		const sheet = {}
		const range = {
			s: {
				c: 10000000,
				r: 10000000,
			},
			e: {
				c: 0,
				r: 0,
			},
		}
		for (let R = 0; R != data.length; ++R) {
			for (let C = 0; C != data[R].length; ++C) {
				if (range.s.r > R) range.s.r = R
				if (range.s.c > C) range.s.c = C
				if (range.e.r < R) range.e.r = R
				if (range.e.c < C) range.e.c = C
				const cell = {
					v: data[R][C],
				}
				if (cell.v == null) continue
				const cell_ref = XLSX.utils.encode_cell({
					c: C,
					r: R,
				})
				if (typeof cell.v === 'number') {
					cell.t = 'n'
				} else if (typeof cell.v === 'boolean') {
					cell.t = 'b'
				} else if (cell.v instanceof Date) {
					cell.t = 'n'
					cell.z = XLSX.SSF._table[14]
					cell.v = this.datenum(cell.v)
				} else {
					cell.t = 's'
				}
				sheet[cell_ref] = cell
			}
		}
		if (range.s.c < 10000000) {
			sheet['!ref'] = XLSX.utils.encode_range(range)
		}
		return sheet
	}

	datenum(v, date1904) {
		if (date1904) {
			v += 1462
		}
		const epoch = Date.parse(v)
		return (epoch - new Date(Date.UTC(1899, 11, 30))) / (24 * 60 * 60 * 1000)
	}

	// 调整excel表格宽度
	adaptWidth(data, sheet) {
		// 设置worksheet每列的最大宽度
		const colWidth = data.map(row =>
			row.map(val => {
				// 先判断是否为null/undefined
				if (val == null) {
					return {
						wch: 10,
					}
				} else if (val.toString().charCodeAt(0) > 255) {
					// 再判断是否为中文
					return {
						wch: val.toString().length * 2,
					}
				} else {
					return {
						wch: val.toString().length,
					}
				}
			})
		)
		// 以第一行为初始值
		const result = colWidth[0]
		for (let i = 1; i < colWidth.length; i++) {
			for (let j = 0; j < colWidth[i].length; j++) {
				if (result[j]['wch'] < colWidth[i][j]['wch']) {
					result[j]['wch'] = colWidth[i][j]['wch']
				}
			}
		}
		sheet['!cols'] = result
	}

	//二进制字符串转字节流
	s2ab(s) {
		const buffer = new ArrayBuffer(s.length)
		const view = new Uint8Array(buffer)
		for (let i = 0; i < s.length; i++) {
			view[i] = s.charCodeAt(i) & 0xff
		}
		return buffer
	}

	// 使用blob对象下载文件
	downloadBlob(binaryString, name) {
		const stream = this.s2ab(binaryString)
		const blob = new Blob([stream], { type: 'application/octet-stream' })
		const link = document.createElement('a')
		const fileName = name || '数据'
		link.href = window.URL.createObjectURL(blob)
		link.download = fileName + '.xlsx'
		link.click()
		// 延时释放
		setTimeout(function () {
			window.URL.revokeObjectURL(link.href)
		}, 100)
	}

	// 使用base64下载文件
	downloadURI(binaryString, name) {
		const header = 'data:application/octet-stream;base64,'
		const dataURI = window.btoa(binaryString)
		const link = document.createElement('a')
		const fileName = name || '数据'
		link.href = header + dataURI
		link.download = fileName + '.xlsx'
		link.click()
	}

	// 将excel表格中的数据读取为json数据
	readExcelData(file, callback) {
		const fileReader = new FileReader()
		let result = {}
		fileReader.onload = event => {
			const dataSource = event.target.result
			// 以二进制流方式读取得到整份excel表格对象
			const workbook = XLSX.read(dataSource, {
				type: 'binary',
			})
			// 只读取第一个sheet中的数据
			const firstSheetName = workbook.SheetNames[0]
			const worksheet = workbook.Sheets[firstSheetName]
			const header = this.getHeaderRow(worksheet)
			const data = XLSX.utils.sheet_to_json(worksheet)
			callback({
				header,
				data,
			})
		}
		fileReader.readAsBinaryString(file)
	}

	// 获取表头
	getHeaderRow(sheet) {
		const header = []
		const range = XLSX.utils.decode_range(sheet['!ref'])
		const R = range.s.r
		for (let C = range.s.c; C <= range.e.c; C++) {
			const cell =
				sheet[
					XLSX.utils.encode_cell({
						c: C,
						r: R,
					})
				]
			let hdr = 'UNKNOWN' + (C + 1) //设置表头不存在时的默认值
			if (cell && cell.t) {
				hdr = XLSX.utils.format_cell(cell)
			}
			header.push(hdr)
		}
		return header
	}
}

const excelHandler = new ExcelHandler()
const { exportExcel, readExcelData } = excelHandler

export { exportExcel, readExcelData }
